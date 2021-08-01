package common

const (
	RequestInfoMsg    string = "Request Information"
	ResponseInfoMsg   string = "Response Information"
	ClientRequestMsg  string = "http client request info"
	ClientResponseMsg string = "http client response info"
)

const (
	ContentType     string = "Content-Type"
	ApplicationJSON string = "application/json"
	XRequestID      string = "X-Request-ID"
	LocaleKey       string = "locale"
	JWTClaimsKey    string = "claims"
	ReferenceOTPKey string = "ReferenceNo"
	OTPKey          string = "OTP"
)

const (
	DateYYYYMMDDHHMMSSFormat string = "2006-01-02 15:04:05"
	RandomStringInt          string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	RandomInt                string = "0123456789"
)

const (
	THBBTCRedis    string = "THB/BTC"
	THBETHRedis    string = "THB/ETH"
	PendingStatus  string = "PENDING"
	ConfirmStatus  string = "CONFIRMED"
	RejectStatus   string = "REJECTED"
	OngoingStatus  string = "ONGOING"
	ClosedStatus   string = "CLOSED"
	DepositStatus  string = "DEPOSIT"
	WithdrawStatus string = "WITHDRAW"
)

const (
	PenaltyRedis string = "Penalty"
)
