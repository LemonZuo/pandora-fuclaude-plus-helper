package repository

import (
	"PandoraFuclaudePlusHelper/internal/model"
	"context"
	"time"
)

type ClaudeTokenRepository interface {
	GetToken(ctx context.Context, id int64) (*model.ClaudeToken, error)
	Update(ctx context.Context, token *model.ClaudeToken) error
	Create(ctx context.Context, token *model.ClaudeToken) error
	SearchToken(ctx context.Context, keyword string) ([]*model.ClaudeToken, error)
	DeleteToken(ctx context.Context, id int64) error
	GetAllToken(ctx context.Context) ([]*model.ClaudeToken, error)
}

func NewClaudeTokenRepository(
	repository *Repository,

) ClaudeTokenRepository {
	return &claudeTokenRepository{
		Repository: repository,
	}
}

type claudeTokenRepository struct {
	*Repository
}

func (r *claudeTokenRepository) SearchToken(ctx context.Context, keyword string) ([]*model.ClaudeToken, error) {
	var tokens []*model.ClaudeToken
	if err := r.DB(ctx).Where("token_name like ?", "%"+keyword+"%").Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *claudeTokenRepository) Update(ctx context.Context, token *model.ClaudeToken) error {
	his, err := r.GetToken(ctx, token.ID)
	if err != nil {
		return err
	}
	token.CreateTime = his.CreateTime
	token.UpdateTime = time.Now()
	if err = r.DB(ctx).Save(token).Error; err != nil {
		return err
	}
	return nil
}

func (r *claudeTokenRepository) Create(ctx context.Context, token *model.ClaudeToken) error {
	if err := r.DB(ctx).Create(token).Error; err != nil {
		return err
	}
	return nil
}

func (r *claudeTokenRepository) DeleteToken(ctx context.Context, id int64) error {
	r.DB(ctx).Delete(&model.ClaudeToken{}, id)
	return nil
}

func (r *claudeTokenRepository) GetToken(ctx context.Context, id int64) (*model.ClaudeToken, error) {
	var token model.ClaudeToken
	if err := r.DB(ctx).Where("id = ?", id).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *claudeTokenRepository) GetAllToken(ctx context.Context) ([]*model.ClaudeToken, error) {
	var tokens []*model.ClaudeToken
	if err := r.DB(ctx).Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}
