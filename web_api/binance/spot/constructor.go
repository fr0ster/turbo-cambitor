package spot_web_api

import (
	"github.com/bitly/go-simplejson"
	request "github.com/fr0ster/turbo-cambitor/web_api/binance/common/request"
	web_api "github.com/fr0ster/turbo-cambitor/web_api/binance/common/web_api"

	"github.com/fr0ster/turbo-restler/web_socket"
	signature "github.com/fr0ster/turbo-signer/signature"
)

type WebApi interface {
	AccountInformation() *request.RequestBuilder
	CancelOrder() *request.RequestBuilder
	CancelReplaceOrder() *request.RequestBuilder
	ExchangeInfo() *request.RequestBuilder
	Logon() *request.RequestBuilder
	Logout() *request.RequestBuilder
	OrderBook() *request.RequestBuilder
	Ping() *request.RequestBuilder
	PlaceOrder() *request.RequestBuilder
	QueryAllOrders() *request.RequestBuilder
	QueryOpenOrders() *request.RequestBuilder
	QueryOrder() *request.RequestBuilder
	Status() *request.RequestBuilder
	SymbolBookTicker() *request.RequestBuilder
	SymbolPriceTicker() *request.RequestBuilder
	Time() *request.RequestBuilder

	UserDataStreamStart() *request.RequestBuilder
	UserDataStreamPing() *request.RequestBuilder
	UserDataStreamStop() *request.RequestBuilder

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
