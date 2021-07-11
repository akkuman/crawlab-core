package interfaces

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TaskHandlerService interface {
	TaskBaseService
	// Run task and execute locally
	Run(taskId primitive.ObjectID) (err error)
	// Cancel task locally
	Cancel(taskId primitive.ObjectID) (err error)
	// ReportStatus periodically report handler status to master
	ReportStatus()
	// Reset reset internals to default
	Reset()
	// IsSyncLocked whether the given spider is locked for files sync
	IsSyncLocked(spiderId primitive.ObjectID) (ok bool)
	// LockSync lock files sync for given spider
	LockSync(spiderId primitive.ObjectID)
	// UnlockSync unlock files sync for given spider
	UnlockSync(spiderId primitive.ObjectID)
	// GetExitWatchDuration get max runners
	GetExitWatchDuration() (duration time.Duration)
	// SetExitWatchDuration set max runners
	SetExitWatchDuration(duration time.Duration)
	// GetReportInterval get report interval
	GetReportInterval() (interval time.Duration)
	// SetReportInterval set report interval
	SetReportInterval(interval time.Duration)
	// GetModelService get model service
	GetModelService() (modelSvc GrpcClientModelService)
	// GetModelSpiderService get model spider service
	GetModelSpiderService() (modelSpiderSvc GrpcClientModelSpiderService)
	// GetModelTaskService get model task service
	GetModelTaskService() (modelTaskSvc GrpcClientModelTaskService)
	// GetModelTaskStatService get model task stat service
	GetModelTaskStatService() (modelTaskStatSvc GrpcClientModelTaskStatService)
	// GetNodeConfigService get node config service
	GetNodeConfigService() (cfgSvc NodeConfigService)
	// GetCurrentNode get node of the handler
	GetCurrentNode() (n Node, err error)
	// GetTaskById get task by id
	GetTaskById(id primitive.ObjectID) (t Task, err error)
	// GetSpiderById get task by id
	GetSpiderById(id primitive.ObjectID) (t Spider, err error)
}
