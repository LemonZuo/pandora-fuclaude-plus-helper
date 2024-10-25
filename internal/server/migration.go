package server

import (
	"PandoraFuclaudePlusHelper/internal/model"
	"PandoraFuclaudePlusHelper/pkg/log"
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Migrate struct {
	db  *gorm.DB
	log *log.Logger
}

func NewMigrate(db *gorm.DB, log *log.Logger) *Migrate {
	return &Migrate{
		db:  db,
		log: log,
	}
}
func (m *Migrate) Start(ctx context.Context) error {
	if err := m.db.AutoMigrate(
		model.OpenaiAccount{},
		model.OpenaiToken{},
		model.User{},
		model.ClaudeToken{},
		model.ClaudeAccount{},
	); err != nil {
		m.log.Error("user migrate error", zap.Error(err))
		return err
	}
	// TODO 待处理
	// m.initClaudeAccount()
	m.log.Info("AutoMigrate success")
	return nil
}
func (m *Migrate) Stop(ctx context.Context) error {
	m.log.Info("AutoMigrate stop")
	return nil
}

func (m *Migrate) initClaudeAccount() {
	var accounts []*model.ClaudeAccount
	// 查出ClaudeAccount 过期时间为空的数据
	m.db.Where("expire_at is null").Find(&accounts)
	// 判断是否有数据
	if len(accounts) == 0 {
		return
	}

	// 遍历数据，处理
	for _, account := range accounts {
		userId := account.UserId
		// 查询关联的用户
		var user model.User
		m.db.Where("id = ?", userId).First(&user)
		// 判断用户是否存在
		if user.ID == 0 {
			continue
		}
		// 更新过期时间
		// account.ExpireAt = user.ExpirationTime
		m.db.Save(&account)
	}
}
