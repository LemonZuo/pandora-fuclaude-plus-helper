package service

import (
	v1 "PandoraFuclaudePlusHelper/api/v1"
	"PandoraFuclaudePlusHelper/internal/model"
	"PandoraFuclaudePlusHelper/internal/repository"
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type ClaudeAccountService interface {
	GetAccount(ctx context.Context, id int64) (*model.ClaudeAccount, error)
	Update(ctx context.Context, account *model.ClaudeAccount) error
	Create(ctx context.Context, account *model.ClaudeAccount) error
	SearchAccount(ctx context.Context, tokenId int64) ([]*model.ClaudeAccount, error)
	DeleteAccount(ctx context.Context, id int64) error
	StatisticAccount(ctx context.Context, id int64) (v1.StatisticOpenaiAccountResponseData, error)
	DisableAccount(ctx context.Context, id int64) error
	EnableAccount(ctx context.Context, id int64) error
}

func NewClaudeAccountService(service *Service, claudeTokenRepository repository.ClaudeTokenRepository, claudeAccountRepository repository.ClaudeAccountRepository, coordinator *Coordinator) ClaudeAccountService {
	return &claudeAccountService{
		Service:                 service,
		claudeTokenRepository:   claudeTokenRepository,
		claudeAccountRepository: claudeAccountRepository,
		claudeAccountService:    coordinator.ClaudeAccountSvc,
	}
}

type claudeAccountService struct {
	*Service
	claudeTokenRepository   repository.ClaudeTokenRepository
	claudeAccountRepository repository.ClaudeAccountRepository
	claudeAccountService    ClaudeAccountService
}

func (s *claudeAccountService) Update(ctx context.Context, account *model.ClaudeAccount) error {
	his, err := s.GetAccount(ctx, account.ID)
	if err != nil {
		s.logger.Error("GetAccount error", zap.Any("err", err))
		return err
	}
	if his == nil {
		s.logger.Error("account not found")
		return fmt.Errorf("account not found")
	}
	// 查询token是否存在
	token, err := s.claudeTokenRepository.GetToken(ctx, account.TokenID)
	if err != nil {
		return err
	}
	if token == nil {
		s.logger.Error("token not found")
		return fmt.Errorf("token not found")
	}

	now := time.Now()
	his.UserId = account.UserId
	his.TokenID = account.TokenID
	his.Account = account.Account
	his.Status = account.Status
	his.UpdateTime = now

	err = s.claudeAccountRepository.Update(ctx, his)
	if err != nil {
		s.logger.Error("Create error", zap.Any("err", err))
		return err
	}
	return nil
}

func (s *claudeAccountService) Create(ctx context.Context, account *model.ClaudeAccount) error {
	// 查询token是否存在
	token, err := s.claudeTokenRepository.GetToken(ctx, account.TokenID)
	if err != nil {
		return err
	}
	if token == nil {
		s.logger.Error("token not found")
		return fmt.Errorf("token not found")
	}

	now := time.Now()

	account.CreateTime = now
	account.UpdateTime = now
	err = s.claudeAccountRepository.Create(ctx, account)
	if err != nil {
		s.logger.Error("Create error", zap.Any("err", err))
		return err
	}
	return nil
}

func (s *claudeAccountService) SearchAccount(ctx context.Context, tokenId int64) ([]*model.ClaudeAccount, error) {
	return s.claudeAccountRepository.SearchAccount(ctx, tokenId)
}

func (s *claudeAccountService) DeleteAccount(ctx context.Context, id int64) error {
	return s.claudeAccountRepository.DeleteAccount(ctx, id)
}

func (s *claudeAccountService) GetAccount(ctx context.Context, id int64) (*model.ClaudeAccount, error) {
	return s.claudeAccountRepository.GetAccount(ctx, id)
}

func (s *claudeAccountService) StatisticAccount(ctx context.Context, id int64) (v1.StatisticOpenaiAccountResponseData, error) {
	return v1.StatisticOpenaiAccountResponseData{}, nil
}

func (s *claudeAccountService) DisableAccount(ctx context.Context, id int64) error {
	account, err := s.GetAccount(ctx, id)
	if err != nil {
		s.logger.Error("GetAccount error", zap.Any("err", err))
		return err
	}
	if account == nil {
		s.logger.Error("account not found")
		return fmt.Errorf("account not found")
	}

	now := time.Now()
	account.Status = 0
	account.UpdateTime = now
	err = s.claudeAccountRepository.Update(ctx, account)
	if err != nil {
		s.logger.Error("Update error", zap.Any("err", err))
		return err
	}
	return nil
}

func (s *claudeAccountService) EnableAccount(ctx context.Context, id int64) error {
	account, err := s.GetAccount(ctx, id)
	if err != nil {
		s.logger.Error("GetAccount error", zap.Any("err", err))
		return err
	}
	if account == nil {
		s.logger.Error("account not found")
		return fmt.Errorf("account not found")
	}
	now := time.Now()
	account.Status = 1
	// 有效期一个月
	account.UpdateTime = now
	err = s.claudeAccountRepository.Update(ctx, account)
	if err != nil {
		s.logger.Error("Update error", zap.Any("err", err))
		return err
	}
	return nil
}
