package model

import (
	"time"
)

type User struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement" comment:"主键" column:"id"`
	UniqueName     string    `json:"uniqueName" gorm:"not null;unique" comment:"唯一名称" column:"unique_name"`
	Password       string    `json:"password" gorm:"not null;unique" comment:"密码" column:"password"`
	Enable         int       `json:"enable" gorm:"default:1" comment:"是否启用, 0:禁用, 1:启用" column:"enable"`
	Openai         int       `json:"openai" gorm:"default:0" comment:"是否开启openai, 0:禁用, 1:启用" column:"openai"`
	OpenaiToken    int64     `json:"openaiToken" gorm:"default:0" comment:"OpenaiToken ID" column:"openai_token"`
	Claude         int       `json:"claude" gorm:"default:0" comment:"是否开启claude, 0:禁用, 1:启用" column:"claude"`
	ClaudeToken    int64     `json:"claudeToken" gorm:"default:0" comment:"ClaudeToken ID" column:"claude_token"`
	ExpirationTime time.Time `json:"expirationTime" gorm:"not null" comment:"过期时间" column:"expiration_time"`
	CreateTime     time.Time `json:"createTime" gorm:"not null" comment:"创建时间" column:"create_time"`
	UpdateTime     time.Time `json:"updateTime" gorm:"not null" comment:"更新时间" column:"update_time"`
}

func (m *User) TableName() string {
	return "tb_user"
}
