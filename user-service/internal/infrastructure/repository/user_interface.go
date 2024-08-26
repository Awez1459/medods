package repository

import (
	"context"
	"medods/user-service/internal/entity"
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
