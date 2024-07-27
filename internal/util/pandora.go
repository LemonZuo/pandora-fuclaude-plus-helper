package util

import (
	commonConfig "PandoraFuclaudePlusHelper/config"
	"PandoraFuclaudePlusHelper/pkg/log"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"strings"
)

// GenAccessToken generates an access token based on the refresh token
func GenAccessToken(refreshToken string, logger *log.Logger) (string, int, error) {
	var resp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"refresh_token": refreshToken,
		}).
		SetResult(&resp).
		Post(commonConfig.GetConfig().TokenUrl)
	if err != nil {
		logger.Error("GenAccessToken error", zap.Any("err", err))
		return "", -1, err
	}
	logger.Info("GenAccessToken resp", zap.Any("resp", resp))
	return resp.AccessToken, resp.ExpiresIn, nil
}

// CheckSubscriptionStatus GenShareToken generates a share token based on the access token
func CheckSubscriptionStatus(accessToken string, logger *log.Logger) int {
	if accessToken == "" {
		logger.Error("check_subscription_status: 1, because of empty access token")
		return 1
	}

	// 定义响应体的匿名结构
	var responseBody struct {
		Models []struct {
			Slug string `json:"slug"`
		} `json:"models"`
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+accessToken).
		SetResult(&responseBody).
		Get(commonConfig.GetConfig().CheckSubscribeUrl)

	if err != nil {
		logger.Error("check_subscription_status: 1, because of http error: ", zap.Any("err", err))
		return 1
	}

	if resp.StatusCode() != 200 && resp.StatusCode() != 401 {
		logger.Info("check_subscription_status: 1, because of status code", zap.Int("status_code", resp.StatusCode()))
		return 1
	}

	if resp.StatusCode() == 401 {
		logger.Info("check_subscription_status: 1, because of status code 401")
		return 2
	}

	subscribePlus := false
	for _, item := range responseBody.Models {
		if item.Slug == "gpt-4" {
			subscribePlus = true
			break
		}
	}

	if subscribePlus {
		logger.Info("check_subscription_status: 3, because of gpt-4 model found in response body")
		return 3
	} else {
		logger.Info("check_subscription_status: 2, because of no gpt-4 model found in response body")
		return 2
	}
}

// GenShareToken generates a share token based on the access token
func GenShareToken(accessToken string,
	uniqueName string,
	expiresIn int,
	gpt35Limit int,
	gpt4Limit int,
	showConversations bool,
	showUserinfo bool,
	resetLimit bool,
	temporaryChat bool,
	logger *log.Logger) (string, int64, error) {
	var resp struct {
		ExpireAt          int64  `json:"expire_at"`
		Gpt35Limit        int    `json:"gpt35_limit"`
		Gpt4Limit         int    `json:"gpt4_limit"`
		ResetLimit        bool   `json:"reset_limit"`
		ShowConversations bool   `json:"show_conversations"`
		TemporaryChat     bool   `json:"temporary_chat"`
		SiteLimit         string `json:"site_limit"`
		TokenKey          string `json:"token_key"`
		UniqueName        string `json:"unique_name"`
	}
	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"unique_name":        uniqueName,
			"access_token":       accessToken,
			"expires_in":         fmt.Sprintf("%d", expiresIn),
			"site_limit":         "",
			"reset_limit":        fmt.Sprintf("%t", resetLimit),
			"show_conversations": fmt.Sprintf("%t", showConversations),
			"temporary_chat":     fmt.Sprintf("%t", temporaryChat),
			"show_userinfo":      fmt.Sprintf("%t", showUserinfo),
			"gpt35_limit":        fmt.Sprintf("%d", gpt35Limit),
			"gpt4_limit":         fmt.Sprintf("%d", gpt4Limit),
		}).
		SetResult(&resp).
		Post(commonConfig.GetConfig().ShareTokenUrl)
	logger.Info("GenShareToken:", zap.Any("ShareTokenUrl", commonConfig.GetConfig().ShareTokenUrl), zap.Any("resp", resp))

	if err != nil {
		logger.Error("GenerateShareToken error", zap.Any("err", err))
		return "", -1, err
	}
	logger.Info("GenerateShareToken resp", zap.Any("resp", resp))
	return resp.TokenKey, resp.ExpireAt, nil
}

