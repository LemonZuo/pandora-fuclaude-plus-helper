package model

import (
	"time"
)

type ClaudeAccount struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement" comment:"主键" column:"id"`
	UserId     int64     `json:"userId" gorm:"not null" comment:"token_id" column:"user_id"`
	TokenID    int64     `json:"tokenId" gorm:"not null" comment:"token_id" column:"token_id"`
	Account    string    `json:"account" gorm:"not null;unique" comment:"唯一名称" column:"account"`
	Status     int       `json:"status" gorm:"not null;default:1" comment:"状态, 1:正常, 0:禁用" column:"status"`
	CreateTime time.Time `json:"createTime" gorm:"not null" comment:"创建时间" column:"create_time"`
	UpdateTime time.Time `json:"updateTime" gorm:"not null" comment:"更新时间" column:"update_time"`
	// ExpireAt   time.Time `json:"expireAt" gorm:"not null" comment:"过期时间" column:"expire_at"`
}

func (m *ClaudeAccount) TableName() string {
	return "tb_claude_account"
}
