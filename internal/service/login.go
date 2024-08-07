package service

import (
	v1 "PandoraFuclaudePlusHelper/api/v1"
	commonConfig "PandoraFuclaudePlusHelper/config"
	"PandoraFuclaudePlusHelper/internal/model"
	"PandoraFuclaudePlusHelper/internal/repository"
	"PandoraFuclaudePlusHelper/internal/util"
	"context"
	"errors"
	"fmt"
	"time"
)

type LoginService interface {
	Login(ctx context.Context, req *v1.LoginRequest) (int, string, map[string]interface{}, string, error)
}

func NewLoginService(service *Service, userRepository repository.UserRepository,
	openaiTokenRepository repository.OpenaiTokenRepository, openaiAccountRepository repository.OpenaiAccountRepository,
	claudeTokenRepository repository.ClaudeTokenRepository, claudeAccountRepository repository.ClaudeAccountRepository) LoginService {
	return &loginService{
		Service:                 service,
		userRepository:          userRepository,
		openaiTokenRepository:   openaiTokenRepository,
		openaiAccountRepository: openaiAccountRepository,
		claudeTokenRepository:   claudeTokenRepository,
		claudeAccountRepository: claudeAccountRepository,
	}
}

type loginService struct {
	*Service
	userRepository          repository.UserRepository
	openaiTokenRepository   repository.OpenaiTokenRepository
	openaiAccountRepository repository.OpenaiAccountRepository
	claudeTokenRepository   repository.ClaudeTokenRepository
	claudeAccountRepository repository.ClaudeAccountRepository
}

func (s *loginService) Login(ctx context.Context, req *v1.LoginRequest) (int, string, map[string]interface{}, string, error) {
	// 登录类型
	loginType := req.Type
	// tokenId或者accountId，用于后台快捷登录
	accountId := req.AccountId
	// 密码：前台登录时所需
	password := req.Password

	if accountId <= 0 && len(password) == 0 {
		return -1, "", nil, "", errors.New("缺少登录参数")
	}

	// 根据登录类型处理请求
	switch loginType {
	case 9999:
		// 管理员登录
		adminPassword := commonConfig.GetConfig().AdminPassword
		if password == adminPassword {
			login := map[string]interface{}{
				"id":          1,
				"loginname":   "admin",
				"email":       "admin@uasm.com",
				"role":        model.ADMIN_ROLE,
				"status":      1,
				"permissions": model.PERMISSION_LIST,
			}
			token, err := s.jwt.GenToken("1", time.Now().Add(time.Hour*24*90))
			if err != nil {
				return -1, "", nil, "", err
			}
			return loginType, token, login, "", nil
		}
		return -1, "", nil, "", v1.ErrLoginFailed
	case 1:
		// 普通用户openai登录
		user, err := s.userRepository.GetUserByPassword(ctx, password)
		if err != nil {
			s.logger.Info(fmt.Sprintf("user %s login failed", password))
			return -1, "", nil, "", v1.ErrLoginFailed
		}
		if user.Enable != 1 {
			s.logger.Info(fmt.Sprintf("user %s is not enable", user.UniqueName))
			return -1, "", nil, "", errors.New("登录失败")
		}
		account, err := s.openaiAccountRepository.GetAccountByUserId(ctx, user.ID)
		if err != nil {
			s.logger.Info(fmt.Sprintf("user %s has no account", user.UniqueName))
			return -1, "", nil, "", errors.New("登录失败")
		}
		if account.Status != 1 {
			s.logger.Info(fmt.Sprintf("account %d is not enable", account.ID))
			return -1, "", nil, "", errors.New("登录失败")
		}
		return gptLogin(account.ShareToken, s, loginType)
	case 2:
		// 管理员 openai快捷登录
		account, err := s.openaiAccountRepository.GetAccountById(ctx, accountId)
		if err != nil {
			return -1, "", nil, "", errors.New("账号不存在")
		}
		return gptLogin(account.ShareToken, s, loginType)
	case 3:
		// 普通用户claude登录
		user, err := s.userRepository.GetUserByPassword(ctx, password)
		if err != nil {
			s.logger.Info(fmt.Sprintf("user %s login failed", password))
			return -1, "", nil, "", v1.ErrLoginFailed
		}
		if user.Enable != 1 {
			s.logger.Info(fmt.Sprintf("user %s is not enable", user.UniqueName))
			return -1, "", nil, "", errors.New("登录失败")
		}
		account, err := s.claudeAccountRepository.GetAccountByUserId(ctx, user.ID)
		if err != nil {
			s.logger.Info(fmt.Sprintf("user %s has no account", user.UniqueName))
			return -1, "", nil, "", errors.New("登录失败")
		}
		if account.Status != 1 {
			s.logger.Info(fmt.Sprintf("account %d is not enable", account.ID))
			return -1, "", nil, "", errors.New("登录失败")
		}
		token, err := s.claudeTokenRepository.GetToken(ctx, user.ClaudeToken)
		if err != nil {
			s.logger.Info(fmt.Sprintf("user %s has no token", user.UniqueName))
			return -1, "", nil, "", errors.New("登录失败")
		}
		return claudeLogin(token.SessionToken, user.UniqueName, s, 3)
	case 4:
		// 管理员 claud account 快捷登录
		account, err := s.claudeAccountRepository.GetAccountById(ctx, accountId)
		if err != nil {
			return -1, "", nil, "", v1.ErrLoginFailed
		}
		token, err := s.claudeTokenRepository.GetToken(ctx, account.TokenID)
		if err != nil {
			return -1, "", nil, "", errors.New("账号不存在")
		}
		return claudeLogin(token.SessionToken, account.Account, s, 4)
	case 5:
		// 管理员 claud token 快捷登录
		token, err := s.claudeTokenRepository.GetToken(ctx, accountId)
		if err != nil {
			return -1, "", nil, "", errors.New("账号不存在")
		}
		return claudeLogin(token.SessionToken, "", s, 5)
	default:
		// 不支持的登录类型
		return -1, "", nil, "", v1.ErrLoginFailed
	}
}

func gptLogin(shareToken string, s *loginService, loginType int) (int, string, map[string]interface{}, string, error) {
	loginUrl, err := util.ExecuteShareAuth(shareToken, s.logger)
	if err != nil {
		return -1, "", nil, "", v1.ErrLoginFailed
	}
	return loginType, "", nil, loginUrl, nil
}

func claudeLogin(sessionToken string, account string, s *loginService, loginType int) (int, string, map[string]interface{}, string, error) {
	if (loginType == 3 && account == "") || (loginType == 4 && account == "") {
		// 通过 account登录时，account不能为空
		return -1, "", nil, "", v1.ErrLoginFailed
	}
	loginUrl, err := util.ExecuteClaudeAuth(sessionToken, account, s.logger)
	if err != nil {
		return -1, "", nil, "", v1.ErrLoginFailed
	}
	return loginType, "", nil, loginUrl, nil
}
