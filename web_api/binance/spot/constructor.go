package spot_web_api

import (
	"github.com/bitly/go-simplejson"
	request "github.com/fr0ster/turbo-cambitor/web_api/binance/common/request"
	web_api "github.com/fr0ster/turbo-cambitor/web_api/binance/common/web_api"

	"github.com/fr0ster/turbo-restler/web_socket"
	signature "github.com/fr0ster/turbo-signer/signature"
)

type WebApi interface {
	AccountInformation() *request.Request
	CancelOrder() *request.Request
	CancelReplaceOrder() *request.Request
	ExchangeInfo() *request.Request
	Logon() *request.Request
	Logout() *request.Request
	OrderBook() *request.Request
	Ping() *request.Request
	PlaceOrder() *request.Request
	QueryAllOrders() *request.Request
	QueryOpenOrders() *request.Request
	QueryOrder() *request.Request
	Status() *request.Request
	SymbolBookTicker() *request.Request
	SymbolPriceTicker() *request.Request
	Time() *request.Request

	UserDataStreamStart() *request.Request
	UserDataStreamPing() *request.Request
	UserDataStreamStop() *request.Request

	Call(*simplejson.Json) (*simplejson.Json, error)

	Lock()
	Unlock()
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
		waHost = "testnet.binance.vision"
		waPath = "/ws-api/v3"
	} else {
		waHost = "ws-api.binance.com:443"
		waPath = "/ws-api/v3"
	}
	return web_api.New(web_socket.WsHost(waHost), web_socket.WsPath(waPath), sign)
}
