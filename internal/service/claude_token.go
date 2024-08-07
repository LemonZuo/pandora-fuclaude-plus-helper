package service

import (
	v1 "PandoraFuclaudePlusHelper/api/v1"
	"PandoraFuclaudePlusHelper/internal/model"
	"PandoraFuclaudePlusHelper/internal/repository"
	"context"
	"go.uber.org/zap"
	"time"
)

type ClaudeTokenService interface {
	RefreshToken(ctx context.Context, id int64) error
	GetToken(ctx context.Context, id int64) (*model.ClaudeToken, error)
	GetAllToken(ctx context.Context) ([]*model.ClaudeToken, error)
	Update(ctx context.Context, token *model.ClaudeToken) error
	Create(ctx context.Context, token *model.ClaudeToken) error
	SearchToken(ctx context.Context, keyword string) ([]*model.ClaudeToken, error)
	DeleteToken(ctx context.Context, id int64) error
	RefreshByToken(ctx context.Context, token *model.ClaudeToken) error
}

func NewClaudeTokenService(service *Service, claudeTokenRepository repository.ClaudeTokenRepository, claudeAccountRepository repository.ClaudeAccountRepository, coordinator *Coordinator) ClaudeTokenService {
	return &claudeTokenService{
		Service:                 service,
		claudeTokenRepository:   claudeTokenRepository,
		claudeAccountRepository: claudeAccountRepository,
		claudeAccountService:    coordinator.ClaudeAccountSvc,
	}
}

type claudeTokenService struct {
	*Service
	claudeTokenRepository   repository.ClaudeTokenRepository
	claudeAccountRepository repository.ClaudeAccountRepository
	claudeAccountService    ClaudeAccountService
}

func (s *claudeTokenService) RefreshToken(ctx context.Context, id int64) error {
	token, err := s.claudeTokenRepository.GetToken(ctx, id)
	if err != nil {
		s.logger.Error("GetToken error", zap.Any("err", err))
		return err
	}
	return s.RefreshByToken(ctx, token)
}

func (s *claudeTokenService) Update(ctx context.Context, token *model.ClaudeToken) error {
	err := s.claudeTokenRepository.Update(ctx, token)
	if err != nil {
		s.logger.Error("Update error", zap.Any("err", err))
		return err
	}
	return nil
}

func (s *claudeTokenService) Create(ctx context.Context, token *model.ClaudeToken) error {
	now := time.Now()
	token.CreateTime = now
	token.UpdateTime = now
	err := s.claudeTokenRepository.Create(ctx, token)
	if err != nil {
		s.logger.Error("Create error", zap.Any("err", err))
		return err
	}
	return nil
}

func (s *claudeTokenService) SearchToken(ctx context.Context, keyword string) ([]*model.ClaudeToken, error) {
	return s.claudeTokenRepository.SearchToken(ctx, keyword)
}

func (s *claudeTokenService) DeleteToken(ctx context.Context, id int64) error {
	accounts, err := s.claudeAccountService.SearchAccount(ctx, id)
	if err != nil {
		s.logger.Error("SearchAccount error", zap.Any("err", err))
		return err
	}
	if len(accounts) > 0 {
		return v1.ErrCannotDeleteToken
	}
	return s.claudeTokenRepository.DeleteToken(ctx, id)
}

func (s *claudeTokenService) GetToken(ctx context.Context, id int64) (*model.ClaudeToken, error) {
	token, err := s.claudeTokenRepository.GetToken(ctx, id)
	return token, err
}

func (s *claudeTokenService) GetAllToken(ctx context.Context) ([]*model.ClaudeToken, error) {
	return s.claudeTokenRepository.GetAllToken(ctx)
}

func (s *claudeTokenService) RefreshByToken(ctx context.Context, token *model.ClaudeToken) error {
	return nil
}
