package blockchain

type TransactionInfo struct {
	TxnHash        string        `json:"txnHash" example:"0xf5a3aa87c40b05e6a308b61186eeded8996b654a9895401b8089a2966b54f618"`
	Status         int64         `json:"status" example:"1"`
	Block          int64         `json:"block" example:"12870267"`
	Timestamp      int64         `json:"timestamp" example:"1527211625"`
	From           string        `json:"from" example:"0x0dcf57635f6562897cba35168b232fb302de0748"`
	InteractedWith string        `json:"interactWith" example:"0x2b54a9350de2bf0be86a09253d9382829e74084a"`
	TokenTransfer  TokenTransfer `json:"tokensTransfer"`
	Value          float64       `json:"value" example:"0.05"`
	TxnFee         float64       `json:"txnFee" example:"0.000462"`
	GasPrice       float64       `json:"gasPrice" example:"0.000000022"`
	GasLimit       int64         `json:"gasLimit" example:"21000"`
	GasUsed        int64         `json:"gasUsed" example:"21000"`
	Nonce          int64         `json:"nonce" example:"629"`
}

type TokenTransfer struct {
	From   string  `json:"from" example:"0xc083eb69aa7215f4afa7a22dcbfcc1a33999371c"`
	To     string  `json:"to" example:"0xa9b6d99ba92d7d691c6ef4f49a1dc909822cee46"`
	Amount float64 `json:"amount" example:"123456789"`
}
