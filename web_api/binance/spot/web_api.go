package spot_web_api

import (
	common "github.com/fr0ster/turbo-cambitor/web_api/binance/common"
	request "github.com/fr0ster/turbo-cambitor/web_api/binance/common/request"

	signature "github.com/fr0ster/turbo-restler/utils/signature"
	web_api "github.com/fr0ster/turbo-restler/web_api"
)

type WebApi interface {
	PlaceOrder() *request.Request
	CancelOrder() *request.Request
	QueryOrder() *request.Request
	CancelReplaceOrder() *request.Request
	QueryOpenOrders() *request.Request
	QueryAllOrders() *request.Request
	ListOfSubscriptions() *request.Request
	Logon() *request.Request
	Logout() *request.Request
	Status() *request.Request
}

func New(apiKey, apiSecret, symbol string, sign signature.Sign, useTestNet ...bool) WebApi {
	var (
		waHost string
		waPath string
	)
	if len(useTestNet) == 0 {
		useTestNet = append(useTestNet, false)
	}
	if useTestNet[0] {
		waHost = "testnet.binance.vision"
		waPath = "/ws-api/v3"
	} else {
		waHost = "ws-api.binance.com:443"
		waPath = "/ws-api/v3"
	}
	return common.New(apiKey, apiSecret, web_api.WsHost(waHost), web_api.WsPath(waPath), symbol, sign)
}
