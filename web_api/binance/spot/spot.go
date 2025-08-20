package spot_web_api

import (
	"github.com/bitly/go-simplejson"
	"github.com/fr0ster/turbo-cambitor/common"
	request "github.com/fr0ster/turbo-cambitor/web_api/binance/common/request"
	web_api "github.com/fr0ster/turbo-cambitor/web_api/binance/common/web_api"

	signature "github.com/fr0ster/turbo-signer/v2/signature"
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

	Call(js *simplejson.Json) (*simplejson.Json, error)

	Lock()
	Unlock()
}

func NewDefault(sign signature.Sign, useTestNet ...bool) WebApi {
	var (
		waHost     common.WsHost
		waEndpoint common.WsEndpoint
		waScheme   = common.WsSchemeWSS
	)
	if len(useTestNet) == 0 {
		useTestNet = append(useTestNet, false)
	}
	if useTestNet[0] {
		waHost = "testnet.binance.vision"
		waEndpoint = "/ws-api/v3"
		waScheme = common.WsSchemeWSS
	} else {
		waHost = "ws-api.binance.com:443"
		waEndpoint = "/ws-api/v3"
		waScheme = common.WsSchemeWSS
	}
	return web_api.New(waHost, waEndpoint, waScheme, sign)
}

func New(host string, endpoint string, scheme string, sign signature.Sign) WebApi {
	return web_api.New(
		common.WsHost(host),
		common.WsEndpoint(endpoint),
		common.WsScheme(scheme),
		sign,
	)
}
