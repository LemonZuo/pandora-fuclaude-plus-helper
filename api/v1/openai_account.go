package v1

import (
	"PandoraFuclaudePlusHelper/internal/model"
)

type SearchOpenaiAccountRequest struct {
	TokenId int64 `json:"tokenId"`
}

type AddOpenaiAccountRequest struct {
	ID                int64  `json:"id"`
	Account           string `json:"account"`
	ExpirationTime    string `json:"expirationTime"`
	Status            int    `json:"status"`
	Gpt35Limit        int    `json:"gpt35Limit"`
	Gpt4Limit         int    `json:"gpt4Limit"`
	ShowConversations int    `json:"showConversations"`
	TemporaryChat     int    `json:"temporaryChat"`
	TokenID           int64  `json:"tokenId"`
	ShareToken        string `json:"shareToken"`
}

type UpdateOpenaiAccountRequest struct {
	Account model.OpenaiAccount `json:"account" binding:"required"`
}

type DeleteOpenaiAccountRequest struct {
	Id int64 `json:"id" binding:"required"`
}

type DisableOpenaiAccountRequest struct {
	Id int64 `json:"id" binding:"required"`
}

type EnableOpenaiAccountRequest struct {
	Id int64 `json:"id" binding:"required"`
}

type SearchOpenaiAccountResponse struct {
	Response
	Data []*model.OpenaiAccount `json:"data"`
}

type OpenaiAccountResponse struct {
	Response
}

type StatisticOpenaiAccountRequest struct {
	TokenId int64 `json:"tokenId" binding:"required"`
}

type StatisticOpenaiAccountResponseData struct {
	Categories []string                 `json:"categories"`
	Series     []map[string]interface{} `json:"series"`
}

type StatisticOpenaiAccountResponse struct {
	Response
	Data StatisticOpenaiAccountResponseData `json:"data"`
}
