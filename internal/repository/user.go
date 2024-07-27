package repository

import (
	"PandoraFuclaudePlusHelper/internal/model"
	"context"
)

type UserRepository interface {
	GetUser(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Create(ctx context.Context, user *model.User) error
	SearchUser(ctx context.Context, keyword string) ([]*model.User, error)
	DeleteUser(ctx context.Context, id int64) error
	GetAllUser(ctx context.Context) ([]*model.User, error)
	GetUserByPassword(ctx context.Context, password string) (model.User, error)
}

func NewUserRepository(
	repository *Repository,

) UserRepository {
	return &userRepository{
		Repository: repository,
	}
}

type userRepository struct {
	*Repository
}

func (r *userRepository) SearchUser(ctx context.Context, keyword string) ([]*model.User, error) {
	var users []*model.User
	if err := r.DB(ctx).Where("unique_name like ?", "%"+keyword+"%").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.DB(ctx).Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.DB(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int64) error {
	r.DB(ctx).Delete(&model.User{}, id)
	return nil
}

func (r *userRepository) GetUser(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	if err := r.DB(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAllUser(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	if err := r.DB(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetUserByPassword(ctx context.Context, password string) (model.User, error) {
	var user model.User
	if err := r.DB(ctx).Where("password = ?", password).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}
