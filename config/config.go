package config

import (
	"fmt"
	"github.com/sethvargo/go-password/password"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DataDir            string
	AdminPassword      string
	ApiKey             string
	TokenUrl           string
	ShareTokenUrl      string
	ShareTokenInfoUrl  string
	CheckSubscribeUrl  string
	OpenAiSite         string
	OpenAiAuthSite     string
	ClaudeSite         string
	ClaudeAuthSite     string
	ModerationEndpoint string
	ModerationApiKey   string
	ModerationMessage  string
	HiddenUserInfo     bool
	EnableTask         bool
	LogFileName        string
	LogLevel           string
	LogMaxSize         int
	LogMaxBackups      int
	LogMaxAge          int
	LogCompress        bool
	LogEncoding        string
	Env                string
	DatabaseDriver     string
	DatabaseDsn        string
	AppKey             string
	AppSecurity        string
	HttpHost           string
	ApiPort            int
	OpenAiPort         int
	ClaudePort         int
	StartTime          time.Time
	Version            string
	Secret             string
}

func (config *Config) ModerationEnable() bool {
	return config.ModerationEndpoint != "" && config.ModerationApiKey != ""
}

var globalConfig *Config
var Version = "0.0.0"
var initMutex sync.Mutex

// InitConfig 初始化全局配置
func InitConfig() {
	initMutex.Lock()
	defer initMutex.Unlock()

	// 加载.env文件
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("No .env file found")
	}

	dataDir := getEnvStr("DATA_DIR", "/data")
	driver, dsn := getDbConfig(dataDir)
	logFileName := fmt.Sprintf("%s/%s", dataDir, getEnvStr("LOG_FILE_NAME", "logs/server.log"))
	apiKey := getEnvStr("API_KEY", "dad04481-fa3f-494e-b90c-b822128073e5")
	defaultCheckSubscribeUrl := fmt.Sprintf("https://chat.oaifree.com/%s/backend-api/models?history_and_training_disabled=false", apiKey)

	globalConfig = &Config{
		DataDir:            dataDir,
		AdminPassword:      getAdminPassword(),
		ApiKey:             apiKey,
		TokenUrl:           getEnvStr("TOKEN_URL", "https://token.oaifree.com/api/auth/refresh"),
		ShareTokenUrl:      getEnvStr("SHARE_TOKEN_URL", "https://chat.oaifree.com/token/register"),
		ShareTokenInfoUrl:  getEnvStr("SHARE_TOKEN_INFO_URL", "https://chat.oaifree.com/token/info"),
		CheckSubscribeUrl:  getEnvStr("CHECK_SUBSCRIBE_URL", defaultCheckSubscribeUrl),
		OpenAiSite:         getEnvStr("OPENAI_SITE", "https://new.oaifree.com"),
		OpenAiAuthSite:     getEnvStr("OPENAI_AUTH_SITE|SHARE_TOKEN_AUTH", "https://new.oaifree.com"),
		ClaudeSite:         getEnvStr("CLAUDE_SITE", "https://demo.fuclaude.com"),
		ClaudeAuthSite:     getEnvStr("CLAUDE_AUTH_SITE|FUCLAUDE_LOGIN_AUTH", "https://demo.fuclaude.com"),
		ModerationEndpoint: getEnvStr("MODERATION_ENDPOINT", "https://api.openai.com"),
		ModerationApiKey:   getEnvStr("MODERATION_API_KEY", ""),
		ModerationMessage:  getEnvStr("MODERATION_MESSAGE", "Your message has been blocked due to inappropriate content"),
		HiddenUserInfo:     getEnvBool("HIDDEN_USER_INFO", false),
		EnableTask:         getEnvBool("ENABLE_TASK", true),
		LogFileName:        logFileName,
		LogLevel:           getEnvStr("LOG_LEVEL", "info"),
		LogMaxSize:         getEnvInt("LOG_MAX_SIZE", 10),
		LogMaxBackups:      getEnvInt("LOG_MAX_BACKUPS", 15),
		LogMaxAge:          getEnvInt("LOG_MAX_AGE", 30),
		LogCompress:        getEnvBool("LOG_COMPRESS", true),
		LogEncoding:        getEnvStr("LOG_ENCODING", "console"),
		Env:                getEnvStr("ENV", "dev"),
		DatabaseDriver:     driver,
		DatabaseDsn:        dsn,
		AppKey:             getEnvStr("APP_KEY", ""),
		AppSecurity:        getEnvStr("APP_SECURITY", ""),
		HttpHost:           getEnvStr("HTTP_HOST", "0.0.0.0"),
		ApiPort:            getEnvInt("HTTP_PORT|API_PORT", 5000),
		OpenAiPort:         getEnvInt("OPENAI_PORT", 5001),
		ClaudePort:         getEnvInt("CLAUDE_PORT", 5002),
		StartTime:          time.Now(),
		Version:            getVersion(),
		Secret:             getSecret(),
	}
}

