package account

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type accountRepositoryDB struct {
	db *sqlx.DB
}

func NewAccountRepositoryDB(db *sqlx.DB) accountRepositoryDB {
	return accountRepositoryDB{
		db: db,
	}
}

func (r accountRepositoryDB) SignUpAccountRepo(ctx context.Context, firstName string, lastName string, phone string, email string, password string, accountNumber string, citizenName string, citizenContext string, bookBankName string, bookBankContext string) (int64, error) {
	var accountId int64
	tx := r.db.MustBeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO lending.public.account
		(
			first_name,
			last_name,
			phone,
			email,
			"password",
			account_number
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)
		RETURNING account_id
	;`, firstName, lastName, phone, email, password, accountNumber).Scan(&accountId); err != nil {
		return 0, err
	}

	tx.MustExecContext(ctx, `
		INSERT INTO lending.public.terms_condition
		(
			account_id,
			current_accept_version
		)
		VALUES
		(
			$1,
			$2
		)
	;`, accountId, "0.0.0")

	insertDocument := `
		INSERT INTO lending.public.account_document
		(
			account_id,
			document_id,
			file_name,
			file_context,
			tag
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4,
			$5
		)
	;`

	tx.MustExecContext(ctx, insertDocument, accountId, 1, citizenName, citizenContext, "ID")
	tx.MustExecContext(ctx, insertDocument, accountId, 2, bookBankName, bookBankContext, "Ledger")
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return accountId, nil
}
func (r accountRepositoryDB) CreateWalletRepo(ctx context.Context, accountId int) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO lending.public.wallet
		(
			account_id,
			btc_volume,
			eth_volume
		)
		VALUES
		(
			$1,
			$2,
			$3
		)
	;`, accountId, 0, 0)
	if err != nil {
		return err
	}
	return nil
}

func (r accountRepositoryDB) ConfirmVerifyEmailRepo(ctx context.Context, accountId int) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE lending.public.account
		SET is_verify = true
		WHERE account_id = $1
	;`, accountId)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affect, nil
}

func (r accountRepositoryDB) GetTermsConditionRepo(ctx context.Context, accountId int) (*TermsCondition, error) {
	var termsCondition TermsCondition
	err := r.db.Get(&termsCondition, `
		SELECT account_id, current_accept_version
		FROM public.terms_condition
		WHERE account_id = $1
	;`, accountId)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return &termsCondition, nil
	}
}

func (r accountRepositoryDB) AcceptTermsConditionRepo(ctx context.Context, accountId int, version string) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE lending.public.terms_condition
		SET current_accept_version = $1
		WHERE account_id = $2
	;`, version, accountId)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affect, nil
}

func (r accountRepositoryDB) GetAccountRepo(ctx context.Context, email string) (*Account, error) {
	var account Account
	err := r.db.Get(&account, `
		SELECT account_id, first_name, last_name, phone, email, "password", account_number, is_verify
		FROM lending.public.account
		WHERE email = $1
	;`, email)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return &account, nil
	}
}
