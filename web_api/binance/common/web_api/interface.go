package common_web_api

import (
	"github.com/bitly/go-simplejson"
	request "github.com/fr0ster/turbo-cambitor/web_api/binance/common/request"
)

type WebApiInterface interface {
	// Account-related
	AccountBalance() *request.RequestBuilder
	AccountInformation() *request.RequestBuilder
	AccountPositions() *request.RequestBuilder
	QueryPosition() *request.RequestBuilder
	QueryPositionV2() *request.RequestBuilder

	// Order-related
	PlaceOrder() *request.RequestBuilder
	CancelOrder() *request.RequestBuilder
	QueryOrder() *request.RequestBuilder
	CancelReplaceOrder() *request.RequestBuilder
	ModifyOrder() *request.RequestBuilder
	QueryOpenOrders() *request.RequestBuilder
	QueryAllOrders() *request.RequestBuilder

	// Session
	Logon() *request.RequestBuilder
	Logout() *request.RequestBuilder
	Status() *request.RequestBuilder

	// Market Data
	OrderBook() *request.RequestBuilder
	SymbolPriceTicker() *request.RequestBuilder
	SymbolBookTicker() *request.RequestBuilder
	Ping() *request.RequestBuilder
	Time() *request.RequestBuilder
	ExchangeInfo() *request.RequestBuilder

	// User Stream
	UserDataStreamStart() *request.RequestBuilder
	UserDataStreamPing() *request.RequestBuilder
	UserDataStreamStop() *request.RequestBuilder

	// Utility
	ListOfSubscriptions() *request.RequestBuilder

	// Concurrency
	Lock()
	Unlock()

	// Generic call
	Call(js *simplejson.Json) (*simplejson.Json, error)
}
