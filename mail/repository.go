package mail

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Account struct {
	AccountID     *int    `db:"account_id"`
	FirstName     *string `db:"first_name"`
	LastName      *string `db:"last_name"`
	Phone         *string `db:"phone"`
	Email         *string `db:"email"`
	Password      *string `db:"password"`
	AccountNumber *string `db:"account_number"`
	IsVerify      *bool   `db:"is_verify"`
	Status        *string `db:"status"`
}

type QueryAccountByIDFn func(ctx context.Context, accountId int) (*Account, error)

func NewQueryAccountByIDFn(db *sqlx.DB) QueryAccountByIDFn {
	return func(ctx context.Context, accountId int) (*Account, error) {
		var account Account
		err := db.GetContext(ctx, &account, `
			SELECT account_id, first_name, last_name, phone, email, "password", account_number, is_verify, status
			FROM lending.public.account
			WHERE account_id = $1
		;`, accountId)
		switch {
		case err == sql.ErrNoRows:
			return nil, nil
		case err != nil:
			return nil, err
		default:
			return &account, nil
		}
	}
}
