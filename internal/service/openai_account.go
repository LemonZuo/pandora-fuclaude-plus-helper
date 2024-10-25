package service

import (
	v1 "PandoraFuclaudePlusHelper/api/v1"
	"PandoraFuclaudePlusHelper/internal/model"
	"PandoraFuclaudePlusHelper/internal/repository"
	"PandoraFuclaudePlusHelper/internal/util"
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type OpenaiAccountService interface {
	GetAccount(ctx context.Context, id int64) (*model.OpenaiAccount, error)
	Update(ctx context.Context, account *model.OpenaiAccount) error
	Create(ctx context.Context, account *model.OpenaiAccount) error
	SearchAccount(ctx context.Context, tokenId int64) ([]*model.OpenaiAccount, error)
	DeleteAccount(ctx context.Context, id int64) error
	StatisticAccount(ctx context.Context, id int64) (v1.StatisticOpenaiAccountResponseData, error)
	DisableAccount(ctx context.Context, id int64) error
	EnableAccount(ctx context.Context, id int64) error
}

func NewOpenaiAccountService(service *Service, openaiTokenRepository repository.OpenaiTokenRepository, openaiAccountRepository repository.OpenaiAccountRepository, coordinator *Coordinator) OpenaiAccountService {
	return &openaiAccountService{
		Service:                 service,
		openaiTokenRepository:   openaiTokenRepository,
		openaiAccountRepository: openaiAccountRepository,
		openaiAccountService:    coordinator.OpenaiAccountSvc,
	}
}

type openaiAccountService struct {
	*Service
	openaiTokenRepository   repository.OpenaiTokenRepository
	openaiAccountRepository repository.OpenaiAccountRepository
	openaiAccountService    OpenaiAccountService
}

func (s *openaiAccountService) Update(ctx context.Context, account *model.OpenaiAccount) error {
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
	token, err := s.openaiTokenRepository.GetToken(ctx, account.TokenID)
	if err != nil {
		return err
	}
	if token == nil {
		s.logger.Error("token not found")
		return fmt.Errorf("token not found")
	}

	now := time.Now()
	his.ExpirationTime = account.ExpirationTime
	his.Gpt35Limit = account.Gpt35Limit
	his.Gpt4Limit = account.Gpt4Limit
	his.Gpt4oLimit = account.Gpt4oLimit
	his.Gpt4oMiniLimit = account.Gpt4oMiniLimit
	his.O1Limit = account.O1Limit
	his.O1MiniLimit = account.O1MiniLimit
	his.ShowConversations = account.ShowConversations
	his.TemporaryChat = account.TemporaryChat
	// account.ShareToken = token.AccessToken
	account.ExpireAt = now.Add(time.Hour * 24 * 365)
	his.UpdateTime = now
	his.Status = account.Status

	// 生成共享token
	shareToken, shareTokenEncrypt, expireIn, err := util.GenShareToken(token.AccessToken,
		account.Account,
		0,
		account.Gpt35Limit,
		account.Gpt4Limit,
		account.Gpt4oLimit,
		account.Gpt4oMiniLimit,
		account.O1Limit,
		account.O1MiniLimit,
		account.ShowConversations == 1,
		false,
		false,
		account.TemporaryChat == 1,
		s.logger)
	if err != nil {
		s.logger.Error("GenerateShareToken error", zap.Any("err", err))
		return err
	}
	his.TokenID = account.TokenID
	his.ShareToken = shareToken
	his.ShareTokenEncrypt = shareTokenEncrypt
	his.ExpireAt = time.Unix(expireIn, 0)

	err = s.openaiAccountRepository.Update(ctx, his)
	if err != nil {
		s.logger.Error("Create error", zap.Any("err", err))
		return err
	}
	return nil
}

func (s *openaiAccountService) Create(ctx context.Context, account *model.OpenaiAccount) error {
	// 查询token是否存在
	token, err := s.openaiTokenRepository.GetToken(ctx, account.TokenID)
	if err != nil {
		return err
	}
	if token == nil {
		s.logger.Error("token not found")
		return fmt.Errorf("token not found")
	}

	now := time.Now()
	account.ShareToken = token.AccessToken
	account.ExpireAt = now.Add(time.Hour * 24 * 365)
	// 生成共享token
	shareToken, shareTokenEncrypt, expireIn, err := util.GenShareToken(token.AccessToken,
		account.Account,
		0,
		account.Gpt35Limit,
		account.Gpt4Limit,
		account.Gpt4oLimit,
		account.Gpt4oMiniLimit,
		account.O1Limit,
		account.O1MiniLimit,
		account.ShowConversations == 1,
		false,
		false,
		account.TemporaryChat == 1,
		s.logger)
	if err != nil {
		s.logger.Error("GenerateShareToken error", zap.Any("err", err))
		return err
	}
	account.ShareToken = shareToken
	account.ShareTokenEncrypt = shareTokenEncrypt
	account.ExpireAt = time.Unix(expireIn, 0)
	account.CreateTime = now
	account.UpdateTime = now
	err = s.openaiAccountRepository.Create(ctx, account)

	if err != nil {
		s.logger.Error("Create error", zap.Any("err", err))
		return err
	}
	return nil
}

func (s *openaiAccountService) SearchAccount(ctx context.Context, tokenId int64) ([]*model.OpenaiAccount, error) {
	return s.openaiAccountRepository.SearchAccount(ctx, tokenId)
}

func (s *openaiAccountService) DeleteAccount(ctx context.Context, id int64) error {
	account, err := s.GetAccount(ctx, id)
	if err != nil {
		s.logger.Error("GetAccount error", zap.Any("err", err))
		return err
	}
	if account == nil {
		s.logger.Error("account not found")
		return fmt.Errorf("account not found")
	}

	// 查询token是否存在
	token, err := s.openaiTokenRepository.GetToken(ctx, account.TokenID)
	if err != nil {
		return err
	}
	if token == nil {
		s.logger.Error("token not found")
		return fmt.Errorf("token not found")
	}

	shareToken, _, expireIn, err := util.GenShareToken(token.AccessToken,
		account.Account,
		-1,
		0,
		0,
		0,
		0,
		0,
		0,
		false,
		false,
		false,
		account.TemporaryChat == 1,
		s.logger)
	s.logger.Info("DeleteAccount", zap.Any("shareToken", shareToken), zap.Any("expireIn", expireIn))
	if err != nil {
		return err
	}
	return s.openaiAccountRepository.DeleteAccount(ctx, id)
}

func (s *openaiAccountService) GetAccount(ctx context.Context, id int64) (*model.OpenaiAccount, error) {
	return s.openaiAccountRepository.GetAccount(ctx, id)
}

func (s *openaiAccountService) StatisticAccount(ctx context.Context, id int64) (v1.StatisticOpenaiAccountResponseData, error) {
	token, err := s.openaiTokenRepository.GetToken(ctx, id)
	if err != nil {
		return v1.StatisticOpenaiAccountResponseData{}, err
	}
	accounts, err := s.openaiAccountService.SearchAccount(ctx, id)
	if err != nil {
		return v1.StatisticOpenaiAccountResponseData{}, err
	}

	// 假设Account是你账户的结构体类型
	var filteredAccounts []*model.OpenaiAccount
	for _, account := range accounts {
		if account.Status == 1 {
			filteredAccounts = append(filteredAccounts, account)
		}
	}
	// 重新赋值只包含状态为1的账户列表到accounts变量
	accounts = filteredAccounts

	uniqueNames := make([]string, len(accounts))
	infoList := make(map[string]util.ShareTokenInfo)
	var models []string
	series := make([]map[string]interface{}, 0)

	for i, account := range accounts {
		uniqueNames[i] = account.Account
		info, err := util.GetShareTokenInfo(account.ShareToken, token.AccessToken, s.logger)
		if err != nil {
			s.logger.Error("GetShareTokenInfo error", zap.Any("err", err))
			continue
		}
		infoList[account.Account] = info
	}

	for _, info := range infoList {
		if info.Usage == nil {
			s.logger.Error("获取分享用户信息失败, 请检查access_token是否有效")
			continue
		}
		if _, ok := info.Usage["range"]; ok {
			delete(info.Usage, "range")
		}
		for item := range info.Usage {
			if !contains(models, item) {
				models = append(models, item)
			}
		}
	}

	for _, item := range models {
		data := make([]interface{}, len(uniqueNames))
		for i, uName := range uniqueNames {
			if info, ok := infoList[uName]; ok {
				if usage, ok := info.Usage[item]; ok {
					data[i] = usage
				} else {
					data[i] = 0
				}
			} else {
				data[i] = 0
			}
		}
		series = append(series, map[string]interface{}{
			"name": item,
			"data": data,
		})
	}

	return v1.StatisticOpenaiAccountResponseData{
		Categories: uniqueNames,
		Series:     series,
	}, nil
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (s *openaiAccountService) DisableAccount(ctx context.Context, id int64) error {
	account, err := s.GetAccount(ctx, id)
	if err != nil {
		s.logger.Error("GetAccount error", zap.Any("err", err))
		return err
	}
	if account == nil {
		s.logger.Error("account not found")
		return fmt.Errorf("account not found")
	}

	// 查询token是否存在
	token, err := s.openaiTokenRepository.GetToken(ctx, account.TokenID)
	if err != nil {
		return err
	}
	if token == nil {
		s.logger.Error("token not found")
		return fmt.Errorf("token not found")
	}

	_, _, _, err = util.GenShareToken(token.AccessToken,
		account.Account,
		-1,
		account.Gpt35Limit,
		account.Gpt4Limit,
		account.Gpt4oLimit,
		account.Gpt4oMiniLimit,
		account.O1Limit,
		account.O1MiniLimit,
		account.ShowConversations == 1,
		false,
		false,
		account.TemporaryChat == 1,
		s.logger)
	if err != nil {
		return err
	}
	now := time.Now()

	account.Status = 0
	account.ExpirationTime = now
	account.ExpireAt = now
	account.UpdateTime = time.Now()
	err = s.openaiAccountRepository.Update(ctx, account)
	if err != nil {
		s.logger.Error("Update error", zap.Any("err", err))
		return err
	}
	return nil
}

func (s *openaiAccountService) EnableAccount(ctx context.Context, id int64) error {
	account, err := s.GetAccount(ctx, id)
	if err != nil {
		s.logger.Error("GetAccount error", zap.Any("err", err))
		return err
	}
	if account == nil {
		s.logger.Error("account not found")
		return fmt.Errorf("account not found")
	}

	// 查询token是否存在
	token, err := s.openaiTokenRepository.GetToken(ctx, account.TokenID)
	if err != nil {
		return err
	}
	if token == nil {
		s.logger.Error("token not found")
		return fmt.Errorf("token not found")
	}

	// 生成共享token
	shareToken, shareTokenEncrypt, expireIn, err := util.GenShareToken(token.AccessToken,
		account.Account,
		0,
		account.Gpt35Limit,
		account.Gpt4Limit,
		account.Gpt4oLimit,
		account.Gpt4oMiniLimit,
		account.O1Limit,
		account.O1MiniLimit,
		account.ShowConversations == 1,
		false,
		false,
		account.TemporaryChat == 1,
		s.logger)
	if err != nil {
		s.logger.Error("GenerateShareToken error", zap.Any("err", err))
		return err
	}
	now := time.Now()
	account.ShareToken = shareToken
	account.ShareTokenEncrypt = shareTokenEncrypt
	account.ExpireAt = time.Unix(expireIn, 0)
	account.Status = 1
	// 有效期一个月
	account.ExpirationTime = now.Add(time.Hour * 24 * 30)
	account.UpdateTime = now
	err = s.openaiAccountRepository.Update(ctx, account)
	if err != nil {
		s.logger.Error("Update error", zap.Any("err", err))
		return err
	}
	return nil
}
