package repository

import (
	"PandoraPlusHelper/internal/model"
	"context"
	"errors"
)

type OpenaiAccountRepository interface {
	GetAccount(ctx context.Context, id int64) (*model.OpenaiAccount, error)
	GetAccountByUserId(ctx context.Context, id int64) (*model.OpenaiAccount, error)
	Update(ctx context.Context, account *model.OpenaiAccount) error
	Create(ctx context.Context, account *model.OpenaiAccount) error
	SearchAccount(ctx context.Context, tokenId int64) ([]*model.OpenaiAccount, error)
	DeleteAccount(ctx context.Context, id int64) error
	GetAccountByPassword(ctx context.Context, password string) (model.OpenaiAccount, error)
	GetAccountById(ctx context.Context, id int64) (model.OpenaiAccount, error)
}

func NewOpenaiAccountRepository(
	repository *Repository,
) OpenaiAccountRepository {
	return &openaiAccountRepository{
		Repository: repository,
	}
}

type openaiAccountRepository struct {
	*Repository
}

func (r *openaiAccountRepository) Update(ctx context.Context, account *model.OpenaiAccount) error {
	if err := r.DB(ctx).Save(account).Error; err != nil {
		return err
	}
	return nil
}

func (r *openaiAccountRepository) Create(ctx context.Context, account *model.OpenaiAccount) error {
	if err := r.DB(ctx).Create(account).Error; err != nil {
		return err
	}
	return nil
}

func (r *openaiAccountRepository) SearchAccount(ctx context.Context, tokenId int64) ([]*model.OpenaiAccount, error) {
	var accounts []*model.OpenaiAccount
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

func (r *openaiAccountRepository) DeleteAccount(ctx context.Context, id int64) error {
	r.DB(ctx).Delete(&model.OpenaiAccount{}, id)
	return nil
}

func (r *openaiAccountRepository) GetAccount(ctx context.Context, id int64) (*model.OpenaiAccount, error) {
	var account model.OpenaiAccount
	if err := r.DB(ctx).Where("id = ?", id).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *openaiAccountRepository) GetAccountByUserId(ctx context.Context, id int64) (*model.OpenaiAccount, error) {
	var account model.OpenaiAccount
	if err := r.DB(ctx).Where("user_id = ?", id).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *openaiAccountRepository) GetAccountByPassword(ctx context.Context, password string) (model.OpenaiAccount, error) {
	var account model.OpenaiAccount
	if err := r.DB(ctx).Where("password = ? and status = 1", password).First(&account).Error; err != nil {
		return model.OpenaiAccount{}, errors.New("PLEASE CHECK YOUR PASSWORD")
	}
	return account, nil
}

func (r *openaiAccountRepository) GetAccountById(ctx context.Context, id int64) (model.OpenaiAccount, error) {
	var account model.OpenaiAccount
	if err := r.DB(ctx).Where("id = ?", id).First(&account).Error; err != nil {
		return model.OpenaiAccount{}, err
	}
	return account, nil
}
