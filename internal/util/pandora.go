package util

import (
	commonConfig "PandoraFuclaudePlusHelper/config"
	"PandoraFuclaudePlusHelper/pkg/log"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

var genAccessTokenMutex sync.Mutex
var checkSubscriptionStatusMutex sync.Mutex
var genShareTokenMutex sync.Mutex
var getShareTokenInfoMutex sync.Mutex
var executeShareAuthMutex sync.Mutex
var executeClaudeAuthMutex sync.Mutex

// GenAccessToken generates an access token based on the refresh token
func GenAccessToken(refreshToken string, logger *log.Logger) (string, int, error) {
	genAccessTokenMutex.Lock()
	defer genAccessTokenMutex.Unlock()

	// 优先使用 Pandora 的刷新令牌生成访问令牌
	accessToken, expiresIn, err := GenAccessTokenPandora(refreshToken, logger)
	if err == nil {
		return accessToken, expiresIn, nil
	}
	// 如果使用 Pandora 的刷新令牌生成访问令牌失败，则使用官方的刷新令牌生成访问令牌
	accessToken, expiresIn, err = GenAccessTokenOfficial(refreshToken, logger)
	if err != nil {
		return "", -1, err
	} else {
		return accessToken, expiresIn, nil
	}
}

// GenAccessTokenPandora generates an access token based on the refresh token
func GenAccessTokenPandora(refreshToken string, logger *log.Logger) (string, int, error) {

	var resp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("User-Agent", fmt.Sprintf("pandora-plus-helper/%s", commonConfig.GetConfig().Version)).
		SetFormData(map[string]string{
			"refresh_token": refreshToken,
		}).
		SetResult(&resp).
		Post(commonConfig.GetConfig().TokenUrl)
	if err != nil {
		logger.Error(fmt.Sprintf("GenAccessToken by pandora error, response: %v, error: %v", response, err))
		return "", -1, err
	}
	logger.Info(fmt.Sprintf("GenAccessToken by pandora, StatusCode: %d, responseContent: %s", response.StatusCode(), string(response.Body())))

	if response.StatusCode() != http.StatusOK {
		logger.Error(fmt.Sprintf("GenAccessToken by pandora error, code: %d", response.StatusCode()))
		return "", -1, errors.New(fmt.Sprintf("GenAccessToken by pandora error, code: %d", response.StatusCode))

	}

	return resp.AccessToken, resp.ExpiresIn, nil
}

// GenAccessTokenOfficial generates an access token based on the refresh token
func GenAccessTokenOfficial(refreshToken string, logger *log.Logger) (string, int, error) {

	// 定义并初始化 RefreshRequest 结构体
	RefreshRequest := struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		RefreshToken string `json:"refresh_token"`
		RedirectURI  string `json:"redirect_uri"`
	}{
		GrantType:    "refresh_token",
		ClientID:     "pdlLIX2Y72MIl2rhLhTE9VV9bN905kBh",
		RefreshToken: refreshToken,
		RedirectURI:  "com.openai.chat://auth0.openai.com/ios/com.openai.chat/callback",
	}

	var resp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	client := resty.New()

	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36").
		SetBody(RefreshRequest).
		SetResult(&resp).
		Post("https://auth0.openai.com/oauth/token")
	if err != nil {
		logger.Error(fmt.Sprintf("GenAccessToken by official error, response: %v, error: %v", response, err))
		return "", -1, err
	}
	logger.Info(fmt.Sprintf("GenAccessToken by official, StatusCode: %d, responseContent: %s", response.StatusCode(), string(response.Body())))

	if response.StatusCode() != http.StatusOK {
		logger.Error(fmt.Sprintf("GenAccessToken by official error, code: %d", response.StatusCode()))
		return "", -1, errors.New(fmt.Sprintf("GenAccessToken by official error, code: %d", response.StatusCode))

	}

	return resp.AccessToken, resp.ExpiresIn, nil
}