type ShareTokenInfo struct {
	Email      string                 `json:"email"`
	ExpireAt   int64                  `json:"expire_at"`
	Gpt35Limit string                 `json:"gpt35_limit"`
	Gpt4Limit  string                 `json:"gpt4_limit"`
	Usage      map[string]interface{} `json:"usage"`
	UserID     string                 `json:"user_id"`
}

// GetShareTokenInfo gets the share token information based on the share token and access token
func GetShareTokenInfo(shareToken string, accessToken string, logger *log.Logger) (ShareTokenInfo, error) {
	if shareToken == "" || accessToken == "" {
		logger.Error("shareToken or accessToken is empty")
		return ShareTokenInfo{}, errors.New("shareToken or accessToken is empty")

	}
	shareTokenInfoUrl := fmt.Sprintf("%s/%s", commonConfig.GetConfig().ShareTokenInfoUrl, shareToken)
	var resp ShareTokenInfo
	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		SetResult(&resp).
		Get(shareTokenInfoUrl)
	if err != nil {
		logger.Error("GetShareTokenInfo error", zap.Any("err", err))
		return ShareTokenInfo{}, err
	}
	logger.Info("GetShareTokenInfo resp", zap.Any("resp", resp))
	return resp, nil
}

// ExecuteShareAuth executes the share auth based on the share token
func ExecuteShareAuth(shareToken string, logger *log.Logger) (string, error) {
	if shareToken == "" {
		logger.Error("shareToken is empty")
		return "", errors.New("shareToken is empty")
	}
	var resp struct {
		LoginUrl   string `json:"login_url"`
		OauthToken string `json:"oauth_token"`
	}
	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Origin", commonConfig.GetConfig().ShareTokenAuth).
		SetBody(fmt.Sprintf(`{"share_token": "%s"}`, shareToken)).
		SetResult(&resp).
		Post(commonConfig.GetConfig().ShareTokenAuthUrl)
	if err != nil {
		logger.Error("ExecuteShareAuth error", zap.Any("err", err))
		return "", err
	}
	if resp.LoginUrl == "" {
		logger.Error("ExecuteShareAuth error", zap.Any("err", "login url is empty"))
		return "", errors.New("login url is empty")
	}
	logger.Info("ExecuteShareAuth resp", zap.Any("resp", resp))
	return resp.LoginUrl, nil
}

func ExecuteClaudeAuth(sessionToken string, accountName string, logger *log.Logger) (string, error) {
	var resp struct {
		LoginUrl string `json:"login_url"`
	}
	requestBody := ""
	if sessionToken == "" {
		requestBody = fmt.Sprintf(`{"session_key": "%s"}`, sessionToken)
	} else {
		requestBody = fmt.Sprintf(`{"session_key": "%s","unique_name":"%s", "expires_in": %d}`, sessionToken, accountName, 60*60*24*7)

	}
	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Origin", commonConfig.GetConfig().ShareTokenAuth).
		SetBody(requestBody).
		SetResult(&resp).
		Post(commonConfig.GetConfig().FuclaudeAuthUrl)
	if err != nil {
		logger.Error("ExecuteShareAuth error", zap.Any("err", err))
		return "", err
	}
	if resp.LoginUrl == "" {
		logger.Error("ExecuteShareAuth error", zap.Any("err", "login url is empty"))
		return "", errors.New("login url is empty")
	}
	logger.Info("ExecuteShareAuth resp", zap.Any("resp", resp))
	return fmt.Sprintf("%s%s", commonConfig.GetConfig().FuclaudeAuth, resp.LoginUrl), nil
}

func CheckShareTokenType(shareToken string) int {
	if shareToken == "" {
		return -1
	}
	if strings.HasPrefix(shareToken, "sk-") {
		return 2
	}
	return 1
}
