package middleware

import (
	commonConfig "PandoraFuclaudePlusHelper/config"
	"PandoraFuclaudePlusHelper/pkg/log"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Author struct {
		Role string `json:"role"`
	} `json:"author"`
	Content struct {
		ContentType string `json:"content_type"`
		Parts       []Part `json:"parts"`
	} `json:"content"`
}

type Part struct {
	StringValue *string            `json:"-"`
	ImageValue  *ImageAssetPointer `json:"-"`
}

type ImageAssetPointer struct {
	ContentType  string `json:"content_type"`
	AssetPointer string `json:"asset_pointer"`
	SizeBytes    int    `json:"size_bytes"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

func (p *Part) UnmarshalJSON(data []byte) error {
	// 尝试解析为字符串
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		p.StringValue = &str
		return nil
	}

	// 尝试解析为 ImageAssetPointer 对象
	var img ImageAssetPointer
	if err := json.Unmarshal(data, &img); err == nil {
		p.ImageValue = &img
		return nil
	}

	// 如果都无法解析，返回错误
	return fmt.Errorf("无法解析 Part: %s", string(data))
}

type ChatGPTConversationRequest struct {
	Messages []Message `json:"messages"`
}

type ClaudeConversationRequest struct {
	Prompt string
}

func OpenAiContentModerationMiddleware(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/backend-api/conversation" {
			// 读取cookie值
			userId, err := c.Request.Cookie("_Secure-next-auth.user-id")
			if err != nil || userId == nil || userId.Value == "" {
				logger.Info(fmt.Sprintf("Missing user ID: %v", userId))
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"detail": gin.H{
						"message": "Missing user ID",
						"flagged": true,
					},
				})
				return
			}
			shareToken, err := c.Request.Cookie("_Secure-next-auth.share-token")
			if err != nil || shareToken == nil || shareToken.Value == "" {
				logger.Info(fmt.Sprintf("Missing share token: %v", shareToken))
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"detail": gin.H{
						"message": "Missing share token",
						"flagged": true,
					},
				})
				return
			}

			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				logger.Info("Failed to read request body")
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"detail": gin.H{
						"message": "Failed to read request body",
						"flagged": true,
					},
				})
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			var requestBody ChatGPTConversationRequest
			if err := json.Unmarshal(body, &requestBody); err != nil {
				logger.Info("Failed to unmarshal request body")
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"detail": gin.H{
						"message": "Failed to unmarshal request body",
						"flagged": true,
					},
				})
				return
			}

			var userMessages []string
			for _, msg := range requestBody.Messages {
				if msg.Author.Role == "user" && (msg.Content.ContentType == "text" || msg.Content.ContentType == "multimodal_text") {
					for _, part := range msg.Content.Parts {
						if part.StringValue != nil {
							userMessages = append(userMessages, *part.StringValue)
						}
					}
				}
			}

			if len(userMessages) > 0 {
				shouldBlock, err := checkContentForModeration(userMessages, logger)
				if err != nil {
					logger.Info(fmt.Sprintf("Failed to check content for moderation: %v", err))
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"detail": gin.H{
							"message": "Failed to check content for moderation",
							"flagged": true,
						},
					})
					return
				}
				if shouldBlock {
					// 异步记录被阻止的消息到日志文件
					go asyncModerationLog(userId.Value, shareToken.Value, userMessages, logger)
					logger.Info(fmt.Sprintf("User %s with share token %s sent a message that was blocked by the moderation system, message: %v", userId.Value, shareToken.Value, userMessages))
					c.AbortWithStatusJSON(http.StatusUnavailableForLegalReasons, gin.H{
						"detail": gin.H{
							"message": commonConfig.GetConfig().ModerationMessage,
							"flagged": true,
						},
					})
					return
				}
			}
		}

		c.Next()
	}
}

func ClaudeContentModerationMiddleware(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		re := regexp.MustCompile(`^/api/organizations/([^/]+)/chat_conversations/([^/]+)/completion$`)
		// Claude道德检查
		if re.MatchString(c.Request.URL.Path) {
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				logger.Error("Failed to read request body")
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"type": "error",
					"error": gin.H{
						"type":    "bad_request",
						"message": "Failed to read request body",
					},
				})
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			var requestBody ClaudeConversationRequest
			if err := json.Unmarshal(body, &requestBody); err != nil {
				logger.Error("Failed to unmarshal request body")
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"type": "error",
					"error": gin.H{
						"type":    "bad_request",
						"message": "Failed to unmarshal request body",
					},
				})
				return
			}

			var userMessages []string
			userMessages = append(userMessages, requestBody.Prompt)

			if len(userMessages) > 0 {
				shouldBlock, err := checkContentForModeration(userMessages, logger)
				if err != nil {
					logger.Error(fmt.Sprintf("Failed to check content for moderation: %v", err))
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"type": "error",
						"error": gin.H{
							"type":    "internal_server_error",
							"message": "Failed to check content for moderation",
						},
					})
					return
				}
				if shouldBlock {
					logger.Info(fmt.Sprintf("User sent a message that was blocked by the moderation system, message: %v", userMessages))
					go asyncModerationLog("claude", "claude", userMessages, logger)
					c.AbortWithStatusJSON(http.StatusUnavailableForLegalReasons, gin.H{
						"type": "error",
						"error": gin.H{
							"type":    "moderation_error",
							"message": commonConfig.GetConfig().ModerationMessage,
						},
					})
					return
				}
			}
		}
		c.Next()
	}
}

func checkContentForModeration(messages []string, logger *log.Logger) (bool, error) {
	client := resty.New()
	client.SetTimeout(time.Second * 10)

	userMessage := strings.Join(messages, " ")
	if len(userMessage) == 0 {
		return false, nil
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", commonConfig.GetConfig().ModerationApiKey)).
		SetBody(map[string]interface{}{
			"input": userMessage,
		}).
		Post(fmt.Sprintf("%s/v1/moderations", commonConfig.GetConfig().ModerationEndpoint))

	if err != nil {
		return false, err
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Error("Moderation API returned an error: {}", zap.Any("body", resp.Body()))
		return true, fmt.Errorf("moderation API returned an error: %d", resp.StatusCode())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		logger.Error("Failed to unmarshal moderation response: {}", zap.Error(err))
		return true, err
	}

	results, ok := result["results"].([]interface{})
	if !ok {
		return true, fmt.Errorf("unexpected response format")
	}

	for _, r := range results {
		if flagged, ok := r.(map[string]interface{})["flagged"].(bool); ok && flagged {
			return true, nil
		}
	}

	return false, nil
}

// asyncModerationLog 异步记录被阻止的消息到日志文件
func asyncModerationLog(userId string, shareToken string, messages []string, logger *log.Logger) {
	logMessage := fmt.Sprintf("%s | %s | %s | %s\n",
		time.Now().Format(time.DateTime),
		userId,
		shareToken,
		strings.Join(messages, ", "))

	moderationLog := fmt.Sprintf("%s/%s", commonConfig.GetConfig().DataDir, "logs/moderation.log")

	file, err := os.OpenFile(moderationLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to open log file: %v", err))
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to close log file: %v", err))
		}
	}(file)

	if _, err := file.WriteString(logMessage); err != nil {
		logger.Error(fmt.Sprintf("Failed to write to log file: %v", err))
	}
}
