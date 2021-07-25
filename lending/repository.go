package lending

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type lendingRepositoryDB struct {
	db *sqlx.DB
}

func NewLendingRepositoryDB(db *sqlx.DB) lendingRepositoryDB {
	return lendingRepositoryDB{
		db: db,
	}
}

func (r lendingRepositoryDB) QueryDepositByAccountIDRepo(ctx context.Context, accountId int) (*[]ConfirmedDeposit, error) {
	var deposits []ConfirmedDeposit
	err := r.db.SelectContext(ctx, &deposits, `
		SELECT id, account_id, address, chain_id, txn_hash, collateral_type, volume, status, created_datetime, updated_datetime
		FROM lending.public.confirmed_deposit
		WHERE account_id = $1
	;`, accountId)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return &deposits, nil
	}
}

func (r lendingRepositoryDB) QueryDepositByIDRepo(ctx context.Context, id int) (*ConfirmedDeposit, error) {
	var deposit ConfirmedDeposit
	err := r.db.GetContext(ctx, &deposit, `
		SELECT id, account_id, address, chain_id, txn_hash, collateral_type, volume, status, created_datetime, updated_datetime
		FROM lending.public.confirmed_deposit
		WHERE id = $1
	;`, id)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return &deposit, nil
	}
}

func (r lendingRepositoryDB) QueryDepositRepo(ctx context.Context, request map[string]interface{}) (*[]ConfirmedDeposit, error) {
	deposits := make([]ConfirmedDeposit, 0)
	query := `
		SELECT	id,
				account_id,
				address,
				chain_id,
				txn_hash,
				collateral_type,
				volume,
				status,
				created_datetime,
				updated_datetime
		FROM lending.public.confirmed_deposit
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
		var deposit ConfirmedDeposit
		if err := rows.StructScan(&deposit); err != nil {
			return nil, err
		}
		deposits = append(deposits, deposit)
	}
	defer rows.Close()
	return &deposits, nil
}

func (r lendingRepositoryDB) InsertDepositRepo(ctx context.Context, accountId int, address string, chainId int, txnHash string, collateralType string, volume float64, status string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO lending.public.confirmed_deposit
		(
			account_id,
			address,
			chain_id,
			txn_hash,
			collateral_type,
			volume,
			status
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7
		)
	;`, accountId, address, chainId, txnHash, collateralType, volume, status)
	if err != nil {
		return err
	}
	return nil
}

func (r lendingRepositoryDB) UpdateDepositRepo(ctx context.Context, id int, status string, timestamp string) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE lending.public.confirmed_deposit
		SET 	status = $1,
				updated_datetime = $2
		WHERE id = $3
	;`, status, timestamp, id)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}

func (r lendingRepositoryDB) QueryWalletRepo(ctx context.Context, accountId int) (*Wallet, error) {
	var wallet Wallet
	err := r.db.GetContext(ctx, &wallet, `
		SELECT account_id, btc_volume, eth_volume, margin_call_date, latest_datetime
		FROM lending.public.wallet
		WHERE account_id = $1
	;`, accountId)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return &wallet, nil
	}
}

func (r lendingRepositoryDB) UpdateWalletRepo(ctx context.Context, accountId int, btc float64, eth float64, margin *string, latest string) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE lending.public.wallet
		SET btc_volume = $1,
			eth_volume = $2,
			margin_call_date = $3,
			latest_datetime = $4
		WHERE account_id = $5
	;`, btc, eth, margin, latest, accountId)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}

func (r lendingRepositoryDB) QueryContractByAccountIDRepo(ctx context.Context, accountId int) (*[]Contract, error) {
	var contracts []Contract
	err := r.db.SelectContext(ctx, &contracts, `
		SELECT contract_id, account_id, interest_code, loan_outstanding, term, status, created_datetime, updated_datetime
		FROM lending.public.contract
		WHERE account_id = $1
	;`, accountId)
	switch {
	case err == sql.ErrNoRows:
		return &contracts, nil
	case err != nil:
		return nil, err
	default:
		return &contracts, nil
	}
}

func (r lendingRepositoryDB) QueryContractByIDRepo(ctx context.Context, id int) (*Contract, error) {
	var contract Contract
	err := r.db.GetContext(ctx, &contract, `
		SELECT contract_id, account_id, interest_code, loan_outstanding, term, status, created_datetime, updated_datetime
		FROM lending.public.contract
		WHERE contract_id = $1
	;`, id)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return &contract, nil
	}
}

func (r lendingRepositoryDB) QueryContractRepo(ctx context.Context, request map[string]interface{}) (*[]Contract, error) {
	contracts := make([]Contract, 0)
	query := `
		SELECT contract_id, account_id, interest_code, loan_outstanding, term, status, created_datetime, updated_datetime
		FROM lending.public.contract
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
		var contract Contract
		if err := rows.StructScan(&contract); err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}
	defer rows.Close()
	return &contracts, nil
}

func (r lendingRepositoryDB) InsertContractRepo(ctx context.Context, accountId int, interestCode int, loan float64, term int) (int64, error) {
	var contractId int64
	if err := r.db.QueryRowContext(ctx, `
		INSERT INTO lending.public.contract
		(
			account_id,
			interest_code,
			loan_outstanding,
			term
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4
		)
		RETURNING contract_id
	;`, accountId, interestCode, loan, term).Scan(&contractId); err != nil {
		return 0, err
	}
	return contractId, nil
}

func (r lendingRepositoryDB) UpdateContractRepo(ctx context.Context, contractId int, status string, timestamp string) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE lending.public.contract
		SET		status = $1,
				updated_datetime = $2
		WHERE contract_id = $3
	;`, status, timestamp, contractId)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}

func (r lendingRepositoryDB) QueryInterestTermRepo(ctx context.Context) (*[]InterestTerm, error) {
	var interestTerms []InterestTerm
	err := r.db.SelectContext(ctx, &interestTerms, `
		SELECT interest_code, interest_rate
		FROM lending.public.interest_term
	;`)
	switch {
	case err == sql.ErrNoRows:
		return &interestTerms, nil
	case err != nil:
		return nil, err
	default:
		return &interestTerms, nil
	}
}
