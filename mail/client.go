package mail

import (
	"encoding/json"
	"fmt"
	"lending-engine/client"
	"net/http"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type RequestMailOtpClientFn func(logger *zap.Logger, xRequestID string, request *SendMailOtpClientRequest) error

func NewRequestMailOtpClientFn(cli *client.Client) RequestMailOtpClientFn {
	return func(logger *zap.Logger, xRequestID string, request *SendMailOtpClientRequest) error {
		byteRequest, err := json.Marshal(&request)
		if err != nil {
			return err
		}
		m := make(map[string]string, 0)
		clientRequest := client.Request{
			URL:                 viper.GetString("client.email-api.otp.url"),
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
		var sendMailOtpClientResult SendMailOtpClientResult
		if err := json.Unmarshal(clientResponse.Body, &sendMailOtpClientResult); err != nil {
			return err
		}
		if sendMailOtpClientResult.Code != 2000 {
			return fmt.Errorf("%s(%s)", sendMailOtpClientResult.Title, sendMailOtpClientResult.Description)
		}
		return nil
	}
}
