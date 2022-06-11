package apps

import "github.com/crawlab-team/crawlab-core/interfaces"

type ServerOption func(app ServerApp)

func WithServerConfigPath(path string) ServerOption {
	return func(app ServerApp) {
		app.SetConfigPath(path)
	}
}

func WithServerGrpcAddress(address interfaces.Address) ServerOption {
	return func(app ServerApp) {
		app.SetGrpcAddress(address)
	}
}
