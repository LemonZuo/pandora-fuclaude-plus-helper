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
	//  token是否存在
	// if token.ID == 0 {
	// 	return errors.New("token not found")
	// }
	//
	// if len(token.RefreshToken) == 0 {
	// 	return errors.New("refresh token is empty")
	// }
	//
	// his, err := s.claudeTokenRepository.GetToken(ctx, token.ID)
	// if err != nil {
	// 	return err
	// }
	// if his == nil {
	// 	return errors.New("token not exist")
	// }
	//
	// now := time.Now()
	// his.TokenName = token.TokenName
	// his.AccessToken = token.RefreshToken
	// his.RefreshToken = token.RefreshToken
	// his.ExpireAt = now.Add(time.Hour * time.Duration(24*365))
	// // 使用RefreshToken生成AccessToken
	// accessToken, expiresIn, err := util.GenAccessToken(token.RefreshToken, s.logger)
	// if err != nil {
	// 	s.logger.Error("GetAccessTokenByRefreshToken error", zap.Any("err", err))
	// 	return err
	// }
	// // 判断订阅状态
	// plusSubscription := util.CheckSubscriptionStatus(accessToken, s.logger)
	// his.PlusSubscription = plusSubscription
	// his.AccessToken = accessToken
	// his.ExpireAt = now.Add(time.Second * time.Duration(expiresIn))
	//
	// his.UpdateTime = now
	//
	// err = s.claudeTokenRepository.Update(ctx, his)
	// if err != nil {
	// 	return err
	// }
	// // 刷新此Token的所有AccountToken
	// accounts, err := s.claudeAccountService.SearchAccount(ctx, his.ID)
	// if err != nil {
	// 	return err
	// }
	// for _, account := range accounts {
	// 	if account.Status == 0 {
	// 		s.logger.Info("ClaudeAccount is disabled", zap.Any("account", account))
	// 		continue
	// 	}
	// 	now := time.Now()
	// 	// 默认设置为AccessToken
	// 	account.ShareToken = his.AccessToken
	// 	account.ExpireAt = now.Add(time.Hour * time.Duration(24*365))
	// 	shareToken, expireIn, err := util.GenShareToken(his.AccessToken,
	// 		account.Account,
	// 		0,
	// 		account.Gpt35Limit,
	// 		account.Gpt4Limit,
	// 		account.ShowConversations == 1,
	// 		false,
	// 		false,
	// 		account.TemporaryChat == 1,
	// 		s.logger)
	// 	if err != nil {
	// 		s.logger.Error("GenerateShareToken error", zap.Any("err", err))
	// 		continue
	// 	}
	// 	account.ShareToken = shareToken
	// 	account.ExpireAt = time.Unix(expireIn, 0)
	// 	account.UpdateTime = now
	// 	err = s.claudeAccountRepository.Update(ctx, account)
	// 	if err != nil {
	// 		s.logger.Error("Update error", zap.Any("err", err))
	// 	}
	// }
	return nil
}
