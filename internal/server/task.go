package server

import (
	"PandoraFuclaudePlusHelper/internal/model"
	"PandoraFuclaudePlusHelper/internal/repository"
	"PandoraFuclaudePlusHelper/internal/service"
	"PandoraFuclaudePlusHelper/internal/util"
	"PandoraFuclaudePlusHelper/pkg/log"
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"time"
)

type Task struct {
	log                     *log.Logger
	scheduler               *gocron.Scheduler
	openaiTokenRepository   repository.OpenaiTokenRepository
	openaiAccountRepository repository.OpenaiAccountRepository
	claudeTokenRepository   repository.ClaudeTokenRepository
	claudeAccountRepository repository.ClaudeAccountRepository
	userRepository          repository.UserRepository
	openaiAccountService    service.OpenaiAccountService
	claudeAccountService    service.ClaudeAccountService
}

func NewTask(log *log.Logger,
	openaiTokenRepository repository.OpenaiTokenRepository, openaiAccountRepository repository.OpenaiAccountRepository,
	claudeTokenRepository repository.ClaudeTokenRepository, claudeAccountRepository repository.ClaudeAccountRepository,
	userRepository repository.UserRepository,
	openaiAccountService service.OpenaiAccountService, claudeAccountService service.ClaudeAccountService,
) *Task {
	return &Task{
		log:                     log,
		openaiTokenRepository:   openaiTokenRepository,
		openaiAccountRepository: openaiAccountRepository,
		claudeTokenRepository:   claudeTokenRepository,
		claudeAccountRepository: claudeAccountRepository,
		userRepository:          userRepository,
		openaiAccountService:    openaiAccountService,
		claudeAccountService:    claudeAccountService,
	}
}

func (t *Task) RefreshAllToken(ctx context.Context) {
	t.log.Info("RefreshAllToken Start")
	tokens, err := t.openaiTokenRepository.GetAllToken(ctx)
	if err != nil {
		t.log.Error(fmt.Sprintf("RefreshAllToken GetAllToken error: %v", err))
	}
	if len(tokens) == 0 {
		t.log.Info("RefreshAllToken No token to refresh")
		return
	}
	t.log.Info(fmt.Sprintf("RefreshAllToken Token: %v", tokens))
	for _, token := range tokens {
		t.refreshAccessToken(ctx, token)
		t.refreshShareToken(ctx, token, false)
	}
	t.log.Info("RefreshAllToken Finish")
}

func (t *Task) refreshAccessToken(ctx context.Context, token *model.OpenaiToken) {
	t.log.Info(fmt.Sprintf("Refresh Token: %s", token.TokenName))
	// 刷新订阅状态
	plusSubscription := util.CheckSubscriptionStatus(token.AccessToken, t.log)
	token.PlusSubscription = plusSubscription

	now := time.Now()
	expireAt := token.ExpireAt
	later := now.Add(time.Hour * 1)
	if expireAt.After(later) {
		t.log.Info(fmt.Sprintf("Token not expired: %s", token.TokenName))
	} else {
		// 如果Token过期时间在1小时之内，刷新Token
		accessToken, expire, err := util.GenAccessToken(token.RefreshToken, t.log)
		if err != nil {
			t.log.Error(fmt.Sprintf("GenAccessToken error: %v", err))
		}
		token.AccessToken = accessToken
		token.ExpireAt = now.Add(time.Second * time.Duration(expire))
	}

	token.UpdateTime = now
	err := t.openaiTokenRepository.Update(ctx, token)
	if err != nil {
		t.log.Error(fmt.Sprintf("Update Token error: %v", err))
	}
}

func (t *Task) refreshShareToken(ctx context.Context, token *model.OpenaiToken, resetLimit bool) {
	accounts, err := t.openaiAccountRepository.SearchAccount(ctx, token.ID)
	if err != nil {
		t.log.Error(fmt.Sprintf("refreshShareToken SearchAccount error: %v", err))
	}
	if len(accounts) == 0 {
		t.log.Info(fmt.Sprintf("No account to refresh: %s", token.TokenName))
		return
	}
	for _, account := range accounts {
		if account.Status == 0 {
			t.log.Info(fmt.Sprintf("Account is disabled: %s", account.Account))
			continue
		}
		expireAt := account.ExpireAt
		later := time.Now().Add(time.Hour * 1)
		if expireAt.After(later) {
			t.log.Info(fmt.Sprintf("ShareToken not expired: %s", account.Account))
		} else {
			// 如果Token过期时间在1小时之内，刷新Token
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
				resetLimit,
				account.TemporaryChat == 1,
				t.log)
			if err != nil {
				t.log.Error(fmt.Sprintf("refreshShareToken GenerateShareToken error: %v", err))
				continue
			}
			account.ShareToken = shareToken
			account.ShareTokenEncrypt = shareTokenEncrypt
			account.ExpireAt = time.Unix(expireIn, 0)
		}
		now := time.Now()
		account.UpdateTime = now
		err = t.openaiAccountRepository.Update(ctx, account)
		if err != nil {
			t.log.Error(fmt.Sprintf("refreshShareToken Update error: %v", err))
		}
	}
}

