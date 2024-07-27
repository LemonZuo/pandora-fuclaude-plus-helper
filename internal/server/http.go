package server

import (
	"PandoraFuclaudePlusHelper"
	commonConfig "PandoraFuclaudePlusHelper/config"
	"PandoraFuclaudePlusHelper/internal/handler"
	"PandoraFuclaudePlusHelper/internal/middleware"
	"PandoraFuclaudePlusHelper/pkg/jwt"
	"PandoraFuclaudePlusHelper/pkg/log"
	"PandoraFuclaudePlusHelper/pkg/server/http"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	httpcore "net/http"
)

func NewHTTPServer(
	logger *log.Logger,
	jwt *jwt.JWT,
	loginHandler *handler.LoginHandler,
	openaiAccountHandler *handler.OpenaiAccountHandler,
	openaiTokenHandler *handler.OpenaiTokenHandler,
	userHandler *handler.UserHandler,
	claudeTokenHandler *handler.ClaudeTokenHandler,
	claudeAccountHandler *handler.ClaudeAccountHandler,
) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(commonConfig.GetConfig().HttpHost),
		http.WithServerPort(commonConfig.GetConfig().HttpPort),
	)

	s.Use(static.Serve("/", static.EmbedFolder(PandoraFuclaudePlusHelper.EmbedFrontendFS, "frontend/dist")))

	v1 := s.Group("/api")
	{
		noAuthRouter := v1.Group("/")
		{
			noAuthRouter.POST("/auth", loginHandler.Login)
			noAuthRouter.POST("/info", func(c *gin.Context) {
				c.JSON(httpcore.StatusOK, gin.H{
					"message": "ok",
					"status":  0,
					"data": gin.H{
						"version":    "1.0.0",
						"systemName": "PandoraFuclaudePlusHelper",
						"startTime":  commonConfig.GetConfig().StartTime.Format("2006-01-02 15:04:05"),
						"status":     true,
					},
				})
			})
		}

		userAuthRouter := v1.Group("/user").Use(middleware.StrictAuth(jwt, logger))
		{
			userAuthRouter.POST("/add", userHandler.CreateUser)
			// userAuthRouter.POST("/refresh", openaiTokenHandler.RefreshToken)
			userAuthRouter.POST("/search", userHandler.SearchUser)
			userAuthRouter.POST("/delete", userHandler.DeleteUser)
			userAuthRouter.POST("/update", userHandler.UpdateUser)
		}

		tokenAuthRouter := v1.Group("/openai-token").Use(middleware.StrictAuth(jwt, logger))
		{
			tokenAuthRouter.POST("/add", openaiTokenHandler.CreateToken)
			tokenAuthRouter.POST("/refresh", openaiTokenHandler.RefreshToken)
			tokenAuthRouter.POST("/search", openaiTokenHandler.SearchToken)
			tokenAuthRouter.POST("/delete", openaiTokenHandler.DeleteToken)
			tokenAuthRouter.POST("/update", openaiTokenHandler.UpdateToken)
		}

		accountAuthRouter := v1.Group("/openai-account").Use(middleware.StrictAuth(jwt, logger))
		{
			accountAuthRouter.POST("/add", openaiAccountHandler.CreateAccount)
			accountAuthRouter.POST("/update", openaiAccountHandler.UpdateAccount)
			accountAuthRouter.POST("/delete", openaiAccountHandler.DeleteAccount)
			accountAuthRouter.POST("/search", openaiAccountHandler.SearchAccount)
			accountAuthRouter.POST("/statistic", openaiAccountHandler.StatisticAccount)
			accountAuthRouter.POST("/disable", openaiAccountHandler.DisableAccount)
			accountAuthRouter.POST("/enable", openaiAccountHandler.EnableAccount)
		}

		claudeTokenAuthRouter := v1.Group("/claude-token").Use(middleware.StrictAuth(jwt, logger))
		{
			claudeTokenAuthRouter.POST("/add", claudeTokenHandler.CreateToken)
			// claudeTokenAuthRouter.POST("/refresh", claudeTokenHandler.RefreshToken)
			claudeTokenAuthRouter.POST("/search", claudeTokenHandler.SearchToken)
			claudeTokenAuthRouter.POST("/delete", claudeTokenHandler.DeleteToken)
			claudeTokenAuthRouter.POST("/update", claudeTokenHandler.UpdateToken)
		}

		claudeAccountAuthRouter := v1.Group("/claude-account").Use(middleware.StrictAuth(jwt, logger))
		{
			claudeAccountAuthRouter.POST("/add", claudeAccountHandler.CreateAccount)
			claudeAccountAuthRouter.POST("/update", claudeAccountHandler.UpdateAccount)
			claudeAccountAuthRouter.POST("/delete", claudeAccountHandler.DeleteAccount)
			claudeAccountAuthRouter.POST("/search", claudeAccountHandler.SearchAccount)
			claudeAccountAuthRouter.POST("/statistic", claudeAccountHandler.StatisticAccount)
			claudeAccountAuthRouter.POST("/disable", claudeAccountHandler.DisableAccount)
			claudeAccountAuthRouter.POST("/enable", claudeAccountHandler.EnableAccount)
		}
	}

	return s
}
