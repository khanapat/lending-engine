package lending

import (
	"encoding/json"
	"fmt"
	"lending-engine/client"
	"net/http"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type RequestLiquidationClientFn func(logger *zap.Logger, xRequestID string, request *SendLiquidationClientRequest) error

func NewRequestLiquidationClientFn(cli *client.Client) RequestLiquidationClientFn {
	return func(logger *zap.Logger, xRequestID string, request *SendLiquidationClientRequest) error {
		byteRequest, err := json.Marshal(&request)
		if err != nil {
			return err
		}
		m := make(map[string]string)
		clientRequest := client.Request{
			URL:                 viper.GetString("client.email-api.liquidation.url"),
			Method:              http.MethodPost,
			XRequestID:          xRequestID,
			Header:              m,
			HideLogRequestBody:  viper.GetBool("client.hidebody"),
			HideLogResponseBody: viper.GetBool("client.hidebody"),
			Logger:              logger,
			Body:                byteRequest,
		}
		clientResponse, err := cli.Do(&clientRequest)
		if err != nil {
			return err
		}
		var sendLiquidationClientResult SendLiquidationClientResult
		if err := json.Unmarshal(clientResponse.Body, &sendLiquidationClientResult); err != nil {
			return err
		}
		if sendLiquidationClientResult.Code != 2000 {
			return fmt.Errorf("%s(%s)", sendLiquidationClientResult.Title, sendLiquidationClientResult.Description)
		}
		return nil
	}
}
