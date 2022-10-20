package server

import (
	"context"
	"encoding/json"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/task/stats"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/crawlab-db/mongo"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"io"
	"strings"
)

type TaskServer struct {
	grpc.UnimplementedTaskServiceServer

	// dependencies
	modelSvc service.ModelService
	cfgSvc   interfaces.NodeConfigService
	statsSvc interfaces.TaskStatsService

	// internals
	server interfaces.GrpcServer
}

// Subscribe to task stream when a task runner in a node starts
func (svr TaskServer) Subscribe(stream grpc.TaskService_SubscribeServer) (err error) {
	for {
		msg, err := stream.Recv()
		utils.LogDebug(msg.String())
		if err == io.EOF {
			return nil
		}
		if err != nil {
			if strings.HasSuffix(err.Error(), "context canceled") {
				return nil
			}
			trace.PrintError(err)
			continue
		}
		switch msg.Code {
		case grpc.StreamMessageCode_INSERT_DATA:
			err = svr.handleInsertData(msg)
		case grpc.StreamMessageCode_INSERT_LOGS:
			err = svr.handleInsertLogs(msg)
		default:
			err = errors.ErrorGrpcInvalidCode
			log.Errorf("invalid stream message code: %d", msg.Code)
			continue
		}
		if err != nil {
			log.Errorf("grpc error[%d]: %v", msg.Code, err)
		}
	}
}

// Fetch tasks to be executed by a task handler
func (svr TaskServer) Fetch(ctx context.Context, request *grpc.Request) (response *grpc.Response, err error) {
	nodeKey := request.GetNodeKey()
	if nodeKey == "" {
		return nil, trace.TraceError(errors.ErrorGrpcInvalidNodeKey)
	}
	n, err := svr.modelSvc.GetNodeByKey(nodeKey, nil)
	if err != nil {
		return nil, trace.TraceError(err)
	}
	var tid primitive.ObjectID
	opts := &mongo.FindOptions{
		Sort: bson.D{
			{"p", 1},
			{"_id", 1},
		},
		Limit: 1,
	}
	if err := mongo.RunTransactionWithContext(ctx, func(sc mongo2.SessionContext) (err error) {
		// get task queue item assigned to this node
		tid, err = svr.getTaskQueueItemIdAndDequeue(bson.M{"nid": n.Id}, opts)
		if err != nil {
			return err
		}
		if !tid.IsZero() {
			return nil
		}

		// get task queue item assigned to any node (random mode)
		tid, err = svr.getTaskQueueItemIdAndDequeue(bson.M{"nid": nil}, opts)
		if !tid.IsZero() {
			return nil
		}
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return HandleSuccessWithData(tid)
}

func (svr TaskServer) handleInsertData(msg *grpc.StreamMessage) (err error) {
	data, err := svr.deserialize(msg)
	if err != nil {
		return err
	}
	var records []interface{}
	for _, d := range data.Records {
		res, ok := d[constants.TaskKey]
		if ok {
			switch res.(type) {
			case string:
				id, err := primitive.ObjectIDFromHex(res.(string))
				if err == nil {
					d[constants.TaskKey] = id
				}
			}
		}
		records = append(records, d)
	}
	return svr.statsSvc.InsertData(data.TaskId, records...)
}

func (svr TaskServer) handleInsertLogs(msg *grpc.StreamMessage) (err error) {
	data, err := svr.deserialize(msg)
	if err != nil {
		return err
	}
	return svr.statsSvc.InsertLogs(data.TaskId, data.Logs...)
}

func (svr TaskServer) getTaskQueueItemIdAndDequeue(query bson.M, opts *mongo.FindOptions) (tid primitive.ObjectID, err error) {
	var tq models.TaskQueueItem
	// get task queue item assigned to this node
	if err := mongo.GetMongoCol(interfaces.ModelColNameTaskQueue).Find(query, opts).One(&tq); err != nil {
		if err == mongo2.ErrNoDocuments {
			return tid, nil
		}
		return tid, trace.TraceError(err)
	}
	if err := delegate.NewModelDelegate(&tq).Delete(); err != nil {
		return tid, trace.TraceError(err)
	}
	return tq.Id, nil
}

func (svr TaskServer) deserialize(msg *grpc.StreamMessage) (data entity.StreamMessageTaskData, err error) {
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return data, trace.TraceError(err)
	}
	if data.TaskId.IsZero() {
		return data, trace.TraceError(errors.ErrorGrpcInvalidType)
	}
	return data, nil
}

func NewTaskServer(opts ...TaskServerOption) (res *TaskServer, err error) {
	// task server
	svr := &TaskServer{}

	// apply options
	for _, opt := range opts {
		opt(svr)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		return nil, err
	}
	if err := c.Provide(stats.ProvideGetTaskStatsService(svr.server.GetConfigPath())); err != nil {
		return nil, err
	}
	if err := c.Provide(config.ProvideConfigService(svr.server.GetConfigPath())); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(
		modelSvc service.ModelService,
		statsSvc interfaces.TaskStatsService,
		cfgSvc interfaces.NodeConfigService,
	) {
		svr.modelSvc = modelSvc
		svr.statsSvc = statsSvc
		svr.cfgSvc = cfgSvc
	}); err != nil {
		return nil, err
	}

	return svr, nil
}

func ProvideTaskServer(server interfaces.GrpcServer, opts ...TaskServerOption) func() (res *TaskServer, err error) {
	return func() (*TaskServer, error) {
		opts = append(opts, WithServerTaskServerService(server))
		return NewTaskServer(opts...)
	}
}
