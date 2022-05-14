package apps

import "github.com/crawlab-team/crawlab-core/interfaces"

type App interface {
	Init()
	Start()
	Wait()
	Stop()
}

type NodeApp interface {
	App
	interfaces.WithConfigPath
	SetGrpcAddress(address interfaces.Address)
}

type MasterApp interface {
	NodeApp
	SetRunOnMaster(ok bool)
	GetApi() (api *Api)
	GetMasterService() (masterSvc interfaces.NodeMasterService)
}

type WorkerApp interface {
	NodeApp
}