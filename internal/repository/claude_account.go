package repository

import (
	"PandoraPlusHelper/internal/model"
	"context"
	"errors"
)

type ClaudeAccountRepository interface {
	GetAccount(ctx context.Context, id int64) (*model.ClaudeAccount, error)
	GetAccountByUserId(ctx context.Context, id int64) (*model.ClaudeAccount, error)
	Update(ctx context.Context, account *model.ClaudeAccount) error
	Create(ctx context.Context, account *model.ClaudeAccount) error
	SearchAccount(ctx context.Context, tokenId int64) ([]*model.ClaudeAccount, error)
	DeleteAccount(ctx context.Context, id int64) error
	GetAccountByPassword(ctx context.Context, password string) (model.ClaudeAccount, error)
	GetAccountById(ctx context.Context, id int64) (model.ClaudeAccount, error)
}

func NewClaudeAccountRepository(
	repository *Repository,
) ClaudeAccountRepository {
	return &claudeAccountRepository{
		Repository: repository,
	}
}

type claudeAccountRepository struct {
	*Repository
}

func (r *claudeAccountRepository) Update(ctx context.Context, account *model.ClaudeAccount) error {
	if err := r.DB(ctx).Save(account).Error; err != nil {
		return err
	}
	return nil
}

func (r *claudeAccountRepository) Create(ctx context.Context, account *model.ClaudeAccount) error {
	if err := r.DB(ctx).Create(account).Error; err != nil {
		return err
	}
	return nil
}

func (r *claudeAccountRepository) SearchAccount(ctx context.Context, tokenId int64) ([]*model.ClaudeAccount, error) {
	var accounts []*model.ClaudeAccount
	query := r.DB(ctx)

	// Apply filter if tokenId is provided
	if tokenId != 0 {
		query = query.Where("token_id = ?", tokenId)
	}

	// Retrieve the accounts based on the query built above
	err := query.Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *claudeAccountRepository) DeleteAccount(ctx context.Context, id int64) error {
	r.DB(ctx).Delete(&model.ClaudeAccount{}, id)
	return nil
}

func (r *claudeAccountRepository) GetAccount(ctx context.Context, id int64) (*model.ClaudeAccount, error) {
	var account model.ClaudeAccount
	if err := r.DB(ctx).Where("id = ?", id).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *claudeAccountRepository) GetAccountByUserId(ctx context.Context, id int64) (*model.ClaudeAccount, error) {
	var account model.ClaudeAccount
	if err := r.DB(ctx).Where("user_id = ?", id).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *claudeAccountRepository) GetAccountByPassword(ctx context.Context, password string) (model.ClaudeAccount, error) {
	var account model.ClaudeAccount
	if err := r.DB(ctx).Where("password = ? and status = 1", password).First(&account).Error; err != nil {
		return model.ClaudeAccount{}, errors.New("PLEASE CHECK YOUR PASSWORD")
	}
	return account, nil
}

func (r *claudeAccountRepository) GetAccountById(ctx context.Context, id int64) (model.ClaudeAccount, error) {
	var account model.ClaudeAccount
	if err := r.DB(ctx).Where("id = ?", id).First(&account).Error; err != nil {
		return model.ClaudeAccount{}, err
	}
	return account, nil
}
