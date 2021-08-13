package response

import (
	"context"
	"lending-engine/common"
)

var (
	EN = Global{
		AuthenBasicWeb:                    ErrResponse{Code: ErrBasicAuthenticationCode, Title: ErrBasicAuthenticationMessageEN, Description: ErrAuthenticationDescEN},
		AuthorizationToken:                ErrResponse{Code: ErrUnauthorizationCode, Title: ErrAuthorizationTokenMessageEN, Description: ErrAuthorizationDescEN},
		SignUpAccountSuccess:              Response{Code: SuccessCode, Title: SuccessSignUpMessageEN},
		SignUpAccountRequest:              ErrResponse{Code: ErrInvalidRequestCode, Title: ErrSignUpMessageEN, Description: ErrRequestDataDescEN},
		SignUpAccountThirdParty:           ErrResponse{Code: ErrThirdPartyCode, Title: ErrSignUpMessageEN, Description: ErrThirdPartyDescEN},
		ConfirmVerifyEmailSuccess:         Response{Code: SuccessCode, Title: SuccessConfirmVerifyEmailMessageEN},
		ConfirmVerifyEmailRequest:         ErrResponse{Code: ErrInvalidRequestCode, Title: ErrConfirmVerifyEmailMessageEN, Description: ErrRequestDataDescEN},
		LoginAccountSuccess:               Response{Code: SuccessCode, Title: SuccessLoginMessageEN},
		LoginAccountRequest:               ErrResponse{Code: ErrInvalidRequestCode, Title: ErrLoginMessageEN, Description: ErrRequestDataDescEN},
		AcceptTermsConditionSuccess:       Response{Code: SuccessCode, Title: SuccessAcceptTermsConditionMessageEN},
		AcceptTermsConditionRequest:       ErrResponse{Code: ErrInvalidRequestCode, Title: ErrAcceptTermsConditionMessageEN, Description: ErrRequestDataDescEN},
		GetTermsConditionSuccess:          Response{Code: SuccessCode, Title: SuccessGetTermsConditionMessageEN},
		GetTermsConditionRequest:          ErrResponse{Code: ErrInvalidRequestCode, Title: ErrAcceptTermsConditionMessageEN, Description: ErrRequestDataDescEN},
		RequestResetPasswordSuccess:       Response{Code: SuccessCode, Title: SuccessRequestResetPasswordMessageEN},
		RequestResetPasswordRequest:       ErrResponse{Code: ErrInvalidRequestCode, Title: ErrRequestResetPasswordMessageEN, Description: ErrRequestDataDescEN},
		RequestResetPasswordThirdParty:    ErrResponse{Code: ErrThirdPartyCode, Title: ErrRequestResetPasswordMessageEN, Description: ErrThirdPartyDescEN},
		ResetPasswordSuccess:              Response{Code: SuccessCode, Title: SuccessResetPasswordMessageEN},
		ResetPasswordRequest:              ErrResponse{Code: ErrInvalidRequestCode, Title: ErrResetPasswordMessageEN, Description: ErrRequestDataDescEN},
		GetDocumentInfoAdminSuccess:       Response{Code: SuccessCode, Title: SuccessGetDocumentInfoAdminMessageEN},
		CreateDocumentInfoAdminSuccess:    Response{Code: SuccessCode, Title: SuccessCreateDocumentInfoAdminMessageEN},
		CreateDocumentInfoAdminRequest:    ErrResponse{Code: ErrInvalidRequestCode, Title: ErrCreateDocumentInfoAdminMessageEN, Description: ErrRequestDataDescEN},
		UpdateDocumentInfoAdminSuccess:    Response{Code: SuccessCode, Title: SuccessUpdateDocumentInfoAdminMessageEN},
		UpdateDocumentInfoAdminRequest:    ErrResponse{Code: ErrInvalidRequestCode, Title: ErrUpdateDocumentInfoAdminMessageEN, Description: ErrRequestDataDescEN},
		GetTokenPriceSuccess:              Response{Code: SuccessCode, Title: SuccessGetToknPriceMessageEN},
		PreCalculationLoanSuccess:         Response{Code: SuccessCode, Title: SuccessPreCalculationLoanMessageEN},
		PreCalculationLoanRequest:         ErrResponse{Code: ErrInvalidRequestCode, Title: ErrPreCalculationLoanMessageEN, Description: ErrRequestDataDescEN},
		GetWalletTransactionSuccess:       Response{Code: SuccessCode, Title: SuccessGetWalletTransactionMessageEN},
		SubmitDepositSuccess:              Response{Code: SuccessCode, Title: SuccessSubmitDepositMessageEN},
		SubmitDepositRequest:              ErrResponse{Code: ErrInvalidRequestCode, Title: ErrSubmitDepositMessageEN, Description: ErrRequestDataDescEN},
		SubmitDepositBlockErr:             ErrResponse{Code: ErrBlockchainCode, Title: ErrSubmitDepositMessageEN, Description: ErrContactAdminDescEN},
		SubmitWithdrawSuccess:             Response{Code: SuccessCode, Title: SuccessSubmitWithdrawMessageEN},
		SubmitWithdrawRequest:             ErrResponse{Code: ErrInvalidRequestCode, Title: ErrSubmitDepositMessageEN, Description: ErrRequestDataDescEN},
		GetCreditAvailableSuccess:         Response{Code: SuccessCode, Title: SuccessGetCreditAvailableMessageEN},
		GetLoanSuccess:                    Response{Code: SuccessCode, Title: SuccessGetLoanMessageEN},
		BorrowLoanSuccess:                 Response{Code: SuccessCode, Title: SuccessBorrowLoanMessageEN},
		BorrowLoanRequest:                 ErrResponse{Code: ErrInvalidRequestCode, Title: ErrBorrowLoanMessageEN, Description: ErrRequestDataDescEN},
		GetInterestTermSuccess:            Response{Code: SuccessCode, Title: SuccessGetInterestTermMessageEN},
		GetRepaymentSuccess:               Response{Code: SuccessCode, Title: SuccessGetRepaymentMessageEN},
		GetRepaymentRequest:               ErrResponse{Code: ErrInvalidRequestCode, Title: ErrGetRepaymentMessageEN, Description: ErrRequestDataDescEN},
		SubmitRepaymentSuccess:            Response{Code: SuccessCode, Title: SuccessSubmitRepaymentMessageEN},
		SubmitRepaymentRequest:            ErrResponse{Code: ErrInvalidRequestCode, Title: ErrSubmitRepaymentMessageEN, Description: ErrRequestDataDescEN},
		GetAccountAdminSuccess:            Response{Code: SuccessCode, Title: SuccessGetAccountAdminMessageEN},
		GetAccountAdminRequest:            ErrResponse{Code: ErrInvalidRequestCode, Title: ErrGetAccountAdminMessageEN, Description: ErrRequestDataDescEN},
		ConfirmAccountAdminSuccess:        Response{Code: SuccessCode, Title: SuccessConfirmAccountAdminMessageEN},
		ConfirmAccountAdminRequest:        ErrResponse{Code: ErrInvalidRequestCode, Title: ErrConfirmAccountAdminMessageEN, Description: ErrRequestDataDescEN},
		RejectAccountAdminSuccess:         Response{Code: SuccessCode, Title: SuccessRejectAccountAdminMessageEN},
		RejectAccountAdminRequest:         ErrResponse{Code: ErrInvalidRequestCode, Title: ErrRejectAccountAdminMessageEN, Description: ErrRequestDataDescEN},
		UpdateAccountDocumentAdminSuccess: Response{Code: SuccessCode, Title: SuccessUpdateAccountDocumentAdminMessageEN},
		UpdateAccountDocumentAdminRequest: ErrResponse{Code: ErrInvalidRequestCode, Title: ErrUpdateAccountDocumentAdminMessageEN, Description: ErrRequestDataDescEN},
		GetWalletTransactionAdminSuccess:  Response{Code: SuccessCode, Title: SuccessGetWalletTransactionAdminMessageEN},
		GetWalletTransactionAdminRequest:  ErrResponse{Code: ErrInvalidRequestCode, Title: ErrGetWalletTransactionAdminMessageEN, Description: ErrRequestDataDescEN},
		ConfirmDepositAdminSuccess:        Response{Code: SuccessCode, Title: SuccessConfirmDepositAdminMessageEN},
		ConfirmDepositAdminRequest:        ErrResponse{Code: ErrInvalidRequestCode, Title: ErrConfirmDepositAdminMessageEN, Description: ErrRequestDataDescEN},
		RejectDepositAdminSuccess:         Response{Code: SuccessCode, Title: SuccessRejectDepositAdminMessageEN},
		RejectDepostiAdminRequest:         ErrResponse{Code: ErrInvalidRequestCode, Title: ErrRejectDepositAdminMessageEN, Description: ErrRequestDataDescEN},
		ConfirmWithdrawAdminSuccess:       Response{Code: SuccessCode, Title: SuccessConfirmWithdrawAdminMessageEN},
		ConfirmWithdrawAdminRequest:       ErrResponse{Code: ErrInvalidRequestCode, Title: ErrConfirmWithdrawAdminMessageEN, Description: ErrRequestDataDescEN},
		RejectWithdrawAdminSuccess:        Response{Code: SuccessCode, Title: SuccessRejectWithdrawAdminMessageEN},
		RejectWithdrawAdminRequest:        ErrResponse{Code: ErrInvalidRequestCode, Title: ErrRejectWithdrawAdminMessageEN, Description: ErrRequestDataDescEN},
		GetContractAdminSuccess:           Response{Code: SuccessCode, Title: SuccessGetContractAdminMessageEN},
		GetContractAdminRequest:           ErrResponse{Code: ErrInvalidRequestCode, Title: ErrGetContractAdminMessageEN, Description: ErrRequestDataDescEN},
		ConfirmContractAdminSuccess:       Response{Code: SuccessCode, Title: SuccessConfirmContractAdminMessageEN},
		ConfirmContractAdminRequest:       ErrResponse{Code: ErrInvalidRequestCode, Title: ErrConfirmContractAdminMessageEN, Description: ErrRequestDataDescEN},
		CreateInterestTermAdminSuccess:    Response{Code: SuccessCode, Title: SuccessCreateInterestTermAdminMessageEN},
		CreateInterestTermAdminRequest:    ErrResponse{Code: ErrInvalidRequestCode, Title: ErrCreateInterestTermAdminMessageEN, Description: ErrRequestDataDescEN},
		UpdateInterestTermAdminSuccess:    Response{Code: SuccessCode, Title: SuccessUpdateInterestTermAdminMessageEN},
		UpdateInterestTermAdminRequest:    ErrResponse{Code: ErrInvalidRequestCode, Title: ErrUpdateInterestTermAdminMessageEN, Description: ErrRequestDataDescEN},
		GetRepaymentAdminSuccess:          Response{Code: SuccessCode, Title: SuccessGetRepaymentAdminMessageEN},
		GetRepaymentAdminRequest:          ErrResponse{Code: ErrInvalidRequestCode, Title: ErrGetRepaymentAdminMessageEN, Description: ErrRequestDataDescEN},
		ConfirmRepaymentAdminSuccess:      Response{Code: SuccessCode, Title: SuccessConfirmRepaymentAdminMessageEN},
		ConfirmRepaymentAdminRequest:      ErrResponse{Code: ErrInvalidRequestCode, Title: ErrConfirmRepaymentAdminMessageEN, Description: ErrRequestDataDescEN},
		RejectRepaymentAdminSuccess:       Response{Code: SuccessCode, Title: SuccessRejectRepaymentAdminMessageEN},
		RejectRepaymentAdminRequest:       ErrResponse{Code: ErrInvalidRequestCode, Title: ErrRejectRepaymentAdminMessageEN, Description: ErrRequestDataDescEN},
		LiquidateFundAdminSuccess:         Response{Code: SuccessCode, Title: SuccessLiquidateFundAdminMessageEN},
		LiquidateFundAdminRequest:         ErrResponse{Code: ErrInvalidRequestCode, Title: ErrLiquidateFundAdminMessageEN, Description: ErrRequestDataDescEN},
		GetOTPSuccess:                     Response{Code: SuccessCode, Title: SuccessOTPRequestMessageEN},
		GetOTPRequest:                     ErrResponse{Code: ErrInvalidRequestCode, Title: ErrOTPRequestMessageEN, Description: ErrRequestDataDescEN},
		GetOTPThirdParty:                  ErrResponse{Code: ErrThirdPartyCode, Title: ErrOTPRequestMessageEN, Description: ErrThirdPartyDescEN},
		OTPRequestInvalid:                 ErrResponse{Code: ErrOTPRequestCode, Title: ErrInvalidOTPMessageEN, Description: ErrRequestDataDescEN},
		OTPRequestFailLimit:               ErrResponse{Code: ErrOTPRequestCode, Title: ErrLimitInvalidOTPMessageEN, Description: ErrCooldownDescEN},
		OTPRequestMaxLimit:                ErrResponse{Code: ErrOTPRequestCode, Title: ErrLimitOTPRequestMessageEN, Description: ErrCooldownDescEN},
		OTPRequestDuplicate:               ErrResponse{Code: ErrOTPRequestCode, Title: ErrDuplicateOTPRequestMessageEN, Description: ErrCooldownDescEN},
		InternalOperation:                 ErrResponse{Code: ErrOperationCode, Title: ErrInternalServerMessageEN, Description: ErrContactAdminDescEN},
		InternalDatabase:                  ErrResponse{Code: ErrDatabaseCode, Title: ErrInternalServerMessageEN, Description: ErrContactAdminDescEN},
		InternalRedis:                     ErrResponse{Code: ErrRedisCode, Title: ErrInternalServerMessageEN, Description: ErrContactAdminDescEN},
	}
	TH = Global{
		AuthenBasicWeb:                    ErrResponse{Code: ErrBasicAuthenticationCode, Title: ErrBasicAuthenticationMessageTH, Description: ErrAuthenticationDescTH},
		AuthorizationToken:                ErrResponse{Code: ErrUnauthorizationCode, Title: ErrAuthorizationTokenMessageTH, Description: ErrAuthorizationDescTH},
		SignUpAccountSuccess:              Response{Code: SuccessCode, Title: SuccessSignUpMessageTH},
		SignUpAccountRequest:              ErrResponse{Code: ErrInvalidRequestCode, Title: ErrSignUpMessageTH, Description: ErrRequestDataDescTH},
		SignUpAccountThirdParty:           ErrResponse{Code: ErrThirdPartyCode, Title: ErrSignUpMessageTH, Description: ErrThirdPartyDescTH},
		ConfirmVerifyEmailSuccess:         Response{Code: SuccessCode, Title: SuccessConfirmVerifyEmailMessageTH},
		ConfirmVerifyEmailRequest:         ErrResponse{Code: ErrInvalidRequestCode, Title: ErrConfirmVerifyEmailMessageTH, Description: ErrRequestDataDescTH},
		LoginAccountSuccess:               Response{Code: SuccessCode, Title: SuccessLoginMessageTH},
		LoginAccountRequest:               ErrResponse{Code: ErrInvalidRequestCode, Title: ErrLoginMessageTH, Description: ErrRequestDataDescTH},
		AcceptTermsConditionSuccess:       Response{Code: SuccessCode, Title: SuccessAcceptTermsConditionMessageTH},
		AcceptTermsConditionRequest:       ErrResponse{Code: ErrInvalidRequestCode, Title: ErrAcceptTermsConditionMessageTH, Description: ErrRequestDataDescTH},
		GetTermsConditionSuccess:          Response{Code: SuccessCode, Title: SuccessGetTermsConditionMessageTH},
		GetTermsConditionRequest:          ErrResponse{Code: ErrInvalidRequestCode, Title: ErrAcceptTermsConditionMessageTH, Description: ErrRequestDataDescTH},
		RequestResetPasswordSuccess:       Response{Code: SuccessCode, Title: SuccessRequestResetPasswordMessageTH},
		RequestResetPasswordRequest:       ErrResponse{Code: ErrInvalidRequestCode, Title: ErrRequestResetPasswordMessageTH, Description: ErrRequestDataDescTH},
		RequestResetPasswordThirdParty:    ErrResponse{Code: ErrThirdPartyCode, Title: ErrRequestResetPasswordMessageTH, Description: ErrThirdPartyDescTH},
		ResetPasswordSuccess:              Response{Code: SuccessCode, Title: SuccessResetPasswordMessageTH},
		ResetPasswordRequest:              ErrResponse{Code: ErrInvalidRequestCode, Title: ErrResetPasswordMessageTH, Description: ErrRequestDataDescTH},
		GetDocumentInfoAdminSuccess:       Response{Code: SuccessCode, Title: SuccessGetDocumentInfoAdminMessageTH},
		CreateDocumentInfoAdminSuccess:    Response{Code: SuccessCode, Title: SuccessCreateDocumentInfoAdminMessageTH},
		CreateDocumentInfoAdminRequest:    ErrResponse{Code: ErrInvalidRequestCode, Title: ErrCreateDocumentInfoAdminMessageTH, Description: ErrRequestDataDescTH},
		UpdateDocumentInfoAdminSuccess:    Response{Code: SuccessCode, Title: SuccessUpdateDocumentInfoAdminMessageTH},
		UpdateDocumentInfoAdminRequest:    ErrResponse{Code: ErrInvalidRequestCode, Title: ErrUpdateDocumentInfoAdminMessageTH, Description: ErrRequestDataDescTH},
		GetTokenPriceSuccess:              Response{Code: SuccessCode, Title: SuccessGetToknPriceMessageTH},
		PreCalculationLoanSuccess:         Response{Code: SuccessCode, Title: SuccessPreCalculationLoanMessageTH},
		PreCalculationLoanRequest:         ErrResponse{Code: ErrInvalidRequestCode, Title: ErrPreCalculationLoanMessageTH, Description: ErrRequestDataDescTH},
		GetWalletTransactionSuccess:       Response{Code: SuccessCode, Title: SuccessGetWalletTransactionMessageTH},
		SubmitDepositSuccess:              Response{Code: SuccessCode, Title: SuccessSubmitDepositMessageTH},
		SubmitDepositRequest:              ErrResponse{Code: ErrInvalidRequestCode, Title: ErrSubmitDepositMessageTH, Description: ErrRequestDataDescTH},
		SubmitDepositBlockErr:             ErrResponse{Code: ErrBlockchainCode, Title: ErrSubmitDepositMessageTH, Description: ErrContactAdminDescTH},
		SubmitWithdrawSuccess:             Response{Code: SuccessCode, Title: SuccessSubmitWithdrawMessageTH},
		SubmitWithdrawRequest:             ErrResponse{Code: ErrInvalidRequestCode, Title: ErrSubmitDepositMessageTH, Description: ErrRequestDataDescTH},
		GetCreditAvailableSuccess:         Response{Code: SuccessCode, Title: SuccessGetCreditAvailableMessageTH},
		GetLoanSuccess:                    Response{Code: SuccessCode, Title: SuccessGetLoanMessageTH},
		BorrowLoanSuccess:                 Response{Code: SuccessCode, Title: SuccessBorrowLoanMessageTH},
		BorrowLoanRequest:                 ErrResponse{Code: ErrInvalidRequestCode, Title: ErrBorrowLoanMessageTH, Description: ErrRequestDataDescTH},
		GetInterestTermSuccess:            Response{Code: SuccessCode, Title: SuccessGetInterestTermMessageTH},
		GetRepaymentSuccess:               Response{Code: SuccessCode, Title: SuccessGetRepaymentMessageTH},
		GetRepaymentRequest:               ErrResponse{Code: ErrInvalidRequestCode, Title: ErrGetRepaymentMessageTH, Description: ErrRequestDataDescTH},
		SubmitRepaymentSuccess:            Response{Code: SuccessCode, Title: SuccessSubmitRepaymentMessageTH},
		SubmitRepaymentRequest:            ErrResponse{Code: ErrInvalidRequestCode, Title: ErrSubmitRepaymentMessageTH, Description: ErrRequestDataDescTH},
		GetAccountAdminSuccess:            Response{Code: SuccessCode, Title: SuccessGetAccountAdminMessageTH},
		GetAccountAdminRequest:            ErrResponse{Code: ErrInvalidRequestCode, Title: ErrGetAccountAdminMessageTH, Description: ErrRequestDataDescTH},
		ConfirmAccountAdminSuccess:        Response{Code: SuccessCode, Title: SuccessConfirmAccountAdminMessageTH},
		ConfirmAccountAdminRequest:        ErrResponse{Code: ErrInvalidRequestCode, Title: ErrConfirmAccountAdminMessageTH, Description: ErrRequestDataDescTH},
		RejectAccountAdminSuccess:         Response{Code: SuccessCode, Title: SuccessRejectAccountAdminMessageTH},
		RejectAccountAdminRequest:         ErrResponse{Code: ErrInvalidRequestCode, Title: ErrRejectAccountAdminMessageTH, Description: ErrRequestDataDescTH},
		UpdateAccountDocumentAdminSuccess: Response{Code: SuccessCode, Title: SuccessUpdateAccountDocumentAdminMessageTH},
		UpdateAccountDocumentAdminRequest: ErrResponse{Code: ErrInvalidRequestCode, Title: ErrUpdateAccountDocumentAdminMessageTH, Description: ErrRequestDataDescTH},
		GetWalletTransactionAdminSuccess:  Response{Code: SuccessCode, Title: SuccessGetWalletTransactionAdminMessageTH},
		GetWalletTransactionAdminRequest:  ErrResponse{Code: ErrInvalidRequestCode, Title: ErrGetWalletTransactionAdminMessageTH, Description: ErrRequestDataDescTH},
		ConfirmDepositAdminSuccess:        Response{Code: SuccessCode, Title: SuccessConfirmDepositAdminMessageTH},
		ConfirmDepositAdminRequest:        ErrResponse{Code: ErrInvalidRequestCode, Title: ErrConfirmDepositAdminMessageTH, Description: ErrRequestDataDescTH},
		RejectDepositAdminSuccess:         Response{Code: SuccessCode, Title: SuccessRejectDepositAdminMessageTH},
		RejectDepostiAdminRequest:         ErrResponse{Code: ErrInvalidRequestCode, Title: ErrRejectDepositAdminMessageTH, Description: ErrRequestDataDescTH},
		ConfirmWithdrawAdminSuccess:       Response{Code: SuccessCode, Title: SuccessConfirmWithdrawAdminMessageTH},
		ConfirmWithdrawAdminRequest:       ErrResponse{Code: ErrInvalidRequestCode, Title: ErrConfirmWithdrawAdminMessageTH, Description: ErrRequestDataDescTH},
		RejectWithdrawAdminSuccess:        Response{Code: SuccessCode, Title: SuccessRejectWithdrawAdminMessageTH},
		RejectWithdrawAdminRequest:        ErrResponse{Code: ErrInvalidRequestCode, Title: ErrRejectWithdrawAdminMessageTH, Description: ErrRequestDataDescTH},
		GetContractAdminSuccess:           Response{Code: SuccessCode, Title: SuccessGetContractAdminMessageTH},
		GetContractAdminRequest:           ErrResponse{Code: ErrInvalidRequestCode, Title: ErrGetContractAdminMessageTH, Description: ErrRequestDataDescTH},
		ConfirmContractAdminSuccess:       Response{Code: SuccessCode, Title: SuccessConfirmContractAdminMessageTH},
		ConfirmContractAdminRequest:       ErrResponse{Code: ErrInvalidRequestCode, Title: ErrConfirmContractAdminMessageTH, Description: ErrRequestDataDescTH},
		CreateInterestTermAdminSuccess:    Response{Code: SuccessCode, Title: SuccessCreateInterestTermAdminMessageTH},
		CreateInterestTermAdminRequest:    ErrResponse{Code: ErrInvalidRequestCode, Title: ErrCreateInterestTermAdminMessageTH, Description: ErrRequestDataDescTH},
		UpdateInterestTermAdminSuccess:    Response{Code: SuccessCode, Title: SuccessUpdateInterestTermAdminMessageTH},
		UpdateInterestTermAdminRequest:    ErrResponse{Code: ErrInvalidRequestCode, Title: ErrUpdateInterestTermAdminMessageTH, Description: ErrRequestDataDescTH},
		GetRepaymentAdminSuccess:          Response{Code: SuccessCode, Title: SuccessGetRepaymentAdminMessageTH},
		GetRepaymentAdminRequest:          ErrResponse{Code: ErrInvalidRequestCode, Title: ErrGetRepaymentAdminMessageTH, Description: ErrRequestDataDescTH},
		ConfirmRepaymentAdminSuccess:      Response{Code: SuccessCode, Title: SuccessConfirmRepaymentAdminMessageTH},
		ConfirmRepaymentAdminRequest:      ErrResponse{Code: ErrInvalidRequestCode, Title: ErrConfirmRepaymentAdminMessageTH, Description: ErrRequestDataDescTH},
		RejectRepaymentAdminSuccess:       Response{Code: SuccessCode, Title: SuccessRejectRepaymentAdminMessageTH},
		RejectRepaymentAdminRequest:       ErrResponse{Code: ErrInvalidRequestCode, Title: ErrRejectRepaymentAdminMessageTH, Description: ErrRequestDataDescTH},
		LiquidateFundAdminSuccess:         Response{Code: SuccessCode, Title: SuccessLiquidateFundAdminMessageTH},
		LiquidateFundAdminRequest:         ErrResponse{Code: ErrInvalidRequestCode, Title: ErrLiquidateFundAdminMessageTH, Description: ErrRequestDataDescTH},
		GetOTPSuccess:                     Response{Code: SuccessCode, Title: SuccessOTPRequestMessageTH},
		GetOTPRequest:                     ErrResponse{Code: ErrInvalidRequestCode, Title: ErrOTPRequestMessageTH, Description: ErrRequestDataDescTH},
		GetOTPThirdParty:                  ErrResponse{Code: ErrThirdPartyCode, Title: ErrOTPRequestMessageTH, Description: ErrThirdPartyDescTH},
		OTPRequestInvalid:                 ErrResponse{Code: ErrOTPRequestCode, Title: ErrInvalidOTPMessageTH, Description: ErrRequestDataDescTH},
		OTPRequestFailLimit:               ErrResponse{Code: ErrOTPRequestCode, Title: ErrLimitInvalidOTPMessageTH, Description: ErrCooldownDescTH},
		OTPRequestMaxLimit:                ErrResponse{Code: ErrOTPRequestCode, Title: ErrLimitOTPRequestMessageTH, Description: ErrCooldownDescTH},
		OTPRequestDuplicate:               ErrResponse{Code: ErrOTPRequestCode, Title: ErrDuplicateOTPRequestMessageTH, Description: ErrCooldownDescTH},
		InternalOperation:                 ErrResponse{Code: ErrOperationCode, Title: ErrInternalServerMessageTH, Description: ErrContactAdminDescTH},
		InternalDatabase:                  ErrResponse{Code: ErrDatabaseCode, Title: ErrInternalServerMessageTH, Description: ErrContactAdminDescTH},
		InternalRedis:                     ErrResponse{Code: ErrRedisCode, Title: ErrInternalServerMessageTH, Description: ErrContactAdminDescTH},
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
	//// User
	SignUpAccountSuccess           Response
	SignUpAccountRequest           ErrResponse
	SignUpAccountThirdParty        ErrResponse
	ConfirmVerifyEmailSuccess      Response
	ConfirmVerifyEmailRequest      ErrResponse
	LoginAccountSuccess            Response
	LoginAccountRequest            ErrResponse
	AcceptTermsConditionSuccess    Response
	AcceptTermsConditionRequest    ErrResponse
	GetTermsConditionSuccess       Response
	GetTermsConditionRequest       ErrResponse
	RequestResetPasswordSuccess    Response
	RequestResetPasswordRequest    ErrResponse
	RequestResetPasswordThirdParty ErrResponse
	ResetPasswordSuccess           Response
	ResetPasswordRequest           ErrResponse
	//// Admin
	GetDocumentInfoAdminSuccess    Response
	CreateDocumentInfoAdminSuccess Response
	CreateDocumentInfoAdminRequest ErrResponse
	UpdateDocumentInfoAdminSuccess Response
	UpdateDocumentInfoAdminRequest ErrResponse
	// Lending
	//// User
	GetTokenPriceSuccess        Response
	PreCalculationLoanSuccess   Response
	PreCalculationLoanRequest   ErrResponse
	GetWalletTransactionSuccess Response
	SubmitDepositSuccess        Response
	SubmitDepositRequest        ErrResponse
	SubmitDepositBlockErr       ErrResponse
	SubmitWithdrawSuccess       Response
	SubmitWithdrawRequest       ErrResponse
	GetCreditAvailableSuccess   Response
	GetLoanSuccess              Response
	BorrowLoanSuccess           Response
	BorrowLoanRequest           ErrResponse
	GetInterestTermSuccess      Response
	GetRepaymentSuccess         Response
	GetRepaymentRequest         ErrResponse
	SubmitRepaymentSuccess      Response
	SubmitRepaymentRequest      ErrResponse
	//// Admin
	GetAccountAdminSuccess            Response
	GetAccountAdminRequest            ErrResponse
	ConfirmAccountAdminSuccess        Response
	ConfirmAccountAdminRequest        ErrResponse
	RejectAccountAdminSuccess         Response
	RejectAccountAdminRequest         ErrResponse
	UpdateAccountDocumentAdminSuccess Response
	UpdateAccountDocumentAdminRequest ErrResponse
	GetWalletTransactionAdminSuccess  Response
	GetWalletTransactionAdminRequest  ErrResponse
	ConfirmDepositAdminSuccess        Response
	ConfirmDepositAdminRequest        ErrResponse
	RejectDepositAdminSuccess         Response
	RejectDepostiAdminRequest         ErrResponse
	ConfirmWithdrawAdminSuccess       Response
	ConfirmWithdrawAdminRequest       ErrResponse
	RejectWithdrawAdminSuccess        Response
	RejectWithdrawAdminRequest        ErrResponse
	GetContractAdminSuccess           Response
	GetContractAdminRequest           ErrResponse
	ConfirmContractAdminSuccess       Response
	ConfirmContractAdminRequest       ErrResponse
	CreateInterestTermAdminSuccess    Response
	CreateInterestTermAdminRequest    ErrResponse
	UpdateInterestTermAdminSuccess    Response
	UpdateInterestTermAdminRequest    ErrResponse
	GetRepaymentAdminSuccess          Response
	GetRepaymentAdminRequest          ErrResponse
	ConfirmRepaymentAdminSuccess      Response
	ConfirmRepaymentAdminRequest      ErrResponse
	RejectRepaymentAdminSuccess       Response
	RejectRepaymentAdminRequest       ErrResponse
	LiquidateFundAdminSuccess         Response
	LiquidateFundAdminRequest         ErrResponse
	// Mail
	GetOTPSuccess       Response
	GetOTPRequest       ErrResponse
	GetOTPThirdParty    ErrResponse
	OTPRequestInvalid   ErrResponse
	OTPRequestFailLimit ErrResponse
	OTPRequestMaxLimit  ErrResponse
	OTPRequestDuplicate ErrResponse
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
