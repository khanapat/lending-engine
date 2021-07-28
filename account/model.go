package account

import (
	"fmt"
	"lending-engine/response"
	"unicode/utf8"

	"github.com/pkg/errors"
)

type SignUpRequest struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	AccountNumber  string `json:"accountNumber"`
	CitizenName    string `json:"citizenName"`
	CitizenCard    string `json:"citizenCard"`
	BookBankName   string `json:"bookBankName"`
	BookBankLedger string `json:"bookBankLedger"`
}

func (req *SignUpRequest) validate() error {
	if utf8.RuneCountInString(req.FirstName) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'firstName' must be REQUIRED field but the input is '%v'.", req.FirstName)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.LastName) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'lastName' must be REQUIRED field but the input is '%v'.", req.LastName)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.Phone) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'phone' must be REQUIRED field but the input is '%v'.", req.Phone)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.Email) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'email' must be REQUIRED field but the input is '%v'.", req.Email)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.Password) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'password' must be REQUIRED field but the input is '%v'.", req.Password)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.AccountNumber) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'accountNumber' must be REQUIRED field but the input is '%v'.", req.AccountNumber)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.CitizenName) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'citizenName' must be REQUIRED field but the input is '%v'.", req.CitizenName)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.CitizenCard) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'citizenCard' must be REQUIRED field but the input is '%v'.", req.CitizenCard)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.BookBankName) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'bookBankName' must be REQUIRED field but the input is '%v'.", req.BookBankName)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.BookBankLedger) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'bookBankLedger' must be REQUIRED field but the input is '%v'.", req.BookBankLedger)), response.ValidateFieldError)
	}
	return nil
}

type GetAccountAdminRequest struct {
	AccountID *int    `json:"accountId" example:"1"`
	Email     *string `json:"email" example:"k.apiwattanawong@gmail.com"`
}

type GetAccountAdminResponse struct {
	AccountID     int        `json:"accountId" example:"1"`
	FirstName     string     `json:"firstName" example:"Khanapat"`
	LastName      string     `json:"lastName" example:"Apiwattanawong"`
	Phone         string     `json:"phone" example:"0811111111"`
	Email         string     `json:"email" example:"k.apiwattanawong@gmail.com"`
	Password      string     `json:"password" example:"<bcrypt>"`
	AccountNumber string     `json:"accountNumber" example:"11111111"`
	IsVerify      bool       `json:"isVerify" example:"1"`
	Status        string     `json:"status" example:"PENDING"`
	TermCondition string     `json:"termCondition" example:"1.0.0"`
	Document      []Document `json:"document"`
}

type Document struct {
	DocumentID   int    `json:"documentId" example:"1"`
	DocumentType string `json:"documentType" example:"CITIZEN ID"`
	FileName     string `json:"fileName" example:"cid.pdf"`
	FileContext  string `json:"fileContext" example:"<Base64>"`
	Tag          string `json:"tag" example:"id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req *LoginRequest) validate() error {
	if utf8.RuneCountInString(req.Email) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'email' must be REQUIRED field but the input is '%v'.", req.Email)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.Password) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'password' must be REQUIRED field but the input is '%v'.", req.Password)), response.ValidateFieldError)
	}
	return nil
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDk2ODI2MzYsImlzcyI6ImFkbWluIn0.7MuvPeUSTlSvto45t3pqM5j7ZDxAVJnlBbEzq9ZJT0k"`
}

type SendVerifyEmailClientRequest struct {
	From     string                    `json:"from" example:"k.apiwattanawong@gmail.com"`
	To       []string                  `json:"to" example:"[yoisak4@gmail.com]"`
	Subject  string                    `json:"subject" example:"otp request"`
	Template string                    `json:"template" example:"otp.html"`
	Body     BodySendVerifyEmailClient `json:"body"`
	Auth     bool                      `json:"auth" example:"true"`
}

type BodySendVerifyEmailClient struct {
	Link string `json:"link" example:"www.lending.com/WERaOJOsfX"`
}

type SendVerifyEmailClientResult struct {
	Code        uint64 `json:"code" example:"2000"`
	Title       string `json:"title" example:"Success."`
	Description string `json:"description" example:"Please contact administrator for more information."`
}

type RequestResetPasswordRequest struct {
	Email string `json:"email" example:"k.apiwattanawong@gmail.com"`
}

func (req *RequestResetPasswordRequest) validate() error {
	if utf8.RuneCountInString(req.Email) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'email' must be REQUIRED field but the input is '%v'.", req.Email)), response.ValidateFieldError)
	}
	return nil
}

type ResetPasswordData struct {
	Otp       string `json:"otp" example:"999999"`
	FailCount int    `json:"failCount" example:"1"`
	AccountID int    `json:"accountId" example:"1"`
}

type RequestResetPasswordResponse struct {
	ReferenceNo string `json:"referenceNo" example:"XXXXXX"`
	ExpiredTime string `json:"expiredTime" example:"2021-01-02 12:13:14"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword" example:"BOBO"`
}

func (req *ResetPasswordRequest) validate() error {
	if utf8.RuneCountInString(req.NewPassword) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'newPassword' must be REQUIRED field but the input is '%v'.", req.NewPassword)), response.ValidateFieldError)
	}
	return nil
}
