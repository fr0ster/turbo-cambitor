package futures_web_api

import (
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/fr0ster/turbo-cambitor/common"
	request "github.com/fr0ster/turbo-cambitor/web_api/binance/common/request"
	web_api "github.com/fr0ster/turbo-cambitor/web_api/binance/common/web_api"

	signature "github.com/fr0ster/turbo-signer/v2/signature"
)

// WithTimeouts експортується для тестів та користувачів API
var WithTimeouts = web_api.WithTimeouts

type WebApi interface {
	AccountBalance() *request.RequestBuilder
	AccountInformation() *request.RequestBuilder
	AccountPositions() *request.RequestBuilder
	CancelOrder() *request.RequestBuilder
	ModifyOrder() *request.RequestBuilder
	Logon() *request.RequestBuilder
	Logout() *request.RequestBuilder
	OrderBook() *request.RequestBuilder
	Ping() *request.RequestBuilder
	PlaceOrder() *request.RequestBuilder
	QueryOrder() *request.RequestBuilder
	QueryPosition() *request.RequestBuilder
	QueryPositionV2() *request.RequestBuilder
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
		waHost = "testnet.binancefuture.com"
		waEndpoint = "/ws-fapi/v1"
		waScheme = common.WsSchemeWSS
	} else {
		waHost = "ws-fapi.binance.com"
		waEndpoint = "/ws-fapi/v1"
		waScheme = common.WsSchemeWSS
	}
	return web_api.New(waHost, waEndpoint, waScheme, sign, WithTimeouts(30*time.Second, 30*time.Second))
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
		waHost = "testnet.binancefuture.com"
		waEndpoint = "/ws-fapi/v1"
	} else {
		waHost = "ws-fapi.binance.com"
		waEndpoint = "/ws-fapi/v1"
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
