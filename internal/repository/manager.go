package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickmn/go-cache"
	"sso-service/internal/domain/repository"
	"sso-service/internal/repository/implementation"
)

type RepositoryManager struct {
	Code    repository.CodeRepository
	Account repository.AccountRepository
}

func NewRepositoryManager(conn *pgxpool.Pool, cache *cache.Cache, codeSalt string) RepositoryManager {
	return RepositoryManager{
		Code:    implementation.NewCodeRepository(conn, cache, codeSalt),
		Account: implementation.NewAccountRepository(conn),
	}
}
