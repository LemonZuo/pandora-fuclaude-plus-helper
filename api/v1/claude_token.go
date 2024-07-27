package v1

import "PandoraFuclaudePlusHelper/internal/model"

type SearchClaudeTokenRequest struct {
	TokenName string `json:"tokenName"`
}

type AddClaudeTokenRequest struct {
	// Token的所有属性
	model.ClaudeToken `json:"token" binding:"required"`
}

type UpdateClaudeTokenRequest struct {
	ClaudeToken model.ClaudeToken `json:"token" binding:"required"`
}

type DeleteClaudeTokenRequest struct {
	Id int64 `json:"id" binding:"required"`
}

type RefreshClaudeTokenRequest struct {
	Id int64 `json:"id" binding:"required"`
}

type SearchClaudeTokenResponseData struct {
	Response
	Data []*model.ClaudeToken `json:"data"`
}

type ClaudeTokenResponseData struct {
	Response
}
