package futures_web_api

import (
	common "github.com/fr0ster/turbo-cambitor/web_api/binance/common"
	request "github.com/fr0ster/turbo-cambitor/web_api/binance/common/request"

	web_api "github.com/fr0ster/turbo-restler/web_api"
	signature "github.com/fr0ster/turbo-signer/signature"
)

type WebApi interface {
	AccountBalance() *request.Request
	AccountBalanceV2() *request.Request
	AccountInformation() *request.Request
	AccountInformationV2() *request.Request
	CancelOrder() *request.Request
	ModifyOrder() *request.Request
	Logon() *request.Request
	Logout() *request.Request
	OrderBook() *request.Request
	Ping() *request.Request
	PlaceOrder() *request.Request
	QueryOrder() *request.Request
	QueryPosition() *request.Request
	QueryPositionV2() *request.Request
	Status() *request.Request
	SymbolBookTicker() *request.Request
	SymbolPriceTicker() *request.Request
	Time() *request.Request
}

func New(sign signature.Sign, useTestNet ...bool) WebApi {
	var (
		waHost string
		waPath string
	)
	if len(useTestNet) == 0 {
		useTestNet = append(useTestNet, false)
	}
	if useTestNet[0] {
		waHost = "testnet.binancefuture.com"
		waPath = "/ws-fapi/v1"
	} else {
		waHost = "ws-fapi.binance.com"
		waPath = "/ws-fapi/v1"
	}
	return common.New(web_api.WsHost(waHost), web_api.WsPath(waPath), sign)
}