// CheckSubscriptionStatus GenShareToken generates a share token based on the access token
func CheckSubscriptionStatus(accessToken string, logger *log.Logger) int {

	checkSubscriptionStatusMutex.Lock()
	defer checkSubscriptionStatusMutex.Unlock()

	if accessToken == "" {
		logger.Error("CheckSubscriptionStatus: 1, because of empty access token")
		return 1
	}

	client := resty.New()
	var responseBody Response

	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", fmt.Sprintf("pandora-plus-helper/%s", commonConfig.GetConfig().Version)).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		SetResult(&responseBody).
		Get("https://new.oaifree.com/backend-api/accounts/check/v4-2023-04-27?timezone_offset_min=-480")

	if err != nil {
		logger.Error(fmt.Sprintf("CheckSubscriptionStatus error, response: %v, error: %v", response, err))
		return 1
	}

	logger.Info(fmt.Sprintf("CheckSubscriptionStatus, StatusCode: %d, responseContent: %s", response.StatusCode(), string(response.Body())))

	if response.StatusCode() != 200 && response.StatusCode() != 401 {
		logger.Info(fmt.Sprintf("CheckSubscriptionStatus: 1, because of status code: %d", response.StatusCode()))
		return 1
	}

	if response.StatusCode() == 401 {
		logger.Info("CheckSubscriptionStatus: 2, because of status code 401")
		return 2
	}

	isSubscribe := false
	for _, item := range responseBody.Accounts {
		entitlement := item.Entitlement
		subscription := entitlement.HasActiveSubscription
		plan := entitlement.SubscriptionPlan
		expiresAt := entitlement.ExpiresAt
		if subscription {
			logger.Info(fmt.Sprintf("CheckSubscriptionStatus plan: %s, expiresAt: %s, subscription: %v", plan, expiresAt.Format(time.DateTime), subscription))
			isSubscribe = true
			break
		}
	}

	if isSubscribe {
		logger.Info("CheckSubscriptionStatus: 3, because of subscription is active")
		return 3
	} else {
		logger.Info("CheckSubscriptionStatus: 2, because of no active subscription")
		return 2
	}
}

// GenShareToken generates a share token based on the access token
func GenShareToken(accessToken string,
	uniqueName string,
	expiresIn int,
	gpt35Limit int,
	gpt4Limit int,
	gpt4oLimit int,
	gpt4oMiniLimit int,
	o1Limit int,
	o1MiniLimit int,
	showConversations bool,
	showUserinfo bool,
	resetLimit bool,
	temporaryChat bool,
	logger *log.Logger) (string, string, int64, error) {

	genShareTokenMutex.Lock()
	defer genShareTokenMutex.Unlock()

	var resp struct {
		ExpireAt          int64  `json:"expire_at"`
		Gpt35Limit        int    `json:"gpt35_limit"`
		Gpt4Limit         int    `json:"gpt4_limit"`
		Gpt4oLimit        int    `json:"gpt4o_limit"`
		Gpt4oMiniLimit    int    `json:"gpt4o_mini_limit"`
		O1Limit           int    `json:"o1_limit"`
		O1MiniLimit       int    `json:"o1_mini_limit"`
		ResetLimit        bool   `json:"reset_limit"`
		ShowConversations bool   `json:"show_conversations"`
		TemporaryChat     bool   `json:"temporary_chat"`
		SiteLimit         string `json:"site_limit"`
		TokenKey          string `json:"token_key"`
		UniqueName        string `json:"unique_name"`
	}
	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("User-Agent", fmt.Sprintf("pandora-plus-helper/%s", commonConfig.GetConfig().Version)).
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
			"gpt4o_limit":        fmt.Sprintf("%d", gpt4oLimit),
			"gpt4o_mini_limit":   fmt.Sprintf("%d", gpt4oMiniLimit),
			"o1_limit":           fmt.Sprintf("%d", o1Limit),
			"o1_mini_limit":      fmt.Sprintf("%d", o1MiniLimit),
		}).
		SetResult(&resp).
		Post(commonConfig.GetConfig().ShareTokenUrl)

	if err != nil {
		logger.Error("GenerateShareToken error", zap.Any("err", err))
		return "", "", -1, err
	}
	logger.Info(fmt.Sprintf("GenerateShareToken, StatusCode: %d, responseContent: %s", response.StatusCode(), string(response.Body())))

	code := response.StatusCode()
	if code != http.StatusOK {
		logger.Error(fmt.Sprintf("GenerateShareToken error, code: %d", code))
		return "", "", -1, errors.New(fmt.Sprintf("GenerateShareToken error, code: %d", code))
	}
	hash := sha1.New()
	hash.Write([]byte(resp.TokenKey))
	enc := hex.EncodeToString(hash.Sum(nil))
	return resp.TokenKey, enc, resp.ExpireAt, nil
}

type ShareTokenInfo struct {
	Email          string                 `json:"email"`
	ExpireAt       int64                  `json:"expire_at"`
	Gpt35Limit     interface{}            `json:"gpt35_limit,omitempty"`
	Gpt4Limit      interface{}            `json:"gpt4_limit,omitempty"`
	Gpt4oLimit     interface{}            `json:"gpt4o_limit,omitempty"`
	Gpt4oMiniLimit interface{}            `json:"gpt4o_mini_limit,omitempty"`
	O1Limit        interface{}            `json:"o1_limit,omitempty"`
	O1MiniLimit    interface{}            `json:"o1_mini_limit,omitempty"`
	Usage          map[string]interface{} `json:"usage,omitempty"`
	UserID         string                 `json:"user_id,omitempty"`
}

