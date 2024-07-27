//go:build wireinject
// +build wireinject

package wire

import (
	"PandoraFuclaudePlusHelper/internal/handler"
	"PandoraFuclaudePlusHelper/internal/repository"
	"PandoraFuclaudePlusHelper/internal/server"
	"PandoraFuclaudePlusHelper/internal/service"
	"PandoraFuclaudePlusHelper/pkg/app"
	"PandoraFuclaudePlusHelper/pkg/jwt"
	"PandoraFuclaudePlusHelper/pkg/log"
	"PandoraFuclaudePlusHelper/pkg/server/http"
	"PandoraFuclaudePlusHelper/pkg/sid"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	// repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewOpenaiAccountRepository,
	repository.NewShareRepository,
)

var serviceCoordinatorSet = wire.NewSet(
	service.NewServiceCoordinator,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	serviceCoordinatorSet,
	service.NewOpenaiAccountService,
	service.NewShareService,
	server.NewTask,
)

var migrateSet = wire.NewSet(
	server.NewMigrate,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewShareHandler,
	handler.NewOpenaiAccountHandler,
)

var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJob,
)

// build App
func newApp(httpServer *http.Server, job *server.Job, task *server.Task, migrate *server.Migrate) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, job, task, migrate),
		app.WithName("demo-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		serverSet,
		migrateSet,
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))

}
