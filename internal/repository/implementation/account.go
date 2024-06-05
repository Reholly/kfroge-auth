package implementation

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"sso-service/internal/domain/entity"
	"sso-service/internal/domain/repository"
	"sso-service/internal/repository/model"
)

const (
	UniqueViolationCode = "23505"
)

type AccountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) repository.AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) FindByID(ctx context.Context, id uuid.UUID) (entity.Account, error) {
	query := `SELECT
				id,
				username, 
				email, 
				name,
				lastname, 
				patronymic,
				is_email_confirmed, 
				is_banned, 
				hashed_password, 
				role
				FROM account 
					WHERE id = $1`

	var account model.AccountModel
	err := r.chooseDomainErrorOrDefault(r.db.QueryRow(ctx, query, id).
		Scan(&account.ID,
			&account.Username,
			&account.Email,
			&account.Name,
			&account.Lastname,
			&account.Patronymic,
			&account.IsEmailConfirmed,
			&account.IsBanned,
			&account.HashedPassword,
			&account.Role),
	)

	if err != nil {
		return entity.Account{}, errors.Wrap(err, fmt.Sprintf("[AccountRepository] ошибка поиска по id: %s", id.String()))
	}

	return account.MapToEntity(), nil
}

func (r *AccountRepository) FindByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (entity.Account, error) {
	query := `SELECT
				id,
				username, 
				email, 
				name,
				lastname, 
				patronymic,
				is_email_confirmed, 
				is_banned, 
				hashed_password, 
				role
				FROM account 
					WHERE username = $1 OR email = $2`

	var account model.AccountModel
	err := r.chooseDomainErrorOrDefault(r.db.QueryRow(ctx, query, usernameOrEmail, usernameOrEmail).
		Scan(&account.ID,
			&account.Username,
			&account.Email,
			&account.Name,
			&account.Lastname,
			&account.Patronymic,
			&account.IsEmailConfirmed,
			&account.IsBanned,
			&account.HashedPassword,
			&account.Role),
	)

	if err != nil {
		return entity.Account{}, errors.Wrap(err, fmt.Sprintf("[AccountRepository] ошибка поиска по логину или почте: %s", usernameOrEmail))
	}

	return account.MapToEntity(), nil
}

func (r *AccountRepository) FindByUsername(ctx context.Context, username string) (entity.Account, error) {
	query := `SELECT
				id,
				username, 
				email, 
				name,
				lastname, 
				patronymic,
				is_email_confirmed, 
				is_banned, 
				hashed_password, 
				role
				FROM account 
					WHERE username = $1`

	var account model.AccountModel
	err := r.chooseDomainErrorOrDefault(r.db.QueryRow(ctx, query, username).
		Scan(&account.ID,
			&account.Username,
			&account.Email,
			&account.Name,
			&account.Lastname,
			&account.Patronymic,
			&account.IsEmailConfirmed,
			&account.IsBanned,
			&account.HashedPassword,
			&account.Role),
	)

	if err != nil {
		return entity.Account{}, errors.Wrap(err, fmt.Sprintf("[AccountRepository] ошибка поиска по логину: %s", username))
	}

	return account.MapToEntity(), nil
}

func (r *AccountRepository) FindByEmail(ctx context.Context, email string) (entity.Account, error) {
	query := `SELECT
				id,
				username, 
				emial, 
				name,
				lastname, 
				patronymic,
				is_email_confirmed, 
				is_banned, 
				hashed_password, 
				role
				FROM account 
					WHERE email = $1`

	var account model.AccountModel
	err := r.chooseDomainErrorOrDefault(r.db.QueryRow(ctx, query, &account, email).
		Scan(&account.ID,
			&account.Username,
			&account.Email,
			&account.Name,
			&account.Lastname,
			&account.Patronymic,
			&account.IsEmailConfirmed,
			&account.IsBanned,
			&account.HashedPassword,
			&account.Role),
	)

	if err != nil {
		return entity.Account{}, errors.Wrap(err, fmt.Sprintf("[AccountRepository] ошибка поиска по почте %s", email))
	}

	return account.MapToEntity(), nil
}

func (r *AccountRepository) UpdatePasswordByID(ctx context.Context, id uuid.UUID, hashedPassword string) error {
	sql := `UPDATE account 
				SET hashed_password = $1
					WHERE id = $2`

	_, err := r.db.Exec(ctx, sql, hashedPassword, id)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf(" [AccountRepository] ошибка обновления пароля аккаунта с id: %s", id.String()))
	}

	return nil
}

func (r *AccountRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	sql := `DELETE FROM account 
       			WHERE id = $1`

	_, err := r.db.Exec(ctx, sql, id)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf(" [AccountRepository] ошибка удаления аккаунта с id: %s", id.String()))
	}

	return nil
}

func (r *AccountRepository) Create(ctx context.Context, account entity.Account) error {
	sql := `INSERT INTO account(
                    username, 
                    email,
                    name,
                    lastname,
                    patronymic,
                    is_email_confirmed, 
                    is_banned,
                    hashed_password,
                    role)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(
		ctx,
		sql,
		account.Username,
		account.Email,
		account.Info.Name,
		account.Info.Lastname,
		account.Info.Patronymic,
		account.IsEmailConfirmed,
		account.IsBanned,
		account.HashedPassword,
		account.Role)

	err = r.chooseDomainErrorOrDefault(err)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("[AccountRepository] ошибка создания аккаунта с почтой %s, логином: %s",
			account.Email, account.Username))
	}

	return nil
}

func (r *AccountRepository) ChangeBanStatusByID(ctx context.Context, id uuid.UUID, status bool) error {
	sql := `UPDATE account
				SET is_banned = $1
					WHERE id = $2`

	_, err := r.db.Exec(ctx, sql, status, id)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("[AccountRepository] ошибка смены статуса бана аккаунта с id: %s на %t ", id.String(), status))
	}

	return nil
}

func (r *AccountRepository) ChangeRoleByID(ctx context.Context, id uuid.UUID, role string) error {
	sql := `UPDATE account 
				SET role = $1
					WHERE id = $2`

	_, err := r.db.Exec(ctx, sql, role, id)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("[AccountRepository] ошибка смены роли аккаунта с id: %s на %s ", id.String(), role))
	}

	return nil
}

func (r *AccountRepository) ConfirmEmailByID(ctx context.Context, id uuid.UUID) error {
	sql := `UPDATE account 
				SET is_email_confirmed = true
					WHERE id = $1`

	_, err := r.db.Exec(ctx, sql, id)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("[AccountRepository] ошибка подтверждения почты аккаунта с id: %s ", id.String()))
	}

	return nil
}

func (r *AccountRepository) UpdateInfoById(ctx context.Context, id uuid.UUID, info entity.AccountInfo) error {
	sql := `UPDATE account
				SET name = $1, 
				    lastname = $2, 
				    patronymic = $3
					WHERE id = $4`

	_, err := r.db.Exec(ctx, sql, info.Name, info.Lastname, info.Patronymic, id)

	err = r.chooseDomainErrorOrDefault(err)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("[AccountRepository] ошибка обновления главной информации аккаунта с id : %s", id.String()))
	}

	return nil
}

func (r *AccountRepository) chooseDomainErrorOrDefault(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch {
		case errors.Is(pgErr, pgx.ErrNoRows):
			return repository.ErrorAccountNotFound
		case pgErr.Code == UniqueViolationCode:
			return repository.ErrorAccountAlreadyExists
		}
	}

	return err
}
