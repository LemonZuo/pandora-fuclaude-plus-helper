package model

import (
	"time"
)

type OpenaiToken struct {
	ID               int64     `json:"id" gorm:"primaryKey;autoIncrement" comment:"主键" column:"id"`
	TokenName        string    `json:"tokenName" gorm:"not null" comment:"token名称" column:"token_name"`
	PlusSubscription int       `json:"plusSubscription" gorm:"default:0" comment:"订阅状态, 0:未知, 1:未订阅, 2:已订阅" column:"plus_subscription"`
	RefreshToken     string    `json:"refreshToken" gorm:"not null;unique" comment:"刷新token" column:"refresh_token"`
	AccessToken      string    `json:"accessToken" gorm:"not null" comment:"访问token" column:"access_token"`
	ExpireAt         time.Time `json:"expireAt" gorm:"not null" comment:"过期时间" column:"expire_at"`
	CreateTime       time.Time `json:"createTime" gorm:"not null" comment:"创建时间" column:"create_time"`
	UpdateTime       time.Time `json:"updateTime" gorm:"not null" comment:"更新时间" column:"update_time"`
}

func (m *OpenaiToken) TableName() string {
	return "tb_openai_token"
}