func (t *Task) ResetLimit(ctx context.Context) {
	t.log.Info("ResetLimit Start")
	tokens, err := t.openaiTokenRepository.GetAllToken(ctx)
	if err != nil {
		t.log.Error(fmt.Sprintf("ResetLimit GetAllToken error: %v", err))
	}
	if len(tokens) == 0 {
		t.log.Info("ResetLimit No token to reset")
		return
	}
	t.log.Info(fmt.Sprintf("ResetLimit Token count: %d", len(tokens)))
	for _, token := range tokens {
		t.refreshShareToken(ctx, token, true)
	}
	t.log.Info("ResetLimit Finish")
}

func (t *Task) DisableUser(ctx context.Context) {
	t.log.Info("DisableUser Start")
	users, err := t.userRepository.GetAllUser(ctx)
	if err != nil {
		t.log.Error(fmt.Sprintf("GetAllUser error: %v", err))
		return
	}
	if len(users) == 0 {
		t.log.Info("No user to disable")
		return
	}
	t.log.Info(fmt.Sprintf("DisableUser Users count: %d", len(users)))

	now := time.Now()

	for _, user := range users {
		if user.Enable == 0 {
			t.log.Info(fmt.Sprintf("User already disabled: %s", user.UniqueName))
			continue
		}
		// 判断是否超过了有效期
		if user.ExpirationTime.After(now) {
			t.log.Info(fmt.Sprintf("User not yet expired: %s", user.UniqueName))
			continue
		}
		// 查询这个用户关联的openai账户
		openaiAccount, err := t.openaiAccountRepository.GetAccountByUserId(ctx, user.ID)
		if err != nil {
			t.log.Error(fmt.Sprintf("DisableUser GetAccountByUserId error: %v", err))
			continue
		}
		if openaiAccount != nil && openaiAccount.Status == 1 {
			// 禁用掉这个账户
			err = t.openaiAccountService.DisableAccount(ctx, openaiAccount.ID)
			if err != nil {
				t.log.Error(fmt.Sprintf("DisableOpenaiAccount error: %v", err))
			}
		}
		// 查询这个用户的 claude 账户
		claudeAccount, err := t.claudeAccountRepository.GetAccountByUserId(ctx, user.ID)
		if err != nil {
			t.log.Error(fmt.Sprintf("DisableUser GetAccountByUserId error: %v", err))
			continue
		}
		if claudeAccount != nil && claudeAccount.Status == 1 {
			// 禁用掉这个账户
			err = t.claudeAccountService.DisableAccount(ctx, claudeAccount.ID)
			if err != nil {
				t.log.Error(fmt.Sprintf("DisableClaudeAccount error: %v", err))
			}
		}
		user.Enable = 0
		user.Openai = 0
		user.Claude = 0
		user.UpdateTime = now
		err = t.userRepository.Update(ctx, user)
		if err != nil {
			t.log.Error(fmt.Sprintf("DisableUser Update error: %v", err))
		}
	}

	t.log.Info("DisableUser Finish")
}

func (t *Task) disableAndLogAccount(ctx context.Context, token *model.OpenaiToken, account *model.OpenaiAccount, now time.Time) error {
	_, _, _, err := util.GenShareToken(token.AccessToken,
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
		t.log)
	if err != nil {
		return err
	}
	account.Status = 0
	account.ExpireAt = now
	account.UpdateTime = now
	return t.openaiAccountRepository.Update(ctx, account)
}

func (t *Task) Start(ctx context.Context) error {
	gocron.SetPanicHandler(func(jobName string, recoverData interface{}) {
		t.log.Error(fmt.Sprintf("Task Panic job: %s, %v", jobName, recoverData))
	})

	t.scheduler = gocron.NewScheduler(time.UTC)

	_, err := t.scheduler.Cron("5 * * * *").Do(t.RefreshAllToken, ctx)
	if err != nil {
		t.log.Error(fmt.Sprintf("RefreshAllToken Task Start Error: %v", err))
	}

	_, err = t.scheduler.Cron("15 0 * * *").Do(t.ResetLimit, ctx)
	if err != nil {
		t.log.Error(fmt.Sprintf("ResetLimit Task Start Error: %v", err))
	}

	_, err = t.scheduler.Cron("2-59/5 * * * *").Do(t.DisableUser, ctx)
	if err != nil {
		t.log.Error(fmt.Sprintf("DisableUser Task Start Error: %v", err))
	}

	t.scheduler.StartBlocking()
	return nil
}

func (t *Task) Stop(ctx context.Context) error {
	t.scheduler.Stop()
	t.log.Info("Task stop...")
	return nil
}
