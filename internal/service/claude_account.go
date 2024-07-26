package service

import (
	v1 "PandoraPlusHelper/api/v1"
	"PandoraPlusHelper/internal/model"
	"PandoraPlusHelper/internal/repository"
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
	// token, err := s.claudeTokenRepository.GetToken(ctx, id)
	// if err != nil {
	// 	return v1.StatisticOpenaiAccountResponseData{}, err
	// }
	// accounts, err := s.claudeAccountService.SearchAccount(ctx, id)
	// if err != nil {
	// 	return v1.StatisticOpenaiAccountResponseData{}, err
	// }
	//
	// // 假设Account是你账户的结构体类型
	// var filteredAccounts []*model.ClaudeAccount
	// for _, account := range accounts {
	// 	if account.Status == 1 {
	// 		filteredAccounts = append(filteredAccounts, account)
	// 	}
	// }
	// // 重新赋值只包含状态为1的账户列表到accounts变量
	// accounts = filteredAccounts
	//
	// uniqueNames := make([]string, len(accounts))
	// infoList := make(map[string]util.ShareTokenInfo)
	// var models []string
	// series := make([]map[string]interface{}, 0)
	//
	// for i, account := range accounts {
	// 	uniqueNames[i] = account.Account
	// 	info, err := util.GetShareTokenInfo(account.ShareToken, token.AccessToken, s.logger)
	// 	if err != nil {
	// 		s.logger.Error("GetShareTokenInfo error", zap.Any("err", err))
	// 		continue
	// 	}
	// 	infoList[account.Account] = info
	// }
	//
	// for _, info := range infoList {
	// 	if info.Usage == nil {
	// 		s.logger.Error("获取分享用户信息失败, 请检查access_token是否有效")
	// 		continue
	// 	}
	// 	if _, ok := info.Usage["range"]; ok {
	// 		delete(info.Usage, "range")
	// 	}
	// 	for item := range info.Usage {
	// 		if !contains(models, item) {
	// 			models = append(models, item)
	// 		}
	// 	}
	// }
	//
	// for _, item := range models {
	// 	data := make([]interface{}, len(uniqueNames))
	// 	for i, uName := range uniqueNames {
	// 		if info, ok := infoList[uName]; ok {
	// 			if usage, ok := info.Usage[item]; ok {
	// 				data[i] = usage
	// 			} else {
	// 				data[i] = 0
	// 			}
	// 		} else {
	// 			data[i] = 0
	// 		}
	// 	}
	// 	series = append(series, map[string]interface{}{
	// 		"name": item,
	// 		"data": data,
	// 	})
	// }
	//
	// return v1.StatisticOpenaiAccountResponseData{
	// 	Categories: uniqueNames,
	// 	Series:     series,
	// }, nil
	return v1.StatisticOpenaiAccountResponseData{}, nil
}

// func contains(slice []string, item string) bool {
// 	for _, v := range slice {
// 		if v == item {
// 			return true
// 		}
// 	}
// 	return false
// }

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
