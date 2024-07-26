package server

import (
	"PandoraPlusHelper/internal/model"
	"PandoraPlusHelper/internal/repository"
	"PandoraPlusHelper/internal/service"
	"PandoraPlusHelper/internal/util"
	"PandoraPlusHelper/pkg/log"
	"context"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
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
		t.log.Error("GetAllToken error", zap.Any("err", err))
	}
	if len(tokens) == 0 {
		t.log.Info("No token to refresh")
		return
	}
	t.log.Info("RefreshAllToken", zap.Int("token count", len(tokens)))
	for _, token := range tokens {
		t.refreshAccessToken(ctx, token)
		t.refreshShareToken(ctx, token)
	}
	t.log.Info("RefreshAllToken Finish")
}

func (t *Task) refreshAccessToken(ctx context.Context, token *model.OpenaiToken) {
	t.log.Info("Refresh Token", zap.String("token", token.TokenName))
	// 刷新订阅状态
	plusSubscription := util.CheckSubscriptionStatus(token.AccessToken, t.log)
	token.PlusSubscription = plusSubscription

	expireAt := token.ExpireAt
	later := time.Now().Add(time.Hour * 1)
	if expireAt.After(later) {
		t.log.Info("Token not expired", zap.String("token", token.TokenName))
	} else {
		// 如果Token过期时间在1小时之内，刷新Token
		accessToken, expire, err := util.GenAccessToken(token.RefreshToken, t.log)
		if err != nil {
			t.log.Error("GenAccessToken error", zap.Any("err", err))
		}
		token.AccessToken = accessToken
		token.ExpireAt = time.Unix(int64(expire), 0)
	}

	token.UpdateTime = time.Now()
	err := t.openaiTokenRepository.Update(ctx, token)
	if err != nil {
		t.log.Error("Update Token error", zap.Any("err", err))
	}
}

func (t *Task) refreshShareToken(ctx context.Context, token *model.OpenaiToken) {
	accounts, err := t.openaiAccountRepository.SearchAccount(ctx, token.ID)
	if err != nil {
		t.log.Error("refreshShareToken SearchAccount error", zap.Any("err", err))
	}
	if len(accounts) == 0 {
		t.log.Info("No account to refresh", zap.String("token", token.TokenName))
		return
	}
	for _, account := range accounts {
		if account.Status == 0 {
			t.log.Info("Account is disabled", zap.String("account", account.Account))
			continue
		}
		expireAt := account.ExpireAt
		later := time.Now().Add(time.Hour * 1)
		if expireAt.After(later) {
			t.log.Info("ShareToken not expired", zap.String("account", account.Account))
		} else {
			// 如果Token过期时间在1小时之内，刷新Token
			shareToken, expireIn, err := util.GenShareToken(token.AccessToken,
				account.Account,
				0,
				account.Gpt35Limit,
				account.Gpt4Limit,
				account.ShowConversations == 1,
				false,
				true,
				account.TemporaryChat == 1,
				t.log)
			if err != nil {
				t.log.Error("refreshShareToken GenerateShareToken error", zap.Any("err", err))
				continue
			}
			account.ShareToken = shareToken
			account.ExpireAt = time.Unix(expireIn, 0)
		}
		now := time.Now()
		account.UpdateTime = now
		err = t.openaiAccountRepository.Update(ctx, account)
		if err != nil {
			t.log.Error("refreshShareToken Update error", zap.Any("err", err))
		}
	}
}

func (t *Task) DisableUser(ctx context.Context) {
	t.log.Info("DisableUser Start")
	users, err := t.userRepository.GetAllUser(ctx)
	if err != nil {
		t.log.Error("GetAllUser error", zap.Any("err", err))
		return
	}
	if len(users) == 0 {
		t.log.Info("No user to disable")
		return
	}
	t.log.Info("DisableUser", zap.Int("user count", len(users)))

	now := time.Now()

	for _, user := range users {
		if user.Enable == 0 {
			t.log.Info("User already disabled", zap.String("user", user.UniqueName))
			continue
		}
		// 判断是否超过了有效期
		if user.ExpirationTime.After(now) {
			t.log.Info("User not yet expired", zap.String("user", user.UniqueName))
			continue
		}
		// 查询这个用户关联的openai账户
		openaiAccount, err := t.openaiAccountRepository.GetAccountByUserId(ctx, user.ID)
		if err != nil {
			t.log.Error("DisableUser GetAccountByUserId error", zap.Any("err", err))
			continue
		}
		if openaiAccount != nil && openaiAccount.Status == 1 {
			// 禁用掉这个账户
			err = t.openaiAccountService.DisableAccount(ctx, openaiAccount.ID)
			if err != nil {
				t.log.Error("DisableOpenaiAccount error", zap.Any("err", err))
			}
		}
		// 查询这个用户的 claude 账户
		claudeAccount, err := t.claudeAccountRepository.GetAccountByUserId(ctx, user.ID)
		if err != nil {
			t.log.Error("DisableUser GetAccountByUserId error", zap.Any("err", err))
			continue
		}
		if claudeAccount != nil && claudeAccount.Status == 1 {
			// 禁用掉这个账户
			err = t.claudeAccountService.DisableAccount(ctx, claudeAccount.ID)
			if err != nil {
				t.log.Error("DisableClaudeAccount error", zap.Any("err", err))
			}
		}
		user.Enable = 0
		user.Openai = 0
		user.Claude = 0
		user.UpdateTime = now
		err = t.userRepository.Update(ctx, user)
		if err != nil {
			t.log.Error("DisableUser Update error", zap.Any("err", err))
		}
	}

	t.log.Info("DisableUser Finish")
}

func (t *Task) disableAndLogAccount(ctx context.Context, token *model.OpenaiToken, account *model.OpenaiAccount, now time.Time) error {
	_, _, err := util.GenShareToken(token.AccessToken,
		account.Account,
		-1,
		account.Gpt35Limit,
		account.Gpt4Limit,
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
		t.log.Error("Task Panic", zap.String("job", jobName), zap.Any("recover", recoverData))
	})

	t.scheduler = gocron.NewScheduler(time.UTC)

	_, err := t.scheduler.Cron("0 * * * *").Do(t.RefreshAllToken, ctx)
	if err != nil {
		t.log.Error("RefreshAllToken Task Start Error", zap.Error(err))
	}

	_, err = t.scheduler.Cron("*/5 * * * *").Do(t.DisableUser, ctx)
	if err != nil {
		t.log.Error("DisableUser Task scheduling error", zap.String("task", "DisableUser"), zap.Error(err))
	}

	t.scheduler.StartBlocking()
	return nil
}

func (t *Task) Stop(ctx context.Context) error {
	t.scheduler.Stop()
	t.log.Info("Task stop...")
	return nil
}
