package mail

type ReferenceData struct {
	Otp       string `json:"otp" example:"999999"`
	FailCount int    `json:"failCount" example:"1"`
}

type SendMailOtpClientRequest struct {
	From     string                `json:"from" example:"Treasury.Admin@gmail.com"`
	To       []string              `json:"to" example:"[k.apiwattanawong@gmail.com]"`
	Subject  string                `json:"subject" example:"Request OTP"`
	Template string                `json:"template" example:"otp.html"`
	Body     BodySendMailOtpClient `json:"body" example:"999999"`
	Auth     bool                  `json:"auth" example:"false"`
}

type BodySendMailOtpClient struct {
	UserName string `json:"userName" example:"trust"`
	RefNo    string `json:"refNo" example:"tog2C7"`
	Otp      string `json:"otp" example:"9999"`
}

type OtpMailResponse struct {
	ReferenceNo string `json:"referenceNo,omitempty" example:"999999"`
	ExpiredTime string `json:"expiredTime" example:"2020-01-01 12:12:12"`
}

type SendMailOtpClientResult struct {
	Code        uint64 `json:"code" example:"2000"`
	Title       string `json:"title" example:"Success."`
	Description string `json:"description" example:"Please contact administrator for more information."`
}
