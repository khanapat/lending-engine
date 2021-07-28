package lending

import (
	"context"
	"time"
)

type ConfirmedDeposit struct {
	ID              *int       `db:"id" json:"id" example:"1"`
	AccountID       *int       `db:"account_id" json:"accountId" example:"1"`
	Address         *string    `db:"address" json:"address" example:"0xa9B6D99bA92D7d691c6EF4f49A1DC909822Cee46"`
	ChainID         *int       `db:"chain_id" json:"chainId" example:"1"`
	TxnHash         *string    `db:"txn_hash" json:"txnHash" example:"0xcbeafcd4c82144f7d1f9b94e4ed43e9ed1aa1434feb65a06fed97fee993ba075"`
	CollateralType  *string    `db:"collateral_type" json:"collateralType" example:"BTC"`
	Volume          *float64   `db:"volume" json:"volume" example:"0.5"`
	Status          *string    `db:"status" json:"status" example:"PENDING"`
	CreatedDatetime *time.Time `db:"created_datetime" json:"createdDatetime" example:"2021-01-02 12:13:14"`
	UpdatedDatetime *time.Time `db:"updated_datetime" json:"updatedDatetime" example:"2021-02-03 12:13:14"`
}

type Wallet struct {
	AccountID      *int       `db:"account_id" json:"accountId" example:"1"`
	BTCVolume      *float64   `db:"btc_volume" json:"btcVolume" example:"0.1"`
	ETHVolume      *float64   `db:"eth_volume" json:"ethVolume" example:"0.1"`
	MarginCallDate *string    `db:"margin_call_date" json:"marginCallDate" example:"2021-01-02"`
	LatestDatetime *time.Time `db:"latest_datetime" json:"latestDatetime" example:"2021-01-02 12:13:14"`
}

type Contract struct {
	ContractID      *int       `db:"contract_id" json:"contractId" example:"1"`
	AccountID       *int       `db:"account_id" json:"accountId" example:"1"`
	InterestCode    *int       `db:"interest_code" json:"interestCode" example:"1"`
	LoanOutstanding *float64   `db:"loan_outstanding" json:"loanOutstanding" example:"20000"`
	Term            *int       `db:"term" json:"term" example:"12"`
	Status          *string    `db:"status" json:"status" example:"CLOSED"`
	CreatedDatetime *time.Time `db:"created_datetime" json:"createdDatetime" example:"2021-01-02 12:13:14"`
	UpdatedDatetime *time.Time `db:"updated_datetime" json:"updatedDatetime" example:"2021-02-03 12:13:14"`
}

type InterestTerm struct {
	InterestCode *int     `db:"interest_code" json:"interestCode" example:"1"`
	InterestRate *float64 `db:"interest_rate" json:"interestRate" example:"0.05"`
}

type RepayTransaction struct {
	ID              *int       `db:"id" json:"id" example:"1"`
	ContractID      *int       `db:"contract_id" json:"contractId" example:"1"`
	AccountID       *int       `db:"account_id" json:"accountId" example:"1"`
	Amount          *float64   `db:"amount" json:"amount" example:"1000"`
	Slip            *string    `db:"slip" json:"slip" example:"<Base64>"`
	Status          *string    `db:"status" json:"status" example:"CONFIRMED"`
	CreatedDatetime *time.Time `db:"created_datetime" json:"createdDatetime" example:"2021-01-02 12:13:14"`
	UpdatedDatetime *time.Time `db:"updated_datetime" json:"updatedDatetime" example:"2021-02-03 12:13:14"`
}

type LendingRepository interface {
	QueryDepositByIDRepo(context.Context, int) (*ConfirmedDeposit, error)
	QueryDepositRepo(context.Context, map[string]interface{}) (*[]ConfirmedDeposit, error)
	InsertDepositRepo(context.Context, int, string, int, string, string, float64, string) (int64, error)
	UpdateDepositRepo(context.Context, int, string, string) (int64, error)
	QueryWalletRepo(context.Context, int) (*Wallet, error)
	UpdateWalletRepo(context.Context, int, float64, float64, *string, string) (int64, error)
	QueryContractByIDRepo(context.Context, int) (*Contract, error)
	QueryContractRepo(context.Context, map[string]interface{}) (*[]Contract, error)
	InsertContractRepo(context.Context, int, int, float64, int) (int64, error)
	UpdateContractRepo(context.Context, int, string, string) (int64, error)
	QueryInterestTermRepo(context.Context) (*[]InterestTerm, error)
	QueryRepayTransactionByIDRepo(context.Context, int) (*RepayTransaction, error)
	QueryRepayTransactionRepo(context.Context, map[string]interface{}) (*[]RepayTransaction, error)
	InsertRepayTransactionRepo(context.Context, int, int, float64, string) (int64, error)
	UpdateRepayTransactionRepo(context.Context, int, string, string) (int64, error)
}