// getDbConfig 获取数据库配置
func getDbConfig(dataDir string) (string, string) {
	driver := getEnvStr("DATABASE_DRIVER", "sqlite")
	dsn := ""
	if driver == "sqlite" {
		dsn = getEnvStr("DATABASE_DSN", "pandora-plus-helper.db")
		dsn = fmt.Sprintf("%s/%s", dataDir, dsn)
	} else if driver == "mysql" {
		dsn = getEnvStr("DATABASE_DSN", "")
		if len(dsn) == 0 {
			// 未配置 mysql dsn,退出程序
			fmt.Println("Database dsn not configured")
			os.Exit(-1)
		}
	} else {
		fmt.Println("Database driver not supported")
		os.Exit(-1)
	}
	return driver, dsn
}

// getAdminPassword 获取管理员密码
func getAdminPassword() string {
	adminPassword := getEnvStr("ADMIN_PASSWORD", "")
	if len(adminPassword) == 0 {
		var err error
		// 生成至少12位的密码
		adminPassword, err = password.Generate(12, 1, 1, false, false)
		if err != nil {
			// 如果生成密码出错，则退出程序
			fmt.Print(fmt.Sprintf("generated admin password error: %v", err))
			os.Exit(-1)
		}
		fmt.Println("Generated admin password:", adminPassword)
	} else if len(adminPassword) < 8 {
		fmt.Println("admin password length must be greater than 8")
		os.Exit(-1)
	}

	return adminPassword
}

// getSecret 获取密钥
func getSecret() string {
	secret := getEnvStr("SECRET", "")
	if len(secret) > 0 {
		return secret
	}
	secret, err := password.Generate(32, 1, 1, false, false)
	if err != nil {
		fmt.Print(fmt.Sprintf("generated secret error: %v", err))
		os.Exit(-1)
	}
	return secret
}

// getVersion 返回应用程序的版本号。
func getVersion() string {
	// 首先检查 Version 变量是否已被设置（即它是否不是默认值）。
	if Version != "0.0.0" {
		return Version
	}
	// 如果没有，则尝试读取环境变量 VERSION，如果都未设置，则默认返回 "0.0.0"。
	return getEnvStr("VERSION", "0.0.0")
}

// getEnvStr 返回第一个存在的环境变量的值，如果都不存在，则返回 defaultValue。
func getEnvStr(key, defaultValue string) string {
	// 用 "|" 分割 key 字符串，处理多个环境变量名。
	keys := strings.Split(key, "|")
	// 遍历所有的键名。
	for _, k := range keys {
		// 检查环境变量是否存在。
		if value, exists := os.LookupEnv(k); exists {
			// 如果找到，返回环境变量的值。
			return value
		}
	}
	// 如果所有环境变量都不存在，返回默认值。
	return defaultValue
}

// getEnvInt 返回第一个存在的环境变量的整数值，如果都不存在或转换失败，则返回默认值。
func getEnvInt(key string, defaultValue int) int {
	// 用 "|" 分割 key 字符串，支持多个环境变量名。
	keys := strings.Split(key, "|")
	// 遍历所有的键名。
	for _, k := range keys {
		// 检查环境变量是否存在。
		if value, exists := os.LookupEnv(k); exists {
			// 尝试将字符串转换为整数。
			intValue, err := strconv.Atoi(value)
			if err == nil {
				// 如果转换成功，返回整数值。
				return intValue
			}
		}
	}
	// 如果所有环境变量都不存在或转换失败，返回默认值。
	return defaultValue
}

// getEnvBool 返回第一个存在的环境变量的布尔值，如果都不存在或转换失败，则返回默认值。
func getEnvBool(key string, defaultValue bool) bool {
	// 用 "|" 分割 key 字符串，支持多个环境变量名。
	keys := strings.Split(key, "|")
	// 遍历所有的键名。
	for _, k := range keys {
		// 检查环境变量是否存在。
		if value, exists := os.LookupEnv(k); exists {
			// 尝试将字符串转换为布尔值。
			boolValue, err := strconv.ParseBool(value)
			if err == nil {
				// 如果转换成功，返回布尔值。
				return boolValue
			}
		}
	}
	// 如果所有环境变量都不存在或转换失败，返回默认值。
	return defaultValue
}

// GetConfig 提供全局配置的访问
func GetConfig() *Config {
	if globalConfig == nil {
		fmt.Printf("Config is not initialized")
		os.Exit(-1)
	}
	return globalConfig
}
