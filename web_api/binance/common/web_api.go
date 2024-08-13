package common_web_api

import (
	"sync"

	request "github.com/fr0ster/turbo-cambitor/web_api/binance/common/request"

	signature "github.com/fr0ster/turbo-restler/utils/signature"
	web_api "github.com/fr0ster/turbo-restler/web_api"
)

func (wa *WebApi) Lock() {
	wa.mutex.Lock()
}

func (wa *WebApi) Unlock() {
	wa.mutex.Unlock()
}

func (wa *WebApi) PlaceOrder() *request.Request {
	return request.New(wa.apiKey, wa.symbol, "order.place", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) CancelOrder() *request.Request {
	return request.New(wa.apiKey, wa.symbol, "order.cancel", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) QueryOrder() *request.Request {
	return request.New(wa.apiKey, wa.symbol, "order.status", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) CancelReplaceOrder() *request.Request {
	return request.New(wa.apiKey, wa.symbol, "order.cancelReplace", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) QueryOpenOrders() *request.Request {
	return request.New(wa.apiKey, wa.symbol, "openOrders.status", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) QueryAllOrders() *request.Request {
	return request.New(wa.apiKey, wa.symbol, "orderList.status", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) ListOfSubscriptions() *request.Request {
	return request.New(wa.apiKey, wa.symbol, "LIST_SUBSCRIPTIONS", wa.waHost, wa.waPath, wa.sign)
}

// Функція для логіну
func (wa *WebApi) Logon() *request.Request {
	return request.New(wa.apiKey, wa.symbol, "session.logon", wa.waHost, wa.waPath, wa.sign)
}

// Функція для логіну
func (wa *WebApi) Logout() *request.Request {
	return request.New(wa.apiKey, wa.symbol, "session.logout", wa.waHost, wa.waPath, wa.sign)
}

// Функція для перевірки статусу сесії
func (wa *WebApi) Status() *request.Request {
	return request.New(wa.apiKey, wa.symbol, "session.status", wa.waHost, wa.waPath, wa.sign)
}

func New(apiKey, apiSecret string, host web_api.WsHost, path web_api.WsPath, symbol string, sign signature.Sign) *WebApi {
	return &WebApi{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		symbol:    symbol,
		waHost:    host,
		waPath:    path,
		mutex:     &sync.Mutex{},
		sign:      sign,
	}
}
