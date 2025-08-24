package spot_web_api

import (
	"time"

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
		waHost = "ws-api.testnet.binance.vision"
		waEndpoint = "/ws-api/v3"
		waScheme = common.WsSchemeWSS
	} else {
		waHost = "ws-api.binance.com"
		waEndpoint = "/ws-api/v3"
		waScheme = common.WsSchemeWSS
	}
	return web_api.New(waHost, waEndpoint, waScheme, sign, web_api.WithTimeouts(30*time.Second, 30*time.Second))
}

func New(host string, endpoint string, scheme string, sign signature.Sign) WebApi {
	return web_api.New(
		common.WsHost(host),
		common.WsEndpoint(endpoint),
		common.WsScheme(scheme),
		sign,
	)
}

// NewDefaultWithOptions mirrors NewDefault but allows passing functional options to override the socket factory, timeouts, etc.
func NewDefaultWithOptions(sign signature.Sign, useTestNet bool, opts ...web_api.Option) WebApi {
	var (
		waHost     common.WsHost
		waEndpoint common.WsEndpoint
		waScheme   = common.WsSchemeWSS
	)
	if useTestNet {
		waHost = "ws-api.testnet.binance.vision"
		waEndpoint = "/ws-api/v3"
	} else {
		waHost = "ws-api.binance.com"
		waEndpoint = "/ws-api/v3"
	}
	opts = append([]web_api.Option{web_api.WithTimeouts(30*time.Second, 30*time.Second)}, opts...)
	return web_api.New(waHost, waEndpoint, waScheme, sign, opts...)
}

// NewWithOptions mirrors New but allows passing functional options to override the socket factory, timeouts, etc.
func NewWithOptions(host string, endpoint string, scheme string, sign signature.Sign, opts ...web_api.Option) WebApi {
	return web_api.New(
		common.WsHost(host),
		common.WsEndpoint(endpoint),
		common.WsScheme(scheme),
		sign,
		opts...,
	)
}
