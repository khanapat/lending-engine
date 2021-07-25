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
