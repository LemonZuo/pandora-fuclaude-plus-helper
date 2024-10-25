package main

import (
	"PandoraFuclaudePlusHelper/cmd/server/wire"
	commonConfig "PandoraFuclaudePlusHelper/config"
	"PandoraFuclaudePlusHelper/pkg/log"
	"context"
	"fmt"
)

// @title           Nunu Example API
// @version         1.0.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8000
// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	commonConfig.InitConfig()

	logger := log.NewLog()

	app, cleanup, err := wire.NewWire(logger)

	defer cleanup()
	if err != nil {
		panic(err)
	}

	// 打印服务启动信息
	printEndpoint()

	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}

// printEndpoint 打印服务启动信息
func printEndpoint() {
	version := commonConfig.GetConfig().Version
	apiEndpoint := fmt.Sprintf("http://%s:%d", commonConfig.GetConfig().HttpHost, commonConfig.GetConfig().ApiPort)
	openAiEndpoint := fmt.Sprintf("http://%s:%d", commonConfig.GetConfig().HttpHost, commonConfig.GetConfig().OpenAiPort)
	claudeEndpoint := fmt.Sprintf("http://%s:%d", commonConfig.GetConfig().HttpHost, commonConfig.GetConfig().ClaudePort)
	fmt.Printf("PandoraFuclaudePlusHelper [%s] API started at %s\n", version, apiEndpoint)
	fmt.Printf("PandoraFuclaudePlusHelper [%s] OpenAI Reverse started at %s\n", version, openAiEndpoint)
	fmt.Printf("PandoraFuclaudePlusHelper [%s] Claude Reverse started at %s\n", version, claudeEndpoint)
}
