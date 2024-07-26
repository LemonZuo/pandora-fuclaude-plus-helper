package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DataDir           string
	AdminPassword     string
	ApiKey            string
	ShareTokenAuth    string
	ShareTokenAuthUrl string
	FuclaudeAuth      string
	FuclaudeAuthUrl   string
	TokenUrl          string
	ShareTokenUrl     string
	ShareTokenInfoUrl string
	CheckSubscribeUrl string
	LogFileName       string
	LogLevel          string
	LogMaxSize        int
	LogMaxBackups     int
	LogMaxAge         int
	LogCompress       bool
	LogEncoding       string
	Env               string
	DatabaseDriver    string
	DatabaseDsn       string
	AppKey            string
	AppSecurity       string
	HttpHost          string
	HttpPort          int
	StartTime         time.Time
}

var globalConfig *Config

// InitConfig 初始化全局配置
func InitConfig() {
	// 加载.env文件
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	dataDir := getEnvStr("DATA_DIR", "/data")
	apiKey := getEnvStr("API_KEY", "dad04481-fa3f-494e-b90c-b822128073e5")
	shareTokenAuth := getEnvStr("SHARE_TOKEN_AUTH", "https://new.oaifree.com")
	shareTokenAuthUrl := fmt.Sprintf("%s/api/auth/oauth_token", shareTokenAuth)
	defaultCheckSubscribeUrl := fmt.Sprintf("https://chat.oaifree.com/%s/backend-api/models?history_and_training_disabled=false", apiKey)
	logFileName := fmt.Sprintf("%s/%s", dataDir, getEnvStr("LOG_FILE_NAME", "logs/server.log"))
	databaseDriver := getEnvStr("DATABASE_DRIVER", "sqlite")
	databaseDsn := ""
	if databaseDriver == "sqlite" {
		dsn := getEnvStr("DATABASE_DSN", "pandora-plus-helper.db")
		databaseDsn = fmt.Sprintf("%s/%s", dataDir, dsn)
	} else if databaseDriver == "mysql" {
		dsn := getEnvStr("DATABASE_DSN", "")
		if len(dsn) == 0 {
			// 未配置 mysql dsn,退出程序
			fmt.Println("Database dsn not configured")
			os.Exit(-1)
		}
		databaseDsn = dsn
	} else {
		fmt.Println("Database driver not supported")
		os.Exit(-1)
	}
	fuclaudeAuth := getEnvStr("FUCLAUDE_LOGIN_AUTH", "https://demo.fuclaude.com")
	fuclaudeAuthUrl := fmt.Sprintf("%s/manage-api/auth/oauth_token", fuclaudeAuth)
	globalConfig = &Config{
		DataDir:           dataDir,
		AdminPassword:     getEnvStr("ADMIN_PASSWORD", ""),
		ApiKey:            apiKey,
		ShareTokenAuth:    shareTokenAuth,
		ShareTokenAuthUrl: shareTokenAuthUrl,
		FuclaudeAuth:      fuclaudeAuth,
		FuclaudeAuthUrl:   fuclaudeAuthUrl,
		TokenUrl:          getEnvStr("TOKEN_URL", "https://token.oaifree.com/api/auth/refresh"),
		ShareTokenUrl:     getEnvStr("SHARE_TOKEN_URL", "https://chat.oaifree.com/token/register"),
		ShareTokenInfoUrl: getEnvStr("SHARE_TOKEN_INFO_URL", "https://chat.oaifree.com/token/info"),
		CheckSubscribeUrl: getEnvStr("CHECK_SUBSCRIBE_URL", defaultCheckSubscribeUrl),
		LogFileName:       logFileName,
		LogLevel:          getEnvStr("LOG_LEVEL", "info"),
		LogMaxSize:        getEnvInt("LOG_MAX_SIZE", 1024),
		LogMaxBackups:     getEnvInt("LOG_MAX_BACKUPS", 30),
		LogMaxAge:         getEnvInt("LOG_MAX_AGE", 7),
		LogCompress:       getEnvBool("LOG_COMPRESS", true),
		LogEncoding:       getEnvStr("LOG_ENCODING", "console"),
		Env:               getEnvStr("ENV", "dev"),
		DatabaseDriver:    databaseDriver,
		DatabaseDsn:       databaseDsn,
		AppKey:            getEnvStr("APP_KEY", ""),
		AppSecurity:       getEnvStr("APP_SECURITY", ""),
		HttpHost:          getEnvStr("HTTP_HOST", "0.0.0.0"),
		HttpPort:          getEnvInt("HTTP_PORT", 5000),
		StartTime:         time.Now(),
	}
}

// getEnvStr 读取环境变量或返回默认值
func getEnvStr(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt 读取环境变量或返回默认值
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool 读取环境变量或返回默认值
func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// GetConfig 提供全局配置的访问
func GetConfig() *Config {
	if globalConfig == nil {
		log.Fatal("Config is not initialized")
	}
	return globalConfig
}
