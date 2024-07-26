package service

import (
	v1 "PandoraPlusHelper/api/v1"
	"PandoraPlusHelper/internal/model"
	"PandoraPlusHelper/internal/repository"
	"PandoraPlusHelper/internal/util"
	"context"
	"errors"
	"go.uber.org/zap"
	"time"
)

type OpenaiTokenService interface {
	RefreshToken(ctx context.Context, id int64) error
	GetToken(ctx context.Context, id int64) (*model.OpenaiToken, error)
	GetAllToken(ctx context.Context) ([]*model.OpenaiToken, error)
	Update(ctx context.Context, token *model.OpenaiToken) error
	Create(ctx context.Context, token *model.OpenaiToken) error
	SearchToken(ctx context.Context, keyword string) ([]*model.OpenaiToken, error)
	DeleteToken(ctx context.Context, id int64) error
	RefreshByToken(ctx context.Context, token *model.OpenaiToken) error
}

func NewOpenaiTokenService(service *Service, openaiTokenRepository repository.OpenaiTokenRepository, openaiAccountRepository repository.OpenaiAccountRepository, coordinator *Coordinator) OpenaiTokenService {
	return &openaiTokenService{
		Service:                 service,
		openaiTokenRepository:   openaiTokenRepository,
		openaiAccountRepository: openaiAccountRepository,
		openaiAccountService:    coordinator.OpenaiAccountSvc,
	}
}

type openaiTokenService struct {
	*Service
	openaiTokenRepository   repository.OpenaiTokenRepository
	openaiAccountRepository repository.OpenaiAccountRepository
	openaiAccountService    OpenaiAccountService
}

func (s *openaiTokenService) RefreshToken(ctx context.Context, id int64) error {
	token, err := s.openaiTokenRepository.GetToken(ctx, id)
	if err != nil {
		s.logger.Error("GetToken error", zap.Any("err", err))
		return err
	}
	return s.RefreshByToken(ctx, token)
}

func (s *openaiTokenService) Update(ctx context.Context, token *model.OpenaiToken) error {
	return s.RefreshByToken(ctx, token)
}

func (s *openaiTokenService) Create(ctx context.Context, token *model.OpenaiToken) error {
	now := time.Now()
	// 默认的类型处理
	token.AccessToken = token.RefreshToken
	token.PlusSubscription = 0
	token.ExpireAt = now.Add(time.Hour * time.Duration(24*365))
	// 使用RefreshToken生成AccessToken
	accessToken, expiresIn, err := util.GenAccessToken(token.RefreshToken, s.logger)
	if err != nil {
		return err
	}
	// 判断订阅状态
	plusSubscription := util.CheckSubscriptionStatus(accessToken, s.logger)
	token.AccessToken = accessToken
	token.PlusSubscription = plusSubscription
	token.ExpireAt = now.Add(time.Second * time.Duration(expiresIn))
	token.CreateTime = now
	token.UpdateTime = now

	err = s.openaiTokenRepository.Create(ctx, token)
	if err != nil {
		s.logger.Error("Create error", zap.Any("err", err))
		return err
	}
	return nil
}

func (s *openaiTokenService) SearchToken(ctx context.Context, keyword string) ([]*model.OpenaiToken, error) {
	return s.openaiTokenRepository.SearchToken(ctx, keyword)
}

func (s *openaiTokenService) DeleteToken(ctx context.Context, id int64) error {
	accounts, err := s.openaiAccountService.SearchAccount(ctx, id)
	if err != nil {
		s.logger.Error("SearchAccount error", zap.Any("err", err))
		return err
	}
	if len(accounts) > 0 {
		return v1.ErrCannotDeleteToken
	}
	return s.openaiTokenRepository.DeleteToken(ctx, id)
}

func (s *openaiTokenService) GetToken(ctx context.Context, id int64) (*model.OpenaiToken, error) {
	token, err := s.openaiTokenRepository.GetToken(ctx, id)
	return token, err
}

func (s *openaiTokenService) GetAllToken(ctx context.Context) ([]*model.OpenaiToken, error) {
	return s.openaiTokenRepository.GetAllToken(ctx)
}

func (s *openaiTokenService) RefreshByToken(ctx context.Context, token *model.OpenaiToken) error {
	//  token是否存在
	if token.ID == 0 {
		return errors.New("token not found")
	}

	if len(token.RefreshToken) == 0 {
		return errors.New("refresh token is empty")
	}

	his, err := s.openaiTokenRepository.GetToken(ctx, token.ID)
	if err != nil {
		return err
	}
	if his == nil {
		return errors.New("token not exist")
	}

	now := time.Now()
	his.TokenName = token.TokenName
	his.AccessToken = token.RefreshToken
	his.RefreshToken = token.RefreshToken
	his.ExpireAt = now.Add(time.Hour * time.Duration(24*365))
	// 使用RefreshToken生成AccessToken
	accessToken, expiresIn, err := util.GenAccessToken(token.RefreshToken, s.logger)
	if err != nil {
		s.logger.Error("GetAccessTokenByRefreshToken error", zap.Any("err", err))
		return err
	}
	// 判断订阅状态
	plusSubscription := util.CheckSubscriptionStatus(accessToken, s.logger)
	his.PlusSubscription = plusSubscription
	his.AccessToken = accessToken
	his.ExpireAt = now.Add(time.Second * time.Duration(expiresIn))

	his.UpdateTime = now

	err = s.openaiTokenRepository.Update(ctx, his)
	if err != nil {
		return err
	}
	// 刷新此Token的所有AccountToken
	accounts, err := s.openaiAccountService.SearchAccount(ctx, his.ID)
	if err != nil {
		return err
	}
	for _, account := range accounts {
		if account.Status == 0 {
			s.logger.Info("OpenaiAccount is disabled", zap.Any("account", account))
			continue
		}
		now := time.Now()
		// 默认设置为AccessToken
		account.ShareToken = his.AccessToken
		account.ExpireAt = now.Add(time.Hour * time.Duration(24*365))
		shareToken, expireIn, err := util.GenShareToken(his.AccessToken,
			account.Account,
			0,
			account.Gpt35Limit,
			account.Gpt4Limit,
			account.ShowConversations == 1,
			false,
			false,
			account.TemporaryChat == 1,
			s.logger)
		if err != nil {
			s.logger.Error("GenerateShareToken error", zap.Any("err", err))
			continue
		}
		account.ShareToken = shareToken
		account.ExpireAt = time.Unix(expireIn, 0)
		account.UpdateTime = now
		err = s.openaiAccountRepository.Update(ctx, account)
		if err != nil {
			s.logger.Error("Update error", zap.Any("err", err))
		}
	}
	return nil
}
