package lending

import (
	"fmt"
	"lending-engine/response"
	"unicode/utf8"

	"github.com/pkg/errors"
)

// deposit
type SubmitDepositRequest struct {
	Address        string  `json:"address" example:"0xc083EB69aa7215f4AFa7a22dcbfCC1a33999371C"`
	ChainID        int     `json:"chainId" example:"1"`
	TxnHash        string  `json:"txnHash" example:"0xf5a3aa87c40b05e6a308b61186eeded8996b654a9895401b8089a2966b54f618"`
	CollateralType string  `json:"collateralType" example:"BTC"`
	Volume         float64 `json:"volume" example:"0.5"`
}

func (req *SubmitDepositRequest) validate() error {
	if utf8.RuneCountInString(req.Address) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'address' must be REQUIRED field but the input is '%v'.", req.Address)), response.ValidateFieldError)
	}
	if req.ChainID == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'chainId' must be REQUIRED field but the input is '%v'.", req.ChainID)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.TxnHash) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'txnHash' must be REQUIRED field but the input is '%v'.", req.TxnHash)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.CollateralType) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'collateralType' must be REQUIRED field but the input is '%v'.", req.CollateralType)), response.ValidateFieldError)
	}
	if req.Volume == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'volume' must be REQUIRED field but the input is '%v'.", req.Volume)), response.ValidateFieldError)
	}
	return nil
}

type SubmitDepositResponse struct {
	DepositID int64 `json:"depositId" example:"1"`
}

type GetDepositAdminRequest struct {
	ID        *int    `json:"id" example:"1"`
	AccountID *int    `json:"accountId" example:"4"`
	Address   *string `json:"address" example:"0xc083EB69aa7215f4AFa7a22dcbfCC1a33999371C"`
}

// credit
type GetCreditAvailableResponse struct {
	BTCVolume       float64 `json:"btcVolume" example:"0.1"`
	ETHVolume       float64 `json:"ethVolume" example:"0.1"`
	CollateralValue float64 `json:"collateralValue" example:"10000"`
	LoanOutstanding float64 `json:"loanOutstanding" example:"0"`
	CreditAvailable float64 `json:"creditAvailable" example:"10000"`
}

// Borrow
type BorrowLoanRequest struct {
	Loan         float64 `json:"loan" example:"1000"`
	InterestCode int     `json:"interestCode" example:"1"`
	Term         int     `json:"term" example:"12"`
}

func (req *BorrowLoanRequest) validate() error {
	if req.Loan == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'loan' must be REQUIRED field but the input is '%v'.", req.Loan)), response.ValidateFieldError)
	}
	if req.InterestCode == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'interestCode' must be REQUIRED field but the input is '%v'.", req.InterestCode)), response.ValidateFieldError)
	}
	if req.Term == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'term' must be REQUIRED field but the input is '%v'.", req.Term)), response.ValidateFieldError)
	}
	return nil
}

type BorrowLoanResponse struct {
	ContractID int64 `json:"contractId" example:"1"`
}

type GetLoanAdminRequest struct {
	ContractID *int `json:"contractId" example:"1"`
	AccountID  *int `json:"accountId" example:"1"`
}

// Repay
type SubmitRepayRequest struct {
	ContractID int    `json:"contractId" example:"1"`
	Amount     int    `json:"amount" example:"1000"`
	Slip       string `json:"slip" example:"<Base64>"`
}

func (req *SubmitRepayRequest) validate() error {
	if req.ContractID == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'contractId' must be REQUIRED field but the input is '%v'.", req.ContractID)), response.ValidateFieldError)
	}
	if req.Amount == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'amount' must be REQUIRED field but the input is '%v'.", req.Amount)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.Slip) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'slip' must be REQUIRED field but the input is '%v'.", req.Slip)), response.ValidateFieldError)
	}
	return nil
}

type SubmitRepayResponse struct {
	RepayID int64 `json:"repayId" example:"1"`
}

type GetRepayAdminRequest struct {
	ID         *int `json:"id" example:"1"`
	ContractID *int `json:"contractId" example:"1"`
	AccountID  *int `json:"accountId" example:"1"`
}

// Price
type GetTokenPriceResponse struct {
	BTC          TokenPrice `json:"btc"`
	ETH          TokenPrice `json:"eth"`
	InterestRate float64    `json:"interestRate" example:"0.05"`
}

type TokenPrice struct {
	Price   float64 `json:"price" example:"1042475.25"`
	Haircut float64 `json:"haircut" example:"0.5"`
}

type PreCalculationLoanRequest struct {
	BTCAmount float64 `json:"btcAmount" example:"0.5"`
	ETHAmount float64 `json:"ethAmount" example:"0.5"`
	Period    int     `json:"period" example:"12"`
}

func (req *PreCalculationLoanRequest) validate() error {
	if req.Period == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'period' must be REQUIRED field but the input is '%v'.", req.Period)), response.ValidateFieldError)
	}
	return nil
}

type PreCalculationLoanResponse struct {
	BTC     TokenPriceRate `json:"btc"`
	ETH     TokenPriceRate `json:"eth"`
	Summary SummaryLoan    `json:"summary"`
}

type TokenPriceRate struct {
	Volume     float64 `json:"volume" example:"0"`
	Haircut    float64 `json:"haircut" example:"0.5"`
	LoanAmount float64 `json:"loanAmount" example:"200000"`
}

type SummaryLoan struct {
	TotalLoanAmount float64 `json:"totalLoanAmount" example:"2000000"`
	InterestRate    float64 `json:"interestRate" example:"0.05"`
	MonthlyInterest float64 `json:"monthlyInterest" example:"1666.67"`
	Period          int     `json:"period" example:"12"`
	TotalInterest   float64 `json:"totalInterest" example:"5000"`
}
