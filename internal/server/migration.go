package server

import (
	"PandoraPlusHelper/internal/model"
	"PandoraPlusHelper/pkg/log"
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
	m.log.Info("AutoMigrate success")
	return nil
}
func (m *Migrate) Stop(ctx context.Context) error {
	m.log.Info("AutoMigrate stop")
	return nil
}
