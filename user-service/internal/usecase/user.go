package usecase

import (
	"context"
	"medods/user-service/internal/entity"
	"medods/user-service/internal/infrastructure/repository"
	"time"
)

type User interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Get(ctx context.Context, params map[string]string) (*entity.User, error)
	ListUsers(ctx context.Context, limit, offset int64, field, value string) ([]*entity.User, int64, error)
	ListDeletedUsers(ctx context.Context, limit, offset int64, field, value string) ([]*entity.User, int64, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	SoftDelete(ctx context.Context, id string) error

	CheckUniquess(ctx context.Context, field, value string) (int32, error)
	Exists(ctx context.Context, field, value string) (*entity.User, error)
}

type UserService struct {
	BaseUseCase
	repo       repository.User
	ctxTimeout time.Duration
}

func NewUserService(ctxTimeout time.Duration, repo repository.User) UserService {
	return UserService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (u UserService) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	u.beforeRequest(nil, &user.CreatedAt, &user.UpdatedAt, nil)

	return u.repo.Create(ctx, user)
}

func (u UserService) Get(ctx context.Context, params map[string]string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.Get(ctx, params)
}

func (u UserService) ListUsers(ctx context.Context, limit, offset int64, field, value string) ([]*entity.User, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.ListUsers(ctx, limit, offset, field, value)
}

func (u UserService) ListDeletedUsers(ctx context.Context, limit, offset int64, field, value string) ([]*entity.User, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.ListDeletedUsers(ctx, limit, offset, field, value)
}

func (u UserService) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	u.beforeRequest(nil, nil, &user.UpdatedAt, nil)

	return u.repo.Update(ctx, user)
}

func (u UserService) SoftDelete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	var user entity.User
	user.Id = id
	user.DeletedAt = time.Now().UTC()
	u.beforeRequest(nil, nil, nil, &user.DeletedAt)

	return u.repo.SoftDelete(ctx, user.Id)
}

func (u UserService) CheckUniquess(ctx context.Context, field, value string) (int32, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.CheckUniquess(ctx, field, value)
}

func (u UserService) Exists(ctx context.Context, field, value string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.Exists(ctx, field, value)
}
