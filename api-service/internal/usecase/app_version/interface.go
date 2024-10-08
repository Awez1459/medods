package app_version

import (
	"context"

	"medods/api-service/internal/entity"
)

type AppVersion interface {
	Get(ctx context.Context) (*entity.AppVersion, error)
	Create(ctx context.Context, m *entity.AppVersion) error
	Update(ctx context.Context, m *entity.AppVersion) error
}