// GetShareTokenInfo gets the share token information based on the share token and access token
func GetShareTokenInfo(shareToken string, accessToken string, logger *log.Logger) (ShareTokenInfo, error) {

	getShareTokenInfoMutex.Lock()
	defer getShareTokenInfoMutex.Unlock()

	if shareToken == "" || accessToken == "" {
		logger.Error("GetShareTokenInfo or accessToken is empty")
		return ShareTokenInfo{}, errors.New("shareToken or accessToken is empty")

	}
	shareTokenInfoUrl := fmt.Sprintf("%s/%s", commonConfig.GetConfig().ShareTokenInfoUrl, shareToken)
	var resp ShareTokenInfo
	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", fmt.Sprintf("pandora-plus-helper/%s", commonConfig.GetConfig().Version)).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		SetResult(&resp).
		Get(shareTokenInfoUrl)
	if err != nil {
		logger.Error("GetShareTokenInfo error", zap.Any("err", err))
		return ShareTokenInfo{}, err
	}

	logger.Info(fmt.Sprintf("GetShareTokenInfo by pandora, StatusCode: %d, responseContent: %s", response.StatusCode(), string(response.Body())))

	if response.StatusCode() != http.StatusOK {
		logger.Error(fmt.Sprintf("GetShareTokenInfo error, code: %d", response.StatusCode()))
		return ShareTokenInfo{}, errors.New(fmt.Sprintf("GetShareTokenInfo error, code: %d", response.StatusCode()))
	}

	return resp, nil
}

// ExecuteShareAuth executes the share auth based on the share token
func ExecuteShareAuth(shareToken string, logger *log.Logger) (string, error) {

	executeShareAuthMutex.Lock()
	defer executeShareAuthMutex.Unlock()

	if shareToken == "" {
		logger.Error("ExecuteShareAuth shareToken is empty")
		return "", errors.New("shareToken is empty")
	}
	var resp struct {
		LoginUrl   string `json:"login_url"`
		OauthToken string `json:"oauth_token"`
	}
	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", fmt.Sprintf("pandora-plus-helper/%s", commonConfig.GetConfig().Version)).
		SetHeader("Origin", commonConfig.GetConfig().OpenAiAuthSite).
		SetBody(fmt.Sprintf(`{"share_token": "%s"}`, shareToken)).
		SetResult(&resp).
		Post(fmt.Sprintf("%s/api/auth/oauth_token", commonConfig.GetConfig().OpenAiAuthSite))
	if err != nil {
		logger.Error("ExecuteShareAuth error", zap.Any("err", err))
		return "", err
	}

	logger.Info(fmt.Sprintf("ExecuteShareAuth, StatusCode: %d, responseContent: %s", response.StatusCode(), string(response.Body())))

	if response.StatusCode() != http.StatusOK {
		logger.Error(fmt.Sprintf("ExecuteShareAuth error, code: %d", response.StatusCode()))
		return "", errors.New(fmt.Sprintf("ExecuteShareAuth error, code: %d", response.StatusCode))
	}

	if resp.LoginUrl == "" {
		logger.Error("ExecuteShareAuth error, login url is empty")
		return "", errors.New("login url is empty")
	}
	return resp.LoginUrl, nil
}

func ExecuteClaudeAuth(sessionToken string, accountName string, seconds int, logger *log.Logger) (string, error) {

	executeClaudeAuthMutex.Lock()
	defer executeClaudeAuthMutex.Unlock()

	var resp struct {
		ExpiresAt  int64  `json:"expires_at"`
		LoginUrl   string `json:"login_url"`
		OauthToken string `json:"oauth_token"`
	}
	requestBody := ""
	if accountName == "" {
		requestBody = fmt.Sprintf(`{"session_key": "%s"}`, sessionToken)
	} else if seconds == -1 {
		requestBody = fmt.Sprintf(`{"session_key": "%s","unique_name":"%s"}`, sessionToken, accountName)
	} else {
		requestBody = fmt.Sprintf(`{"session_key": "%s","unique_name":"%s", "expires_in": %d}`, sessionToken, accountName, seconds)
	}

	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", fmt.Sprintf("pandora-plus-helper/%s", commonConfig.GetConfig().Version)).
		SetBody(requestBody).
		SetResult(&resp).
		Post(fmt.Sprintf("%s/manage-api/auth/oauth_token", commonConfig.GetConfig().ClaudeAuthSite))
	if err != nil {
		logger.Error(fmt.Sprintf("ExecuteClaudeAuth error, response: %v, error: %v", resp, err))
		return "", err
	}

	logger.Info(fmt.Sprintf("ExecuteClaudeAuth, StatusCode: %d, responseContent: %s", response.StatusCode(), string(response.Body())))

	if response.StatusCode() != http.StatusOK {
		logger.Error(fmt.Sprintf("ExecuteClaudeAuth error, code: %d", response.StatusCode()))
		return "", errors.New(fmt.Sprintf("ExecuteClaudeAuth error, code: %d", response.StatusCode()))
	}

	if resp.LoginUrl == "" {
		logger.Error("ExecuteClaudeAuth error, login url is empty")
		return "", errors.New("login url is empty")
	}
	return fmt.Sprintf("%s%s", commonConfig.GetConfig().ClaudeAuthSite, resp.LoginUrl), nil
}
