package v1

import "PandoraFuclaudePlusHelper/internal/model"

type SearchOpenaiTokenRequest struct {
	TokenName string `json:"tokenName"`
}

type AddOpenaiTokenRequest struct {
	// Token的所有属性
	model.OpenaiToken `json:"token" binding:"required"`
}

type UpdateOpenaiTokenRequest struct {
	OpenaiToken model.OpenaiToken `json:"token" binding:"required"`
}

type DeleteOpenaiTokenRequest struct {
	Id int64 `json:"id" binding:"required"`
}

type RefreshOpenaiTokenRequest struct {
	Id int64 `json:"id" binding:"required"`
}

type SearchOpenaiTokenResponseData struct {
	Response
	Data []*model.OpenaiToken `json:"data"`
}

type OpenaiTokenResponseData struct {
	Response
}
