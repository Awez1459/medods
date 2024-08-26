package postgresql

import (
	"medods/api-service/internal/pkg/postgres"
	"medods/api-service/internal/usecase/refresh_token"
)

func NewRefreshTokenRepo(db *postgres.PostgresDB) refresh_token.RefreshTokenRepo {
	return nil
}
