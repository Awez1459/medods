package refresh_token

import (
	"context"
	"medods/api-service/internal/entity"
	"time"

	"github.com/google/uuid"
)

type refreshTokenService struct {
	ctxTimeout time.Duration
	repo       RefreshTokenRepo
}

func NewRefreshTokenService(ctxTimeout time.Duration, repo RefreshTokenRepo) RefreshToken {
	return &refreshTokenService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (r *refreshTokenService) beforeCreate(m *entity.RefreshToken) error {
	m.GUID = uuid.New().String()
	m.CreatedAt = time.Now().UTC()
	return nil
}

func (r *refreshTokenService) Get(ctx context.Context, refreshToken string) (*entity.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	return r.repo.Get(ctx, refreshToken)
}

func (r *refreshTokenService) Create(ctx context.Context, m *entity.RefreshToken) error {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	r.beforeCreate(m)
	return r.repo.Create(ctx, m)
}

func (r *refreshTokenService) Delete(ctx context.Context, refreshToken string) error {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	return r.repo.Delete(ctx, refreshToken)
}

func (r *refreshTokenService) GenerateToken(ctx context.Context, sub, tokenType, jwtSecret string, accessTTL, refreshTTL time.Duration, optionalFields ...map[string]interface{}) (string, string, error) {
	return "", "", nil
}
