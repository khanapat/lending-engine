package account

import (
	"fmt"
	"lending-engine/common"
	"lending-engine/response"
	"unicode/utf8"

	"github.com/pkg/errors"
)

// signup
type SignUpRequest struct {
	FirstName      string `json:"firstName" example:"Frank"`
	LastName       string `json:"lastName" example:"Style"`
	Phone          string `json:"phone" example:"0812345678"`
	Email          string `json:"email" example:"k.apiwattanawong@gmail.com"`
	Password       string `json:"password" example:"bobo"`
	AccountNumber  string `json:"accountNumber" example:"000000000"`
	CitizenName    string `json:"citizenName" example:"identity.jpg"`
	CitizenCard    string `json:"citizenCard" example:"<Base64>"`
	BookBankName   string `json:"bookBankName" example:"book.jpg"`
	BookBankLedger string `json:"bookBankLedger" example:"<Base64>"`
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
	if !common.EmailRegexp.MatchString(req.Email) {
		return errors.Wrapf(errors.New(fmt.Sprintf("'email' must be in format standard email but the input is '%v'.", req.Email)), response.ValidateFieldError)
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

type SignUpResponse struct {
	AccountID int64 `json:"accountId" example:"1"`
}

// get account admin
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

// confirm account admin
type ConfirmAccountAdminRequest struct {
	AccountID int `json:"accountId" example:"1"`
}

func (req *ConfirmAccountAdminRequest) validate() error {
	if req.AccountID == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'accountId' must be REQUIRED field but the input is '%v'.", req.AccountID)), response.ValidateFieldError)
	}
	return nil
}

// reject account admin
type RejectAccountAdminRequest struct {
	AccountID int `json:"accountId" example:"1"`
}

func (req *RejectAccountAdminRequest) validate() error {
	if req.AccountID == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'accountId' must be REQUIRED field but the input is '%v'.", req.AccountID)), response.ValidateFieldError)
	}
	return nil
}

// update account document admin
type UpdateAccountDocumentAdminRequest struct {
	AccountID   int    `json:"accountId" example:"1"`
	DocumentID  int    `json:"documentId" example:"1"`
	FileName    string `json:"fileName" example:"ID.jpg"`
	FileContext string `json:"fileContext" example:"<Base64>"`
}

func (req *UpdateAccountDocumentAdminRequest) validate() error {
	if req.AccountID == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'accountId' must be REQUIRED field but the input is '%v'.", req.AccountID)), response.ValidateFieldError)
	}
	if req.DocumentID == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'documentId' must be REQUIRED field but the input is '%v'.", req.DocumentID)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.FileName) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'fileName' must be REQUIRED field but the input is '%v'.", req.FileName)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.FileContext) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'fileContext' must be REQUIRED field but the input is '%v'.", req.FileContext)), response.ValidateFieldError)
	}
	return nil
}

// login
type LoginRequest struct {
	Email    string `json:"email" example:"k.apiwattanawong@gmail.com"`
	Password string `json:"password" example:"password"`
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

// term
type AcceptTermsConditionRequest struct {
	Version string `json:"version" example:"1.0.0"`
}

func (req *AcceptTermsConditionRequest) validate() error {
	if utf8.RuneCountInString(req.Version) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'version' must be REQUIRED field but the input is '%v'.", req.Version)), response.ValidateFieldError)
	}
	return nil
}

// request verify email
type SendVerifyEmailClientRequest struct {
	From     string                    `json:"from" example:"k.apiwattanawong@gmail.com"`
	To       []string                  `json:"to" example:"[yoisak4@gmail.com]"`
	Subject  string                    `json:"subject" example:"otp request"`
	Template string                    `json:"template" example:"otp.html"`
	Body     BodySendVerifyEmailClient `json:"body"`
	Auth     bool                      `json:"auth" example:"true"`
}

type BodySendVerifyEmailClient struct {
	Name string `json:"name" example:"trust momo"`
	Link string `json:"link" example:"www.lending.com/WERaOJOsfX"`
}

type SendVerifyEmailClientResult struct {
	Code        uint64 `json:"code" example:"2000"`
	Title       string `json:"title" example:"Success."`
	Description string `json:"description" example:"Please contact administrator for more information."`
}

// request reset password
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

// reset password
type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword" example:"BOBO"`
}

func (req *ResetPasswordRequest) validate() error {
	if utf8.RuneCountInString(req.NewPassword) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'newPassword' must be REQUIRED field but the input is '%v'.", req.NewPassword)), response.ValidateFieldError)
	}
	return nil
}

// create document info admin
type CreateDocumentInfoAdminRequest struct {
	DocumentType string `json:"documentType" example:"Citizen ID"`
}

func (req *CreateDocumentInfoAdminRequest) validate() error {
	if utf8.RuneCountInString(req.DocumentType) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'documentType' must be REQUIRED field but the input is '%v'.", req.DocumentType)), response.ValidateFieldError)
	}
	return nil
}

type CreateDocumentInfoAdminResponse struct {
	DocumentID int64 `json:"documentId" example:"1"`
}

// update document info admin
type UpdateDocumentInfoAdminRequest struct {
	DocumentID   int    `json:"documentId" example:"1"`
	DocumentType string `json:"documentType" example:"Citizen ID"`
}

func (req *UpdateDocumentInfoAdminRequest) validate() error {
	if req.DocumentID == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'documentId' must be REQUIRED field but the input is '%v'.", req.DocumentID)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.DocumentType) == 0 {
		return errors.Wrapf(errors.New(fmt.Sprintf("'documentType' must be REQUIRED field but the input is '%v'.", req.DocumentType)), response.ValidateFieldError)
	}
	return nil
}
