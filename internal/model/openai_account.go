package model

import (
	"time"
)

type OpenaiAccount struct {
	ID                int64     `json:"id" gorm:"primaryKey;autoIncrement" comment:"主键" column:"id"`
	UserId            int64     `json:"userId" gorm:"not null" comment:"token_id" column:"user_id"`
	TokenID           int64     `json:"tokenId" gorm:"not null" comment:"token_id" column:"token_id"`
	Account           string    `json:"account" gorm:"not null;unique" comment:"唯一名称" column:"account"`
	ExpirationTime    time.Time `json:"expirationTime" comment:"到期时间" column:"expiration_time"`
	Status            int       `json:"status" gorm:"not null;default:1" comment:"状态, 1:正常, 0:禁用" column:"status"`
	Gpt35Limit        int       `json:"gpt35Limit" gorm:"default:-1" comment:"GPT-3.5次数(为0无法使用，负数不限制)" column:"gpt35_limit"`
	Gpt4Limit         int       `json:"gpt4Limit" gorm:"default:-1" comment:"GPT-4.0次数(为0无法使用，负数不限制)" column:"gpt4_limit"`
	Gpt4oLimit        int       `json:"gpt4oLimit" gorm:"default:-1" comment:"GPT-4o次数(为0无法使用，负数不限制)" column:"gpt4o_limit"`
	Gpt4oMiniLimit    int       `json:"gpt4oMiniLimit" gorm:"default:-1" comment:"GPT-4o mini次数(为0无法使用，负数不限制)" column:"gpt4o_mini_limit"`
	O1Limit           int       `json:"o1Limit" gorm:"default:-1" comment:"o1次数(为0无法使用，负数不限制)" column:"o1_limit"`
	O1MiniLimit       int       `json:"o1MiniLimit" gorm:"default:-1" comment:"o1 mini次数(为0无法使用，负数不限制)" column:"o1_mini_limit"`
	ShowConversations int       `json:"showConversations" gorm:"default:0" comment:"会话无需隔离，1:不隔离,0:隔离" column:"show_conversations"`
	TemporaryChat     int       `json:"temporaryChat" gorm:"default:0" comment:"临时聊天，1:强制使用,0:非强制使用" column:"temporary_chat"`
	ShareToken        string    `json:"shareToken" gorm:"not null" comment:"共享token" column:"share_token"`
	ShareTokenEncrypt string    `json:"shareTokenEncrypt" gorm:"not null;default:0" comment:"加密共享token" column:"share_token_encrypt"`
	ExpireAt          time.Time `json:"expireAt" gorm:"not null" comment:"过期时间" column:"expire_at"`
	CreateTime        time.Time `json:"createTime" gorm:"not null" comment:"创建时间" column:"create_time"`
	UpdateTime        time.Time `json:"updateTime" gorm:"not null" comment:"更新时间" column:"update_time"`
}

func (m *OpenaiAccount) TableName() string {
	return "tb_openai_account"
}
