package v1

import (
	"PandoraFuclaudePlusHelper/internal/model"
)

type SearchClaudeAccountRequest struct {
	TokenId int64 `json:"tokenId"`
}

type AddClaudeAccountRequest struct {
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

type UpdateClaudeAccountRequest struct {
	Account model.ClaudeAccount `json:"account" binding:"required"`
}

type DeleteClaudeAccountRequest struct {
	Id int64 `json:"id" binding:"required"`
}

type DisableClaudeAccountRequest struct {
	Id int64 `json:"id" binding:"required"`
}

type EnableClaudeAccountRequest struct {
	Id int64 `json:"id" binding:"required"`
}

type SearchClaudeAccountResponse struct {
	Response
	Data []*model.ClaudeAccount `json:"data"`
}

type ClaudeAccountResponse struct {
	Response
}

type StatisticClaudeAccountRequest struct {
	TokenId int64 `json:"tokenId" binding:"required"`
}

type StatisticClaudeAccountResponseData struct {
	Categories []string                 `json:"categories"`
	Series     []map[string]interface{} `json:"series"`
}

type StatisticClaudeAccountResponse struct {
	Response
	Data StatisticClaudeAccountResponseData `json:"data"`
}
