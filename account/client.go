package account

import (
	"encoding/json"
	"fmt"
	"lending-engine/client"
	"net/http"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type RequestVerifyEmailClientFn func(logger *zap.Logger, XRequestID string, request *SendVerifyEmailClientRequest) error

func NewRequestVerifyEmailClientFn(cli *client.Client) RequestVerifyEmailClientFn {
	return func(logger *zap.Logger, xRequestID string, request *SendVerifyEmailClientRequest) error {
		byteRequest, err := json.Marshal(&request)
		if err != nil {
			return err
		}
		m := make(map[string]string)
		clientRequest := client.Request{
			URL:                 viper.GetString("client.email-api.url"),
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
		var sendVerifyEmailClientResult SendVerifyEmailClientResult
		if err := json.Unmarshal(clientResponse.Body, &sendVerifyEmailClientResult); err != nil {
			return err
		}
		if sendVerifyEmailClientResult.Code != 2000 {
			return fmt.Errorf("%s(%s)", sendVerifyEmailClientResult.Title, sendVerifyEmailClientResult.Description)
		}
		return nil
	}
}
