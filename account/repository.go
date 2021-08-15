package account

import (
	"context"
	"database/sql"
	"fmt"

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

func (r accountRepositoryDB) GetAccountByEmailRepo(ctx context.Context, email string) (*Account, error) {
	var account Account
	err := r.db.GetContext(ctx, &account, `
		SELECT account_id, first_name, last_name, phone, email, "password", account_number, is_verify, status
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

func (r accountRepositoryDB) GetAccountByIDRepo(ctx context.Context, accountId int) (*Account, error) {
	var account Account
	err := r.db.GetContext(ctx, &account, `
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

func (r accountRepositoryDB) GetAccountRepo(ctx context.Context, request map[string]interface{}) (*[]AccountDetail, error) {
	details := make([]AccountDetail, 0)
	query := `
		SELECT	x.account_id,
				x.first_name,
				x.last_name,
				x.phone,
				x.email,
				x."password",
				x.account_number,
				x.is_verify,
				x.status,
				y.current_accept_version,
				z.document_id,
				i.document_type,
				z.file_context,
				z.file_name,
				z.tag 
		FROM lending.public.account x 
		INNER JOIN lending.public.terms_condition y ON x.account_id = y.account_id
		INNER JOIN lending.public.account_document z ON x.account_id = z.account_id
		INNER JOIN lending.public.document_info i ON z.document_id = i.document_id
		WHERE 1 = 1
	`
	for key, _ := range request {
		query = fmt.Sprintf("%s AND %s = :%s", query, key, key)
	}
	rows, err := r.db.NamedQueryContext(ctx, query, request)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var detail AccountDetail
		if err := rows.StructScan(&detail); err != nil {
			return nil, err
		}
		details = append(details, detail)
	}
	defer rows.Close()
	return &details, nil
}

func (r accountRepositoryDB) UpdateAccountRepo(ctx context.Context, accountId int, status string) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE lending.public.account
		SET		status = $1
		WHERE account_id = $2
	;`, status, accountId)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}

func (r accountRepositoryDB) UpdateAccountDocumentRepo(ctx context.Context, accountId int, documentId int, fileName string, fileContext string) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE lending.public.account_document
		SET		file_name = $1,
				file_context = $2
		WHERE	account_id = $3
		AND		document_id = $4
	;`, fileName, fileContext, accountId, documentId)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
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
	err := r.db.GetContext(ctx, &termsCondition, `
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

func (r accountRepositoryDB) ConfirmChangePasswordRepo(ctx context.Context, accountId int, password string) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE lending.public.account
		SET password = $1
		WHERE account_id = $2
	;`, password, accountId)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affect, nil
}

func (r accountRepositoryDB) QueryDocumentInfoAdminRepo(ctx context.Context) (*[]DocumentInfo, error) {
	documentInfos := make([]DocumentInfo, 0)
	err := r.db.SelectContext(ctx, &documentInfos, `
		SELECT document_id, document_type
		FROM lending.public.document_info
	;`)
	switch {
	case err == sql.ErrNoRows:
		return &documentInfos, nil
	case err != nil:
		return nil, err
	default:
		return &documentInfos, nil
	}
}

func (r accountRepositoryDB) InsertDocumentInfoAdminRepo(ctx context.Context, docType string) (int64, error) {
	var documentId int64
	if err := r.db.QueryRowContext(ctx, `
		INSERT INTO lending.public.document_info
		(
			document_type
		)
		VALUES
		(
			$1
		)
	;`, docType).Scan(&documentId); err != nil {
		return 0, err
	}
	return documentId, nil
}

func (r accountRepositoryDB) UpdateDocumentInfoAdminRepo(ctx context.Context, docId int, docType string) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE lending.public.document_info
		SET document_type = $1
		WHERE document_id = $2
	;`, docType, docId)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}

func (r accountRepositoryDB) CreateUserSubscriptionRepo(ctx context.Context, first string, last string, phone string, email string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO lending.public.user_subscription
		(
			first_name,
			last_name,
			phone,
			email
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4
		)
	;`, first, last, phone, email)
	if err != nil {
		return err
	}
	return nil
}
