package response

const (
	SuccessCode                uint64 = 200
	ErrInvalidRequestCode      uint64 = 1000
	ErrRequestExpireCode       uint64 = 1001
	ErrUnauthorizationCode     uint64 = 4001
	ErrOTPRequestCode          uint64 = 4002
	ErrBasicAuthenticationCode uint64 = 4007
	ErrDatabaseCode            uint64 = 5000
	ErrRedisCode               uint64 = 5001
	ErrOperationCode           uint64 = 5002
	ErrBlockchainCode          uint64 = 5003
	ErrThirdPartyCode          uint64 = 5004
)

const (
	SuccessMessageEN           string = "Success."
	ErrInternalServerMessageEN string = "Internal server error."
	// Account
	SuccessSignUpMessageEN               string = "Success sign up account."
	ErrSignUpMessageEN                   string = "Cannot sign up account."
	SuccessConfirmVerifyEmailMessageEN   string = "Success verify email."
	ErrConfirmVerifyEmailMessageEN       string = "Cannot verify email."
	SuccessLoginMessageEN                string = "Success login account."
	ErrLoginMessageEN                    string = "Cannot login account."
	SuccessAcceptTermsConditionMessageEN string = "Success accept terms & condition."
	ErrAcceptTermsConditionMessageEN     string = "Cannot accept terms & condition."
	SuccessGetTermsConditionMessageEN    string = "Success get terms & condition."
	ErrGetTermsConditionMessageEN        string = "Cannot get terms & condition."
	SuccessRequestResetPasswordMessageEN string = "Success request reset password."
	ErrRequestResetPasswordMessageEN     string = "Cannot request reset password."
	SuccessResetPasswordMessageEN        string = "Success reset password."
	ErrResetPasswordMessageEN            string = "Cannot reset password."
	// Lending
	//// User
	SuccessGetToknPriceMessageEN         string = "Success get token price."
	SuccessPreCalculationLoanMessageEN   string = "Success calculate loan."
	ErrPreCalculationLoanMessageEN       string = "Cannot calculate loan."
	SuccessGetWalletTransactionMessageEN string = "Success get wallet transaction."
	ErrGetWalletTransactionMessageEN     string = "Cannot get wallet transaction."
	SuccessSubmitDepositMessageEN        string = "Success submit deposit token."
	ErrSubmitDepositMessageEN            string = "Cannot submit deposit token."
	SuccessSubmitWithdrawMessageEN       string = "Success submit withdraw token."
	ErrSubmitWithdrawMessageEN           string = "Cannot submit withdraw token."
	SuccessGetCreditAvailableMessageEN   string = "Success get credit available."
	SuccessGetLoanMessageEN              string = "Success get loan."
	ErrGetLoanMessageEN                  string = "Cannot get loan."
	SuccessBorrowLoanMessageEN           string = "Success borrow loan."
	ErrBorrowLoanMessageEN               string = "Cannot borrow loan."
	SuccessGetInterestTermMessageEN      string = "Success get interest term."
	SuccessGetRepaymentMessageEN         string = "Success get repayment."
	ErrGetRepaymentMessageEN             string = "Cannot get repayment."
	SuccessSubmitRepaymentMessageEN      string = "Success submit repayment."
	ErrSubmitRepaymentMessageEN          string = "Cannot submit repayment."
	//// Admin
	SuccessGetAccountAdminMessageEN           string = "Success get account detail."
	ErrGetAccountAdminMessageEN               string = "Cannot get account detail."
	SuccessConfirmAccountAdminMessageEN       string = "Success confirm account detail."
	ErrConfirmAccountAdminMessageEN           string = "Cannot confirm account detail."
	SuccessGetWalletTransactionAdminMessageEN string = "Success get wallet transaction."
	ErrGetWalletTransactionAdminMessageEN     string = "Cannot get wallet transaction."
	SuccessConfirmDepositAdminMessageEN       string = "Success confirm deposit token."
	ErrConfirmDepositAdminMessageEN           string = "Cannot confirm deposit token."
	SuccessConfirmWithdrawAdminMessageEN      string = "Success confirm withdraw token."
	ErrConfirmWithdrawAdminMessageEN          string = "Cannot confirm withdraw token."
	SuccessGetContractAdminMessageEN          string = "Success get loan contract."
	ErrGetContractAdminMessageEN              string = "Cannot get loan contract."
	SuccessConfirmContractAdminMessageEN      string = "Success confirm loan contract."
	ErrConfirmContractAdminMessageEN          string = "Cannot confirm loan contract."
	SuccessGetRepaymentAdminMessageEN         string = "Success get repayment."
	ErrGetRepaymentAdminMessageEN             string = "Cannot get repayment."
	SuccessConfirmRepaymentAdminMessageEN     string = "Success confirm repayment."
	ErrConfirmRepaymentAdminMessageEN         string = "Cannot confirm repayment."
	// Mail
	SuccessOTPRequestMessageEN      string = "Success request otp."
	ErrOTPRequestMessageEN          string = "Cannot request otp."
	ErrInvalidOTPMessageEN          string = "OTP is invalid."
	ErrLimitInvalidOTPMessageEN     string = "Maximum failed otp attempts."
	ErrLimitOTPRequestMessageEN     string = "Maximum daily requested otp."
	ErrDuplicateOTPRequestMessageEN string = "OTP request is duplicate."
	// BasicAuthen
	ErrBasicAuthenticationMessageEN string = "Authentication failed."
	// AuthorizeToken
	ErrAuthorizationTokenMessageEN string = "Unauthorization token."
	// Desc
	ErrCooldownDescEN       string = "Please try again later."
	ErrRequestDataDescEN    string = "Please check request data again."
	ErrContactAdminDescEN   string = "Please contact administrator for more information."
	ErrThirdPartyDescEN     string = "Service is unavailable. Please try again later."
	ErrAuthenticationDescEN string = "Unable to access data. Please check user & password."
	ErrAuthorizationDescEN  string = "Unauthorized service. Please check access token."
)

