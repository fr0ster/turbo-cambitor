package common_web_api

import (
	request "github.com/fr0ster/turbo-cambitor/web_api/binance/common/request"
)

func (wa *WebApiWrapper) ListOfSubscriptions() *request.Request {
	return request.New("LIST_SUBSCRIPTIONS", wa.waHost, wa.waPath, wa.sign)
}

// Функція для логіну
func (wa *WebApiWrapper) Logon() *request.Request {
	return request.New("session.logon", wa.waHost, wa.waPath, wa.sign)
}

// Функція для логіну
func (wa *WebApiWrapper) Logout() *request.Request {
	return request.New("session.logout", wa.waHost, wa.waPath, wa.sign)
}

// Функція для перевірки статусу сесії
func (wa *WebApiWrapper) Status() *request.Request {
	return request.New("session.status", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) OrderBook() *request.Request {
	return request.New("depth", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) SymbolPriceTicker() *request.Request {
	return request.New("ticker.price", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) SymbolBookTicker() *request.Request {
	return request.New("ticker.book", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) Ping() *request.Request {
	return request.New("ping", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) Time() *request.Request {
	return request.New("time", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) ExchangeInfo() *request.Request {
	return request.New("exchangeInfo", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) UserDataStreamStart() *request.Request {
	return request.New("userDataStream.start", wa.waHost, wa.waPath, wa.sign)
}

// UserDataStreamPing(listenKey string) (newListenKey string, err error)
func (wa *WebApiWrapper) UserDataStreamPing() *request.Request {
	return request.New("userDataStream.ping", wa.waHost, wa.waPath, wa.sign)
}

// UserDataStreamStop(listenKey string) (err error)
func (wa *WebApiWrapper) UserDataStreamStop() *request.Request {
	return request.New("userDataStream.stop", wa.waHost, wa.waPath, wa.sign)
}
