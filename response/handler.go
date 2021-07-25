package response

import (
	"context"
	"lending-engine/common"
)

var (
	EN = Global{
		AuthenBasicWeb:              ErrResponse{Code: ErrBasicAuthenticationCode, Title: ErrBasicAuthenticationMessageEN, Description: ErrAuthenticationDescEN},
		AuthorizationToken:          ErrResponse{Code: ErrUnauthorizationCode, Title: ErrAuthorizationTokenMessageEN, Description: ErrAuthorizationDescEN},
		SignUpAccountSuccess:        Response{Code: SuccessCode, Title: SuccessSignUpMessageEN},
		SignUpAccountRequest:        ErrResponse{Code: ErrInvalidRequestCode, Title: ErrSignUpMessageEN, Description: ErrRequestDataDescEN},
		ConfirmVerifyEmailSuccess:   Response{Code: SuccessCode, Title: SuccessConfirmVerifyEmailMessageEN},
		ConfirmVerifyEmailRequest:   ErrResponse{Code: ErrInvalidRequestCode, Title: ErrConfirmVerifyEmailMessageEN, Description: ErrRequestDataDescEN},
		LoginAccountSuccess:         Response{Code: SuccessCode, Title: SuccessLoginMessageEN},
		LoginAccountRequest:         ErrResponse{Code: ErrInvalidRequestCode, Title: ErrLoginMessageEN, Description: ErrRequestDataDescEN},
		AcceptTermsConditionSuccess: Response{Code: SuccessCode, Title: SuccessAcceptTermsConditionMessageEN},
		AcceptTermsConditionRequest: ErrResponse{Code: ErrInvalidRequestCode, Title: ErrAcceptTermsConditionMessageEN, Description: ErrRequestDataDescEN},
		GetTermsConditionSuccess:    Response{Code: SuccessCode, Title: SuccessGetTermsConditionMessageEN},
		GetTermsConditionRequest:    ErrResponse{Code: ErrInvalidRequestCode, Title: ErrAcceptTermsConditionMessageEN, Description: ErrRequestDataDescEN},
		GetTokenPriceSuccess:        Response{Code: SuccessCode, Title: SuccessGetToknPriceMessageEN},
		PreCalculationLoanSuccess:   Response{Code: SuccessCode, Title: SuccessPreCalculationLoanMessageEN},
		PreCalculationLoanRequest:   ErrResponse{Code: ErrInvalidRequestCode, Title: ErrPreCalculationLoanMessageEN, Description: ErrRequestDataDescEN},
		GetDepositSuccess:           Response{Code: SuccessCode, Title: SuccessGetDepositStatusMessageEN},
		SubmitDepositSuccess:        Response{Code: SuccessCode, Title: SuccessSubmitDepositMessageEN},
		SubmitDepositRequest:        ErrResponse{Code: ErrInvalidRequestCode, Title: ErrSubmitDepositMessageEN, Description: ErrRequestDataDescEN},
		SubmitDepositBlockErr:       ErrResponse{Code: ErrBlockchainCode, Title: ErrSubmitDepositMessageEN, Description: ErrContactAdminDescEN},
		GetCreditAvailableSuccess:   Response{Code: SuccessCode, Title: SuccessGetCreditAvailableMessageEN},
		BorrowLoanSuccess:           Response{Code: SuccessCode, Title: SuccessBorrowLoanMessageEN},
		BorrowLoanRequest:           ErrResponse{Code: ErrInvalidRequestCode, Title: ErrBorrowLoanMessageEN, Description: ErrRequestDataDescEN},
		GetInterestTermSuccess:      Response{Code: SuccessCode, Title: SuccessGetInterestTermMessageEN},
		GetDepositAdminSuccess:      Response{Code: SuccessCode, Title: SuccessGetDepositAdminMessageEN},
		GetDepositAdminRequest:      ErrResponse{Code: ErrInvalidRequestCode, Title: ErrGetDepositAdminMessageEN, Description: ErrRequestDataDescEN},
		ConfirmDepositAdminSuccess:  Response{Code: SuccessCode, Title: SuccessConfirmDepositAdminMessageEN},
		ConfirmDepositAdminRequest:  ErrResponse{Code: ErrInvalidRequestCode, Title: ErrConfirmDepositAdminMessageEN, Description: ErrRequestDataDescEN},
		GetContractAdminSuccess:     Response{Code: SuccessCode, Title: SuccessGetContractAdminMessageEN},
		GetContractAdminRequest:     ErrResponse{Code: ErrInvalidRequestCode, Title: ErrGetContractAdminMessageEN, Description: ErrRequestDataDescEN},
		ConfirmContractAdminSuccess: Response{Code: SuccessCode, Title: SuccessConfirmContractAdminMessageEN},
		ConfirmContractAdminRequest: ErrResponse{Code: ErrInvalidRequestCode, Title: ErrConfirmContractAdminMessageEN, Description: ErrRequestDataDescEN},
		InternalOperation:           ErrResponse{Code: ErrOperationCode, Title: ErrInternalServerMessageEN, Description: ErrContactAdminDescEN},
		InternalDatabase:            ErrResponse{Code: ErrDatabaseCode, Title: ErrInternalServerMessageEN, Description: ErrContactAdminDescEN},
		InternalRedis:               ErrResponse{Code: ErrRedisCode, Title: ErrInternalServerMessageEN, Description: ErrContactAdminDescEN},
	}
	TH = Global{
		AuthenBasicWeb:              ErrResponse{Code: ErrBasicAuthenticationCode, Title: ErrBasicAuthenticationMessageTH, Description: ErrAuthenticationDescTH},
		AuthorizationToken:          ErrResponse{Code: ErrUnauthorizationCode, Title: ErrAuthorizationTokenMessageTH, Description: ErrAuthorizationDescTH},
		SignUpAccountSuccess:        Response{Code: SuccessCode, Title: SuccessSignUpMessageTH},
		SignUpAccountRequest:        ErrResponse{Code: ErrInvalidRequestCode, Title: ErrSignUpMessageTH, Description: ErrRequestDataDescTH},
		ConfirmVerifyEmailSuccess:   Response{Code: SuccessCode, Title: SuccessConfirmVerifyEmailMessageTH},
		ConfirmVerifyEmailRequest:   ErrResponse{Code: ErrInvalidRequestCode, Title: ErrConfirmVerifyEmailMessageTH, Description: ErrRequestDataDescTH},
		LoginAccountSuccess:         Response{Code: SuccessCode, Title: SuccessLoginMessageTH},
		LoginAccountRequest:         ErrResponse{Code: ErrInvalidRequestCode, Title: ErrLoginMessageTH, Description: ErrRequestDataDescTH},
		AcceptTermsConditionSuccess: Response{Code: SuccessCode, Title: SuccessAcceptTermsConditionMessageTH},
		AcceptTermsConditionRequest: ErrResponse{Code: ErrInvalidRequestCode, Title: ErrAcceptTermsConditionMessageTH, Description: ErrRequestDataDescTH},
		GetTermsConditionSuccess:    Response{Code: SuccessCode, Title: SuccessGetTermsConditionMessageTH},
		GetTermsConditionRequest:    ErrResponse{Code: ErrInvalidRequestCode, Title: ErrAcceptTermsConditionMessageTH, Description: ErrRequestDataDescTH},
		GetTokenPriceSuccess:        Response{Code: SuccessCode, Title: SuccessGetToknPriceMessageTH},
		PreCalculationLoanSuccess:   Response{Code: SuccessCode, Title: SuccessPreCalculationLoanMessageTH},
		PreCalculationLoanRequest:   ErrResponse{Code: ErrInvalidRequestCode, Title: ErrPreCalculationLoanMessageTH, Description: ErrRequestDataDescTH},
		GetDepositSuccess:           Response{Code: SuccessCode, Title: SuccessGetDepositStatusMessageTH},
		SubmitDepositSuccess:        Response{Code: SuccessCode, Title: SuccessSubmitDepositMessageTH},
		SubmitDepositRequest:        ErrResponse{Code: ErrInvalidRequestCode, Title: ErrSubmitDepositMessageTH, Description: ErrRequestDataDescTH},
		SubmitDepositBlockErr:       ErrResponse{Code: ErrBlockchainCode, Title: ErrSubmitDepositMessageTH, Description: ErrContactAdminDescTH},
		GetCreditAvailableSuccess:   Response{Code: SuccessCode, Title: SuccessGetCreditAvailableMessageTH},
		BorrowLoanSuccess:           Response{Code: SuccessCode, Title: SuccessBorrowLoanMessageTH},
		BorrowLoanRequest:           ErrResponse{Code: ErrInvalidRequestCode, Title: ErrBorrowLoanMessageTH, Description: ErrRequestDataDescTH},
		GetInterestTermSuccess:      Response{Code: SuccessCode, Title: SuccessGetInterestTermMessageTH},
		GetDepositAdminSuccess:      Response{Code: SuccessCode, Title: SuccessGetDepositAdminMessageTH},
		GetDepositAdminRequest:      ErrResponse{Code: ErrInvalidRequestCode, Title: ErrGetDepositAdminMessageTH, Description: ErrRequestDataDescTH},
		ConfirmDepositAdminSuccess:  Response{Code: SuccessCode, Title: SuccessConfirmDepositAdminMessageTH},
		ConfirmDepositAdminRequest:  ErrResponse{Code: ErrInvalidRequestCode, Title: ErrConfirmDepositAdminMessageTH, Description: ErrRequestDataDescTH},
		GetContractAdminSuccess:     Response{Code: SuccessCode, Title: SuccessGetContractAdminMessageTH},
		GetContractAdminRequest:     ErrResponse{Code: ErrInvalidRequestCode, Title: ErrGetContractAdminMessageTH, Description: ErrRequestDataDescTH},
		ConfirmContractAdminSuccess: Response{Code: SuccessCode, Title: SuccessConfirmContractAdminMessageTH},
		ConfirmContractAdminRequest: ErrResponse{Code: ErrInvalidRequestCode, Title: ErrConfirmContractAdminMessageTH, Description: ErrRequestDataDescTH},
		InternalOperation:           ErrResponse{Code: ErrOperationCode, Title: ErrInternalServerMessageTH, Description: ErrContactAdminDescTH},
		InternalDatabase:            ErrResponse{Code: ErrDatabaseCode, Title: ErrInternalServerMessageTH, Description: ErrContactAdminDescTH},
		InternalRedis:               ErrResponse{Code: ErrRedisCode, Title: ErrInternalServerMessageTH, Description: ErrContactAdminDescTH},
	}

	Language = map[interface{}]Global{
		"en": EN,
		"th": TH,
	}
)

