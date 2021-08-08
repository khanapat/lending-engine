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

func (r lendingRepositoryDB) QueryWalletTransactionByIDRepo(ctx context.Context, id int) (*WalletTransaction, error) {
	var walletTransaction WalletTransaction
	err := r.db.GetContext(ctx, &walletTransaction, `
		SELECT	id,
				account_id,
				address,
				chain_id,
				txn_hash,
				collateral_type,
				volume,
				txn_type,
				status,
				created_datetime,
				updated_datetime
		FROM lending.public.wallet_transaction
		WHERE id = $1
	;`, id)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return &walletTransaction, nil
	}
}

func (r lendingRepositoryDB) QueryWalletTransactionRepo(ctx context.Context, request map[string]interface{}) (*[]WalletTransaction, error) {
	walletTransactions := make([]WalletTransaction, 0)
	query := `
		SELECT	id,
				account_id,
				address,
				chain_id,
				txn_hash,
				collateral_type,
				volume,
				txn_type,
				status,
				created_datetime,
				updated_datetime
		FROM lending.public.wallet_transaction
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
		var walletTransaction WalletTransaction
		if err := rows.StructScan(&walletTransaction); err != nil {
			return nil, err
		}
		walletTransactions = append(walletTransactions, walletTransaction)
	}
	defer rows.Close()
	return &walletTransactions, nil
}

func (r lendingRepositoryDB) InsertDepositRepo(ctx context.Context, accountId int, address string, chainId int, txnHash string, collateralType string, volume float64, txnType string, status string) (int64, error) {
	var depositId int64
	if err := r.db.QueryRowContext(ctx, `
		INSERT INTO lending.public.wallet_transaction
		(
			account_id,
			address,
			chain_id,
			txn_hash,
			collateral_type,
			volume,
			txn_type,
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
			$7,
			$8
		)
		RETURNING id
	;`, accountId, address, chainId, txnHash, collateralType, volume, txnType, status).Scan(&depositId); err != nil {
		return 0, err
	}
	return depositId, nil
}

func (r lendingRepositoryDB) UpdateDepositRepo(ctx context.Context, id int, status string, timestamp string) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE lending.public.wallet_transaction
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

func (r lendingRepositoryDB) InsertWithdrawRepo(ctx context.Context, accountId int, address string, chainId int, collateralType string, volume float64, txnType string, status string) (int64, error) {
	var withdrawId int64
	if err := r.db.QueryRowContext(ctx, `
		INSERT INTO lending.public.wallet_transaction
		(
			account_id,
			address,
			chain_id,
			collateral_type,
			volume,
			txn_type,
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
		RETURNING id
	;`, accountId, address, chainId, collateralType, volume, txnType, status).Scan(&withdrawId); err != nil {
		return 0, err
	}
	return withdrawId, nil
}

func (r lendingRepositoryDB) UpdateWithdrawRepo(ctx context.Context, id int, txnHash string, status string, timestamp string) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE lending.public.wallet_transaction
		SET 	status = $1,
				txn_hash = $2,
				updated_datetime = $3
		WHERE id = $4
	;`, status, txnHash, timestamp, id)
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
	interestTerms := make([]InterestTerm, 0)
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

func (r lendingRepositoryDB) InsertInterestTermRepo(ctx context.Context, interestRate float64) (int64, error) {
	var interestCode int64
	if err := r.db.QueryRowContext(ctx, `
		INSERT INTO lending.public.interest_term
		(
			interest_rate
		)
		VALUES
		(
			$1
		)
		RETURNING interest_code
	;`, interestRate).Scan(&interestCode); err != nil {
		return 0, err
	}
	return interestCode, nil
}

func (r lendingRepositoryDB) UpdateInterestTermRepo(ctx context.Context, code int, interestRate float64) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE lending.public.interest_term
		SET	interest_rate = $1
		WHERE interest_code = $2
	;`, interestRate, code)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}

func (r lendingRepositoryDB) QueryRepayTransactionByIDRepo(ctx context.Context, id int) (*RepayTransaction, error) {
	var repay RepayTransaction
	err := r.db.GetContext(ctx, &repay, `
		SELECT id, contract_id, account_id, amount, slip, status, created_datetime, updated_datetime
		FROM lending.public.repay_transaction
		WHERE id = $1
	;`, id)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return &repay, nil
	}
}

func (r lendingRepositoryDB) QueryRepayTransactionRepo(ctx context.Context, request map[string]interface{}) (*[]RepayTransaction, error) {
	repays := make([]RepayTransaction, 0)
	query := `
		SELECT id, contract_id, account_id, amount, slip, status, created_datetime, updated_datetime
		FROM lending.public.repay_transaction
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
		var repay RepayTransaction
		if err := rows.StructScan(&repay); err != nil {
			return nil, err
		}
		repays = append(repays, repay)
	}
	defer rows.Close()
	return &repays, nil
}

func (r lendingRepositoryDB) InsertRepayTransactionRepo(ctx context.Context, contractId int, accountId int, amount float64, slip string) (int64, error) {
	var repaymentId int64
	if err := r.db.QueryRowContext(ctx, `
		INSERT INTO lending.public.repay_transaction
		(
			contract_id,
			account_id,
			amount,
			slip
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4
		)
		RETURNING id
	;`, contractId, accountId, amount, slip).Scan(&repaymentId); err != nil {
		return 0, err
	}
	return repaymentId, nil
}

func (r lendingRepositoryDB) UpdateRepayTransactionRepo(ctx context.Context, repayId int, status string, timestamp string) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE lending.public.repay_transaction
		SET status = $1,
			updated_datetime = $2
		WHERE id = $3
	;`, status, timestamp, repayId)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}
