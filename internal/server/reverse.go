package server

import (
	commonConfig "PandoraFuclaudePlusHelper/config"
	"PandoraFuclaudePlusHelper/internal/middleware"
	"PandoraFuclaudePlusHelper/pkg/log"
	"PandoraFuclaudePlusHelper/pkg/server/reverse/claude"
	"PandoraFuclaudePlusHelper/pkg/server/reverse/openai"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewChatGPTReverseProxyServer 创建 ChatGPT 反向代理服务器
func NewChatGPTReverseProxyServer(
	logger *log.Logger,
	conversationLoggerMiddleware *middleware.ConversationLoggerMiddleware,
) *openai.Server {
	r := gin.Default()

	// 创建反向代理处理函数
	proxyHandler := reverseProxy(commonConfig.GetConfig().OpenAiSite)

	if commonConfig.GetConfig().ModerationEnable() {
		r.POST("/backend-api/conversation", middleware.OpenAiContentModerationMiddleware(logger), proxyHandler)
	} else {
		r.POST("/backend-api/conversation", proxyHandler)
	}

	if commonConfig.GetConfig().HiddenUserInfo {
		// 为 /backend-api/me 设置处理器
		r.GET("/backend-api/me", middleware.CreateProxyHandler(commonConfig.GetConfig().OpenAiSite, middleware.ProcessBackendApiMeResponse))
	} else {
		// 对于 /backend-api/me 的请求，直接使用反向代理
		r.GET("/backend-api/me", proxyHandler)
	}

	// 处理所有请求
	r.Use(func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/backend-api/conversation" && c.Request.Method == "POST" {
			// 已经在上面处理过了，直接返回
			return
		}

		if path == "/backend-api/me" && c.Request.Method == "GET" {
			// 已经在上面处理过了，直接返回
			return
		}

		// 对于其他所有请求，使用反向代理
		proxyHandler(c)
	})

	s := openai.NewServer(
		r,
		logger,
		openai.WithServerHost(commonConfig.GetConfig().HttpHost),
		openai.WithServerPort(commonConfig.GetConfig().OpenAiPort),
	)

	return s
}

// NewClaudeReverseProxyServer 创建 Claude 反向代理服务器
func NewClaudeReverseProxyServer(
	logger *log.Logger,
	conversationLoggerMiddleware *middleware.ConversationLoggerMiddleware,
) *claude.Server {
	r := gin.Default()

	// 创建反向代理处理函数
	proxyHandler := reverseProxy(commonConfig.GetConfig().ClaudeSite)

	if commonConfig.GetConfig().ModerationEnable() {
		r.POST("/api/organizations/:id1/chat_conversations/:id2/completion", middleware.ClaudeContentModerationMiddleware(logger), proxyHandler)
	} else {
		r.POST("/api/organizations/:id1/chat_conversations/:id2/completion", proxyHandler)
	}

	// 处理所有请求
	r.Use(proxyHandler)

	s := claude.NewServer(
		r,
		logger,
		claude.WithServerHost(commonConfig.GetConfig().HttpHost),
		claude.WithServerPort(commonConfig.GetConfig().ClaudePort),
	)

	return s
}

// 创建反向代理处理函数
func reverseProxy(target string) gin.HandlerFunc {
	parse, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(parse)

	// 修改默认的Director函数
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// 保持原始的Host头
		req.Host = parse.Host
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
