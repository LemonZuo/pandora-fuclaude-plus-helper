package service

import (
	"PandoraPlusHelper/internal/model"
	"PandoraPlusHelper/internal/repository"
	"context"
	"errors"
	"go.uber.org/zap"
	"time"
)

type UserService interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, id int64) (*model.User, error)
	GetAllUser(ctx context.Context) ([]*model.User, error)
	SearchUser(ctx context.Context, keyword string) ([]*model.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

func NewUserService(service *Service, userRepository repository.UserRepository,
	openaiTokenRepository repository.OpenaiTokenRepository, openaiAccountRepository repository.OpenaiAccountRepository,
	claudeTokenRepository repository.ClaudeTokenRepository, claudeAccountRepository repository.ClaudeAccountRepository,
	coordinator *Coordinator) UserService {
	return &userService{
		Service:                 service,
		userRepository:          userRepository,
		openaiTokenRepository:   openaiTokenRepository,
		openaiAccountRepository: openaiAccountRepository,
		claudeTokenRepository:   claudeTokenRepository,
		claudeAccountRepository: claudeAccountRepository,
		openaiAccountService:    coordinator.OpenaiAccountSvc,
		claudeAccountService:    coordinator.ClaudeAccountSvc,
	}
}

type userService struct {
	*Service
	userRepository          repository.UserRepository
	openaiTokenRepository   repository.OpenaiTokenRepository
	openaiAccountRepository repository.OpenaiAccountRepository
	claudeTokenRepository   repository.ClaudeTokenRepository
	claudeAccountRepository repository.ClaudeAccountRepository
	openaiAccountService    OpenaiAccountService
	claudeAccountService    ClaudeAccountService
}

func (s *userService) Create(ctx context.Context, user *model.User) error {
	now := time.Now()
	// 默认的类型处理
	if user.ExpirationTime.IsZero() {
		user.ExpirationTime = now.Add(time.Hour * time.Duration(24*365))
	}
	user.CreateTime = now
	user.UpdateTime = now

	err := s.userRepository.Create(ctx, user)
	if err != nil {
		s.logger.Error("Create error", zap.Any("err", err))
		return err
	}
	// 未启用账户，新增完毕直接返回
	if user.Enable != 1 {
		return nil
	}
	// 处理openai
	if user.Openai == 1 && user.OpenaiToken > 0 {
		token, err := s.openaiTokenRepository.GetToken(ctx, user.OpenaiToken)
		if err != nil {
			s.logger.Error("GetToken error", zap.Any("err", err))
			return err
		}
		if token == nil {
			s.logger.Error("token not found")
			return errors.New("token not found")
		}
		// 组装一个Account
		account := &model.OpenaiAccount{
			UserId:            user.ID,
			Account:           user.UniqueName,
			ExpirationTime:    user.ExpirationTime,
			Status:            1,
			Gpt35Limit:        -1,
			Gpt4Limit:         -1,
			ShowConversations: 0,
			TemporaryChat:     0,
			TokenID:           token.ID,
		}
		err = s.openaiAccountService.Create(ctx, account)
		if err != nil {
			s.logger.Error("Create error", zap.Any("err", err))
			return err
		}
	}

	if user.Claude == 1 && user.ClaudeToken > 0 {
		token, err := s.claudeTokenRepository.GetToken(ctx, user.ClaudeToken)
		if err != nil {
			s.logger.Error("GetToken error", zap.Any("err", err))
			return err
		}
		if token == nil {
			s.logger.Error("token not found")
			return errors.New("token not found")
		}
		// 组装一个Account
		account := &model.ClaudeAccount{
			UserId:  user.ID,
			Account: user.UniqueName,
			Status:  1,
			TokenID: token.ID,
		}
		err = s.claudeAccountService.Create(ctx, account)
		if err != nil {
			s.logger.Error("Create error", zap.Any("err", err))
			return err
		}
	}
	return nil
}

func (s *userService) Update(ctx context.Context, user *model.User) error {
	his, err := s.userRepository.GetUser(ctx, user.ID)
	if err != nil {
		s.logger.Error("Update error", zap.Any("err", err))
		return err
	}

	account, err := s.openaiAccountRepository.GetAccountByUserId(ctx, user.ID)
	if err != nil {
		s.logger.Error("GetAccountByUserId error", zap.Any("err", err))
	}

	claudeAccount, err := s.claudeAccountRepository.GetAccountByUserId(ctx, user.ID)
	if err != nil {
		s.logger.Error("GetAccountByUserId error", zap.Any("err", err))
	}

	// 处理 Enable 状态
	if user.Enable != 1 {
		// 禁用用户
		user.Openai = 0
		user.Claude = 0
		if account != nil {
			err := s.openaiAccountService.DisableAccount(ctx, account.ID)
			if err != nil {
				s.logger.Error("DisableAccount error", zap.Any("err", err))
				return err
			}
		}
		if claudeAccount != nil {
			err := s.claudeAccountService.DisableAccount(ctx, claudeAccount.ID)
			if err != nil {
				s.logger.Error("DisableAccount error", zap.Any("err", err))
				return err
			}
		}
	} else {
		// 处理 OpenAI
		if user.Openai == 1 && user.OpenaiToken > 0 {
			token, err := s.openaiTokenRepository.GetToken(ctx, user.OpenaiToken)
			if err != nil {
				s.logger.Error("GetToken error", zap.Any("err", err))
				return err
			}
			if token == nil {
				s.logger.Error("token not found")
				return errors.New("token not found")
			}

			if account == nil {
				account = &model.OpenaiAccount{
					UserId:            user.ID,
					Account:           user.UniqueName,
					ExpirationTime:    user.ExpirationTime,
					Status:            1,
					Gpt35Limit:        -1,
					Gpt4Limit:         -1,
					ShowConversations: 0,
					TemporaryChat:     0,
					TokenID:           user.OpenaiToken,
				}
				err = s.openaiAccountService.Create(ctx, account)
				if err != nil {
					s.logger.Error("Create error", zap.Any("err", err))
					return err
				}
			} else if his.OpenaiToken != user.OpenaiToken {
				account.TokenID = user.OpenaiToken
				account.ExpirationTime = user.ExpirationTime
				account.Status = 1
				err := s.openaiAccountService.Update(ctx, account)
				if err != nil {
					s.logger.Error("Update error", zap.Any("err", err))
					return err
				}
			} else {
				err := s.openaiAccountService.EnableAccount(ctx, account.ID)
				if err != nil {
					s.logger.Error("EnableAccount error", zap.Any("err", err))
					return err
				}
			}
		} else if account != nil {
			err := s.openaiAccountService.DisableAccount(ctx, account.ID)
			if err != nil {
				s.logger.Error("DisableAccount error", zap.Any("err", err))
				return err
			}
		}

		// 处理 Claude
		if user.Claude == 1 && user.ClaudeToken > 0 {
			token, err := s.claudeTokenRepository.GetToken(ctx, user.ClaudeToken)
			if err != nil {
				s.logger.Error("GetToken error", zap.Any("err", err))
				return err
			}
			if token == nil {
				s.logger.Error("token not found")
				return errors.New("token not found")
			}

			if claudeAccount == nil {
				claudeAccount = &model.ClaudeAccount{
					UserId:  user.ID,
					Account: user.UniqueName,
					Status:  1,
					TokenID: user.ClaudeToken,
				}
				err = s.claudeAccountService.Create(ctx, claudeAccount)
				if err != nil {
					s.logger.Error("Create error", zap.Any("err", err))
					return err
				}
			} else if his.ClaudeToken != user.ClaudeToken {
				claudeAccount.TokenID = user.ClaudeToken
				claudeAccount.Status = 1
				err := s.claudeAccountService.Update(ctx, claudeAccount)
				if err != nil {
					s.logger.Error("Update error", zap.Any("err", err))
					return err
				}
			} else {
				err := s.claudeAccountService.EnableAccount(ctx, claudeAccount.ID)
				if err != nil {
					s.logger.Error("EnableAccount error", zap.Any("err", err))
					return err
				}
			}
		} else if claudeAccount != nil {
			err := s.claudeAccountService.DisableAccount(ctx, claudeAccount.ID)
			if err != nil {
				s.logger.Error("DisableAccount error", zap.Any("err", err))
				return err
			}
		}
	}

	// 更新属性
	his.UniqueName = user.UniqueName
	his.Password = user.Password
	his.Enable = user.Enable
	his.Openai = user.Openai
	his.OpenaiToken = user.OpenaiToken
	his.Claude = user.Claude
	his.ClaudeToken = user.ClaudeToken
	if !user.ExpirationTime.IsZero() {
		his.ExpirationTime = user.ExpirationTime
	}
	his.UpdateTime = time.Now()
	err = s.userRepository.Update(ctx, his)
	if err != nil {
		s.logger.Error("Update error", zap.Any("err", err))
		return err
	}

	return nil
}

func (s *userService) SearchUser(ctx context.Context, keyword string) ([]*model.User, error) {
	return s.userRepository.SearchUser(ctx, keyword)
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	account, err := s.openaiAccountRepository.GetAccountByUserId(ctx, id)
	if err != nil {
		s.logger.Error("GetAccountByUserId error", zap.Any("err", err))
		return err
	}
	// 1.删除 openai account
	if account != nil {
		err = s.openaiAccountService.DeleteAccount(ctx, account.ID)
		if err != nil {
			s.logger.Error("DeleteAccount error", zap.Any("err", err))
			return err
		}
	}
	// 2.删除 claude account
	claudeAccount, err := s.claudeAccountRepository.GetAccountByUserId(ctx, id)
	if err != nil {
		s.logger.Error("GetAccountByUserId error", zap.Any("err", err))
		return err
	}
	if claudeAccount != nil {
		err = s.claudeAccountService.DeleteAccount(ctx, claudeAccount.ID)
		if err != nil {
			s.logger.Error("DeleteAccount error", zap.Any("err", err))
			return err
		}
	}

	// 3.删除user
	err = s.userRepository.DeleteUser(ctx, id)
	if err != nil {
		s.logger.Error("DeleteUser error", zap.Any("err", err))
		return err
	}
	return nil

}

func (s *userService) GetUser(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.userRepository.GetUser(ctx, id)
	return user, err
}

func (s *userService) GetAllUser(ctx context.Context) ([]*model.User, error) {
	return s.userRepository.GetAllUser(ctx)
}
