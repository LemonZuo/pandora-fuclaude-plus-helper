package model

import (
	"time"
)

type ClaudeToken struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement" comment:"主键" column:"id"`
	TokenName    string    `json:"tokenName" gorm:"not null" comment:"token名称" column:"token_name"`
	SessionToken string    `json:"sessionToken" gorm:"not null" comment:"sessionToken" column:"session_token"`
	CreateTime   time.Time `json:"createTime" gorm:"not null" comment:"创建时间" column:"create_time"`
	UpdateTime   time.Time `json:"updateTime" gorm:"not null" comment:"更新时间" column:"update_time"`
}

func (m *ClaudeToken) TableName() string {
	return "tb_claude_token"
}
