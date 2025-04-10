package common_web_api

import (
	request "github.com/fr0ster/turbo-cambitor/web_api/binance/common/request"
)

func (wa *WebApiWrapper) ListOfSubscriptions() *request.RequestBuilder {
	return request.New("LIST_SUBSCRIPTIONS", wa.waHost, wa.waEndpoint, wa.sign)
}

// Функція для логіну
func (wa *WebApiWrapper) Logon() *request.RequestBuilder {
	return request.New("session.logon", wa.waHost, wa.waEndpoint, wa.sign)
}

// Функція для логіну
func (wa *WebApiWrapper) Logout() *request.RequestBuilder {
	return request.New("session.logout", wa.waHost, wa.waEndpoint, wa.sign)
}

// Функція для перевірки статусу сесії
func (wa *WebApiWrapper) Status() *request.RequestBuilder {
	return request.New("session.status", wa.waHost, wa.waEndpoint, wa.sign)
}

func (wa *WebApiWrapper) OrderBook() *request.RequestBuilder {
	return request.New("depth", wa.waHost, wa.waEndpoint, wa.sign)
}

func (wa *WebApiWrapper) SymbolPriceTicker() *request.RequestBuilder {
	return request.New("ticker.price", wa.waHost, wa.waEndpoint, wa.sign)
}

func (wa *WebApiWrapper) SymbolBookTicker() *request.RequestBuilder {
	return request.New("ticker.book", wa.waHost, wa.waEndpoint, wa.sign)
}

func (wa *WebApiWrapper) Ping() *request.RequestBuilder {
	return request.New("ping", wa.waHost, wa.waEndpoint, wa.sign)
}

func (wa *WebApiWrapper) Time() *request.RequestBuilder {
	return request.New("time", wa.waHost, wa.waEndpoint, wa.sign)
}

func (wa *WebApiWrapper) ExchangeInfo() *request.RequestBuilder {
	return request.New("exchangeInfo", wa.waHost, wa.waEndpoint, wa.sign)
}

func (wa *WebApiWrapper) UserDataStreamStart() *request.RequestBuilder {
	return request.New("userDataStream.start", wa.waHost, wa.waEndpoint, wa.sign)
}

// UserDataStreamPing(listenKey string) (newListenKey string, err error)
func (wa *WebApiWrapper) UserDataStreamPing() *request.RequestBuilder {
	return request.New("userDataStream.ping", wa.waHost, wa.waEndpoint, wa.sign)
}

// UserDataStreamStop(listenKey string) (err error)
func (wa *WebApiWrapper) UserDataStreamStop() *request.RequestBuilder {
	return request.New("userDataStream.stop", wa.waHost, wa.waEndpoint, wa.sign)
}