const (
	SuccessMessageTH           string = "สำเร็จ."
	ErrInternalServerMessageTH string = "มีข้อผิดพลาดภายในเซิร์ฟเวอร์."
	// Account
	SuccessSignUpMessageTH               string = "สมัครบัญชีเข้าใช้งานสำเร็จ."
	ErrSignUpMessageTH                   string = "ไม่สามารถสมัครบัญชีเข้าใช้งานได้."
	SuccessConfirmVerifyEmailMessageTH   string = "ยืนยันอีเมล์สำเร็จ."
	ErrConfirmVerifyEmailMessageTH       string = "ไม่สามารถยืนยันอีเมล์ได้."
	SuccessLoginMessageTH                string = "เข้าสู่บัญชีผู้ใช้งานสำเร็จ."
	ErrLoginMessageTH                    string = "ไม่สามารถเข้าสู่บัญชีผู้ใช้งานได้."
	SuccessAcceptTermsConditionMessageTH string = "ยอมรับข้อกำหนดและเงื่อนไขสำเร็จ."
	ErrAcceptTermsConditionMessageTH     string = "ไม่สามารถยอมรับข้อกำหนดและเงื่อนไขได้."
	SuccessGetTermsConditionMessageTH    string = "แสดงการยอมรับข้อกำหนดและเงื่อนไขสำเร็จ."
	ErrGetTermsConditionMessageTH        string = "ไม่สามารถแสดงการยอมรับข้อกำหนดและเงื่อนไขได้."
	SuccessRequestResetPasswordMessageTH string = "ส่งคำขอร้องเพื่อตั้งรหัสบัญชีใหม่สำเร็จ."
	ErrRequestResetPasswordMessageTH     string = "ไม่สามารถส่งคำขอร้องเพื่อตั้งรหัสบัญชีใหม่ได้."
	SuccessResetPasswordMessageTH        string = "ตั้งรหัสบัญชีใหม่สำเร็จ."
	ErrResetPasswordMessageTH            string = "ไม่สามารถตั้งรหัสบัญชีใหม่ได้."
	// Lending
	//// User
	SuccessGetToknPriceMessageTH         string = "แสดงราคาซื้อขายโทเคนสำเร็จ."
	SuccessPreCalculationLoanMessageTH   string = "คำนวณอัตราเงินกู้สำเร็จ."
	ErrPreCalculationLoanMessageTH       string = "ไม่สามารถคำนวณอัตราเงินกู้ได้."
	SuccessGetWalletTransactionMessageTH string = "แสดงรายการฝากถอนโทเคนสำเร็จ."
	ErrGetWalletTransactionMessageTH     string = "ไม่สามารถแสดงรายการฝากถอนโทเคนได้."
	SuccessSubmitDepositMessageTH        string = "ส่งหลักฐานยืนยันการฝากโทเคนสำเร็จ."
	ErrSubmitDepositMessageTH            string = "ไม่สามารถส่งหลักฐานยืนยันการฝากโทเคนได้."
	SuccessSubmitWithdrawMessageTH       string = "ส่งคำร้องขอถอนโทเคนสำเร็จ."
	ErrSubmitWithdrawMessageTH           string = "ไม่สามารถส่งคำร้องขอถอนโทเคนได้."
	SuccessGetCreditAvailableMessageTH   string = "แสดงเครดิตคงเหลือสำเร็จ."
	SuccessGetLoanMessageTH              string = "แสดงการกู้ยืมเงินสำเร็จ."
	ErrGetLoanMessageTH                  string = "ไม่สามารถแสดงการกู้ยืมเงินได้."
	SuccessBorrowLoanMessageTH           string = "กู้ยืมเงินสำเร็จ."
	ErrBorrowLoanMessageTH               string = "ไม่สามารถกู้ยืมเงินได้."
	SuccessGetInterestTermMessageTH      string = "แสดงอัตราดอกเบี้ยสำเร็จ."
	SuccessGetRepaymentMessageTH         string = "แสดงรายการจ่ายเงินคืนสำเร็จ."
	ErrGetRepaymentMessageTH             string = "ไม่สามารถแสดงรายการจ่ายเงินคืนได้."
	SuccessSubmitRepaymentMessageTH      string = "ส่งหลักฐานยืนยันการจ่ายเงินคืนสำเร็จ."
	ErrSubmitRepaymentMessageTH          string = "ไม่สามารถส่งหลักฐานยืนยันการจ่ายเงินคืนได้."
	//// Admin
	SuccessGetAccountAdminMessageTH           string = "แสดงข้อมูลบัญชีผู้ใช้งานสำเร็จ."
	ErrGetAccountAdminMessageTH               string = "ไม่สามารถแสดงข้อมูลบัญชีผู้ใช้งานได้."
	SuccessConfirmAccountAdminMessageTH       string = "ยืนยันข้อมูลบัญชีผู้ใช้งานสำเร็จ."
	ErrConfirmAccountAdminMessageTH           string = "ไม่สามารถยืนยันข้อมูลบัญชีผู้ใช้งานได้."
	SuccessGetWalletTransactionAdminMessageTH string = "แสดงรายการฝากถอนโทเคนสำเร็จ."
	ErrGetWalletTransactionAdminMessageTH     string = "ไม่สามารถแสดงรายการฝากถอนโทเคนได้."
	SuccessConfirmDepositAdminMessageTH       string = "ยืนยันการฝากโทเคนสำเร็จ."
	ErrConfirmDepositAdminMessageTH           string = "ไม่สามารถยืนยันการฝากโทเคนได้."
	SuccessConfirmWithdrawAdminMessageTH      string = "ยืนยันการถอนโทเคนสำเร็จ."
	ErrConfirmWithdrawAdminMessageTH          string = "ไม่สามารถยืนยันการถอนโทเคนได้."
	SuccessGetContractAdminMessageTH          string = "แสดงสัญญากู้ยืมสำเร็จ."
	ErrGetContractAdminMessageTH              string = "ไม่สามารถแสดงสัญญากู้ยืมได้."
	SuccessConfirmContractAdminMessageTH      string = "ยืนยันการกู้ยืมสำเร็จ."
	ErrConfirmContractAdminMessageTH          string = "ไม่สามารถยืนยันการกู้ยืมได้."
	SuccessGetRepaymentAdminMessageTH         string = "แสดงรายการจ่ายเงินคืนสำเร็จ."
	ErrGetRepaymentAdminMessageTH             string = "ไม่สามารถแสดงรายการจ่ายเงินคืนได้."
	SuccessConfirmRepaymentAdminMessageTH     string = "ยืนยันการจ่ายเงินคืนสำเร็จ."
	ErrConfirmRepaymentAdminMessageTH         string = "ไม่สามารถยืนยันการจ่ายเงินคืนได้."
	// Mail
	SuccessOTPRequestMessageTH      string = "ขอรหัส OTP สำเร็จ."
	ErrOTPRequestMessageTH          string = "ไม่สามารถขอรหัส OTP ได้."
	ErrInvalidOTPMessageTH          string = "รหัส OTP ไม่ถูกต้อง."
	ErrLimitInvalidOTPMessageTH     string = "กรอกรหัส OTP ผิดเกินจำนวนที่กำหนด."
	ErrLimitOTPRequestMessageTH     string = "ขอรหัส OTP มากเกินจำนวนที่กำหนดต่อวัน."
	ErrDuplicateOTPRequestMessageTH string = "ขอรหัส OTP ซ้ำ."
	// BasicAuthen
	ErrBasicAuthenticationMessageTH string = "ยืนยันตัวตนล้มเหลว."
	// AuthorizeToken
	ErrAuthorizationTokenMessageTH string = "ตรวจสอบสิทธิ์ล้มเหลว."
	// Desc
	ErrCooldownDescTH       string = "กรุณาทำรายการใหม่อีกครั้งภายหลัง."
	ErrRequestDataDescTH    string = "กรุณาตรวจสอบข้อมูลอีกครั้ง."
	ErrContactAdminDescTH   string = "กรุณาติดต่อเจ้าหน้าที่ดูแลระบบเพื่อรับข้อมูลเพิ่มเติม."
	ErrThirdPartyDescTH     string = "ไม่สามารถใช้บริการได้. กรุณาทำรายการใหม่อีกครั้งภายหลัง."
	ErrAuthenticationDescTH string = "ไม่สามารถเข้าถึงข้อมูลได้. กรุณาตรวจสอบรหัสผู้ใช้งานใหม่อีกครั้ง."
	ErrAuthorizationDescTH  string = "ไม่สามารถใช้งานระบบได้. กรุณาตรวจสอบสิทธิ์การเข้าใช้งานอีกครั้ง."
)

const (
	ValidateFieldError string = "Invalid Parameters"
	OperationError     string = "Invalid Operation"
)
