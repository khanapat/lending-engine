package account

import "context"

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

type TermsCondition struct {
	AccountID            *int    `db:"account_id" json:"accountId"`
	CurrentAcceptVersion *string `db:"current_accept_version" json:"currentAcceptVersion"`
}

type AccountDocument struct {
	AccountID   *int    `db:"account_id"`
	DocumentID  *int    `db:"document_id"`
	FileName    *string `db:"file_name"`
	FileContext *string `db:"file_context"`
	Tag         *string `db:"tag"`
}

type DocumentInfo struct {
	DocumentID   *int    `db:"document_id"`
	DocumentType *string `db:"document_type"`
}

type AccountDetail struct {
	AccountID            *int    `db:"account_id"`
	FirstName            *string `db:"first_name"`
	LastName             *string `db:"last_name"`
	Phone                *string `db:"phone"`
	Email                *string `db:"email"`
	Password             *string `db:"password"`
	AccountNumber        *string `db:"account_number"`
	IsVerify             *bool   `db:"is_verify"`
	Status               *string `db:"status"`
	CurrentAcceptVersion *string `db:"current_accept_version" json:"currentAcceptVersion"`
	DocumentID           *int    `db:"document_id"`
	DocumentType         *string `db:"document_type"`
	FileName             *string `db:"file_name"`
	FileContext          *string `db:"file_context"`
	Tag                  *string `db:"tag"`
}

type AccountRepository interface {
	SignUpAccountRepo(context.Context, string, string, string, string, string, string, string, string, string, string) (int64, error)
	GetAccountByEmailRepo(context.Context, string) (*Account, error)
	GetAccountByIDRepo(context.Context, int) (*Account, error)
	GetAccountRepo(context.Context, map[string]interface{}) (*[]AccountDetail, error)
	UpdateAccountRepo(context.Context, int, string) (int64, error)
	CreateWalletRepo(context.Context, int) error
	ConfirmVerifyEmailRepo(context.Context, int) (int64, error)
	GetTermsConditionRepo(context.Context, int) (*TermsCondition, error)
	AcceptTermsConditionRepo(context.Context, int, string) (int64, error)
	ConfirmChangePasswordRepo(context.Context, int, string) (int64, error)
}
