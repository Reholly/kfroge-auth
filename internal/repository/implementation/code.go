package implementation

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"sso-service/internal/domain/repository"
	"sso-service/internal/repository/model"
	"sso-service/lib/code"
)

type CodeRepository struct {
	db        *pgxpool.Pool
	codeCache *cache.Cache
	codeSalt  string
}

func NewCodeRepository(db *pgxpool.Pool, codeCache *cache.Cache, codeSalt string) repository.CodeRepository {
	return &CodeRepository{
		db:        db,
		codeCache: codeCache,
		codeSalt:  codeSalt,
	}
}

func (r *CodeRepository) CreatePasswordResetCode(ctx context.Context, account uuid.UUID) (string, error) {
	resetPasswordCode := code.GenerateCodeFromSalt(r.codeSalt)
	r.codeCache.Set(resetPasswordCode, account, cache.DefaultExpiration)

	return resetPasswordCode, nil
}

func (r *CodeRepository) CreateEmailConfirmationCode(ctx context.Context, account uuid.UUID) (string, error) {
	emailConfirmationCode := code.GenerateCodeFromSalt(r.codeSalt)

	sql := `INSERT INTO code(account_id, value)
				VALUES ($1, $2)`

	_, err := r.db.Exec(ctx, sql, account, emailConfirmationCode)

	if err != nil {
		return "", errors.Wrap(err, "[CodeRepository] ошибка создания кода подтверждения почты")
	}

	return emailConfirmationCode, nil
}

func (r *CodeRepository) GetAccountIDByResetPasswordCode(ctx context.Context, code string) (uuid.UUID, error) {
	resetPasswordCode, ok := r.codeCache.Get(code)
	if !ok {
		return uuid.UUID{}, repository.ErrorInvalidResetPasswordCode
	}

	id, ok := resetPasswordCode.(uuid.UUID)
	if !ok {
		return uuid.UUID{}, repository.ErrorInvalidResetPasswordCodeValue
	}

	return id, nil
}

func (r *CodeRepository) GetAccountIDByEmailConfirmationCode(ctx context.Context, code string) (uuid.UUID, error) {
	query := `SELECT account_id
				FROM code
					WHERE value = $1`

	var codeModel model.Code
	err := r.db.QueryRow(ctx, query, code).Scan(&codeModel.AccountId)

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch {
		case errors.Is(pgErr, pgx.ErrNoRows):
			return uuid.UUID{}, repository.ErrorNotFoundAccountForCode
		}
	}

	if err != nil {
		return uuid.UUID{}, errors.Wrap(err, fmt.Sprintf("[CodeRepository] ошибка получения айди аккаунта по коду %s", code))
	}

	return codeModel.AccountId, nil
}
