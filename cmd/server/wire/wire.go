//go:build wireinject
// +build wireinject

package wire

import (
	commonConfig "PandoraFuclaudePlusHelper/config"
	"PandoraFuclaudePlusHelper/internal/handler"
	"PandoraFuclaudePlusHelper/internal/middleware"
	"PandoraFuclaudePlusHelper/internal/repository"
	"PandoraFuclaudePlusHelper/internal/server"
	"PandoraFuclaudePlusHelper/internal/service"
	"PandoraFuclaudePlusHelper/pkg/app"
	"PandoraFuclaudePlusHelper/pkg/jwt"
	"PandoraFuclaudePlusHelper/pkg/log"
	serverType "PandoraFuclaudePlusHelper/pkg/server"
	"PandoraFuclaudePlusHelper/pkg/server/http"
	"PandoraFuclaudePlusHelper/pkg/server/reverse/claude"
	"PandoraFuclaudePlusHelper/pkg/server/reverse/openai"
	"PandoraFuclaudePlusHelper/pkg/sid"
	"github.com/google/wire"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	// repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewOpenaiTokenRepository,
	repository.NewOpenaiAccountRepository,
	repository.NewClaudeTokenRepository,
	repository.NewClaudeAccountRepository,
	repository.NewConversationRepository,
	repository.NewUserRepository,
)

var serviceCoordinatorSet = wire.NewSet(
	service.NewServiceCoordinator,
)

var serviceSet = wire.NewSet(
	service.NewService,
	serviceCoordinatorSet,
	service.NewLoginService,
	service.NewUserService,
	service.NewOpenaiTokenService,
	service.NewOpenaiAccountService,
	service.NewClaudeTokenService,
	service.NewClaudeAccountService,
	server.NewTask,
)

var migrateSet = wire.NewSet(
	server.NewMigrate,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewLoginHandler,
	handler.NewUserHandler,
	handler.NewOpenaiTokenHandler,
	handler.NewOpenaiAccountHandler,
	handler.NewClaudeTokenHandler,
	handler.NewClaudeAccountHandler,
)

var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewChatGPTReverseProxyServer,
	server.NewClaudeReverseProxyServer,
	server.NewJob,
)

// build App
func newApp(httpServer *http.Server, openaiServer *openai.Server, claudeServer *claude.Server, job *server.Job, task *server.Task, migrate *server.Migrate) *app.App {
	servers := []serverType.Server{
		httpServer,
		job,
		migrate,
		openaiServer,
		claudeServer,
	}
	if commonConfig.GetConfig().EnableTask {
		servers = append(servers, task)
	}
	return app.NewApp(
		app.WithServer(servers...),
		app.WithName("demo-server"),
	)
}

func NewWire(*log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		serverSet,
		migrateSet,
		sid.NewSid,
		jwt.NewJwt,
		middleware.NewConversationLoggerMiddleware,
		newApp,
	))

}
