package repository

import (
	"PandoraPlusHelper/internal/model"
	"context"
)

type OpenaiTokenRepository interface {
	GetToken(ctx context.Context, id int64) (*model.OpenaiToken, error)
	Update(ctx context.Context, token *model.OpenaiToken) error
	Create(ctx context.Context, token *model.OpenaiToken) error
	SearchToken(ctx context.Context, keyword string) ([]*model.OpenaiToken, error)
	DeleteToken(ctx context.Context, id int64) error
	GetAllToken(ctx context.Context) ([]*model.OpenaiToken, error)
}

func NewOpenaiTokenRepository(
	repository *Repository,

) OpenaiTokenRepository {
	return &openaiTokenRepository{
		Repository: repository,
	}
}

type openaiTokenRepository struct {
	*Repository
}

func (r *openaiTokenRepository) SearchToken(ctx context.Context, keyword string) ([]*model.OpenaiToken, error) {
	var tokens []*model.OpenaiToken
	if err := r.DB(ctx).Where("token_name like ?", "%"+keyword+"%").Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *openaiTokenRepository) Update(ctx context.Context, token *model.OpenaiToken) error {
	if err := r.DB(ctx).Save(token).Error; err != nil {
		return err
	}
	return nil
}

func (r *openaiTokenRepository) Create(ctx context.Context, token *model.OpenaiToken) error {
	if err := r.DB(ctx).Create(token).Error; err != nil {
		return err
	}
	return nil
}

func (r *openaiTokenRepository) DeleteToken(ctx context.Context, id int64) error {
	r.DB(ctx).Delete(&model.OpenaiToken{}, id)
	return nil
}

func (r *openaiTokenRepository) GetToken(ctx context.Context, id int64) (*model.OpenaiToken, error) {
	var token model.OpenaiToken
	if err := r.DB(ctx).Where("id = ?", id).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *openaiTokenRepository) GetAllToken(ctx context.Context) ([]*model.OpenaiToken, error) {
	var tokens []*model.OpenaiToken
	if err := r.DB(ctx).Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}