type Global struct {
	AuthenBasicWeb     ErrResponse
	AuthorizationToken ErrResponse
	// Account
	SignUpAccountSuccess        Response
	SignUpAccountRequest        ErrResponse
	ConfirmVerifyEmailSuccess   Response
	ConfirmVerifyEmailRequest   ErrResponse
	LoginAccountSuccess         Response
	LoginAccountRequest         ErrResponse
	AcceptTermsConditionSuccess Response
	AcceptTermsConditionRequest ErrResponse
	GetTermsConditionSuccess    Response
	GetTermsConditionRequest    ErrResponse
	// Lending
	//// User
	GetTokenPriceSuccess      Response
	PreCalculationLoanSuccess Response
	PreCalculationLoanRequest ErrResponse
	GetDepositSuccess         Response
	SubmitDepositSuccess      Response
	SubmitDepositRequest      ErrResponse
	SubmitDepositBlockErr     ErrResponse
	GetCreditAvailableSuccess Response
	BorrowLoanSuccess         Response
	BorrowLoanRequest         ErrResponse
	GetInterestTermSuccess    Response
	//// Admin
	GetDepositAdminSuccess      Response
	GetDepositAdminRequest      ErrResponse
	ConfirmDepositAdminSuccess  Response
	ConfirmDepositAdminRequest  ErrResponse
	GetContractAdminSuccess     Response
	GetContractAdminRequest     ErrResponse
	ConfirmContractAdminSuccess Response
	ConfirmContractAdminRequest ErrResponse
	// Basic
	InternalOperation ErrResponse
	InternalDatabase  ErrResponse
	InternalRedis     ErrResponse
}

func ResponseContextLocale(ctx context.Context) *Global {
	v := ctx.Value(common.LocaleKey)
	if v == nil {
		return nil
	}
	l, ok := Language[v]
	if ok {
		return &l
	}
	return &EN
}
