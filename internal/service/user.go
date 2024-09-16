package service

import (
	"PandoraFuclaudePlusHelper/internal/model"
	"PandoraFuclaudePlusHelper/internal/repository"
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
	// 获取当前用户信息
	his, err := s.userRepository.GetUser(ctx, user.ID)
	if err != nil {
		s.logger.Error("Failed to get user", zap.Any("err", err))
		return err
	}

	// 获取 OpenAI 账号信息
	account, err := s.openaiAccountRepository.GetAccountByUserId(ctx, user.ID)
	if err != nil {
		s.logger.Error("Failed to get OpenAI account", zap.Any("err", err))
	}

	// 获取 Claude 账号信息
	claudeAccount, err := s.claudeAccountRepository.GetAccountByUserId(ctx, user.ID)
	if err != nil {
		s.logger.Error("Failed to get Claude account", zap.Any("err", err))
	}

	// 处理用户启用状态
	if user.Enable != 1 {
		// 禁用用户的 OpenAI 和 Claude 服务
		user.Openai = 0
		user.Claude = 0

		// 禁用 OpenAI 账号
		if account != nil {
			if err := s.openaiAccountService.DisableAccount(ctx, account.ID); err != nil {
				s.logger.Error("Failed to disable OpenAI account", zap.Any("err", err))
				return err
			}
		}

		// 禁用 Claude 账号
		if claudeAccount != nil {
			if err := s.claudeAccountService.DisableAccount(ctx, claudeAccount.ID); err != nil {
				s.logger.Error("Failed to disable Claude account", zap.Any("err", err))
				return err
			}
		}
	} else {
		// 处理 OpenAI 服务
		if user.Openai == 1 && user.OpenaiToken > 0 {
			// 获取 OpenAI Token
			token, err := s.openaiTokenRepository.GetToken(ctx, user.OpenaiToken)
			if err != nil || token == nil {
				s.logger.Error("OpenAI token not found", zap.Any("err", err))
				return errors.New("OpenAI token not found")
			}

			if account == nil {
				// 创建新的 OpenAI 账号
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
				if err := s.openaiAccountService.Create(ctx, account); err != nil {
					s.logger.Error("Failed to create OpenAI account", zap.Any("err", err))
					return err
				}
			} else {
				// 更新或启用现有的 OpenAI 账号
				if his.OpenaiToken != user.OpenaiToken {
					account.TokenID = user.OpenaiToken
					account.ExpirationTime = user.ExpirationTime
					account.Status = 1
					if err := s.openaiAccountService.Update(ctx, account); err != nil {
						s.logger.Error("Failed to update OpenAI account", zap.Any("err", err))
						return err
					}
				} else if err := s.openaiAccountService.EnableAccount(ctx, account.ID); err != nil {
					s.logger.Error("Failed to enable OpenAI account", zap.Any("err", err))
					return err
				}
			}
		} else if account != nil {
			// 禁用不需要的 OpenAI 账号
			if err := s.openaiAccountService.DisableAccount(ctx, account.ID); err != nil {
				s.logger.Error("Failed to disable OpenAI account", zap.Any("err", err))
				return err
			}
		}

		// 处理 Claude 服务
		if user.Claude == 1 && user.ClaudeToken > 0 {
			// 获取 Claude Token
			token, err := s.claudeTokenRepository.GetToken(ctx, user.ClaudeToken)
			if err != nil || token == nil {
				s.logger.Error("Claude token not found", zap.Any("err", err))
				return errors.New("Claude token not found")
			}

			if claudeAccount == nil {
				// 创建新的 Claude 账号
				claudeAccount = &model.ClaudeAccount{
					UserId:  user.ID,
					Account: user.UniqueName,
					Status:  1,
					TokenID: user.ClaudeToken,
				}
				if err := s.claudeAccountService.Create(ctx, claudeAccount); err != nil {
					s.logger.Error("Failed to create Claude account", zap.Any("err", err))
					return err
				}
			} else {
				// 更新或启用现有的 Claude 账号
				if his.ClaudeToken != user.ClaudeToken {
					claudeAccount.TokenID = user.ClaudeToken
					claudeAccount.Status = 1
					if err := s.claudeAccountService.Update(ctx, claudeAccount); err != nil {
						s.logger.Error("Failed to update Claude account", zap.Any("err", err))
						return err
					}
				} else if err := s.claudeAccountService.EnableAccount(ctx, claudeAccount.ID); err != nil {
					s.logger.Error("Failed to enable Claude account", zap.Any("err", err))
					return err
				}
			}
		} else if claudeAccount != nil {
			// 禁用不需要的 Claude 账号
			if err := s.claudeAccountService.DisableAccount(ctx, claudeAccount.ID); err != nil {
				s.logger.Error("Failed to disable Claude account", zap.Any("err", err))
				return err
			}
		}
	}

	// 更新用户信息
	his.UniqueName = user.UniqueName
	his.Password = user.Password
	his.Enable = user.Enable
	his.Openai = user.Openai
	his.OpenaiToken = user.OpenaiToken
	his.Claude = user.Claude
	his.ClaudeToken = user.ClaudeToken

	// 更新过期时间（如果有）
	if !user.ExpirationTime.IsZero() {
		his.ExpirationTime = user.ExpirationTime
	}
	his.UpdateTime = time.Now()

	// 保存更新后的用户信息
	if err := s.userRepository.Update(ctx, his); err != nil {
		s.logger.Error("Failed to update user", zap.Any("err", err))
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
