package common_web_api

import (
	request "github.com/fr0ster/turbo-cambitor/web_api/binance/common/request"
)

func (wa *WebApi) ListOfSubscriptions() *request.Request {
	return request.New("LIST_SUBSCRIPTIONS", wa.waHost, wa.waPath, wa.sign)
}

// Функція для логіну
func (wa *WebApi) Logon() *request.Request {
	return request.New("session.logon", wa.waHost, wa.waPath, wa.sign)
}

// Функція для логіну
func (wa *WebApi) Logout() *request.Request {
	return request.New("session.logout", wa.waHost, wa.waPath, wa.sign)
}

// Функція для перевірки статусу сесії
func (wa *WebApi) Status() *request.Request {
	return request.New("session.status", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) OrderBook() *request.Request {
	return request.New("depth", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) SymbolPriceTicker() *request.Request {
	return request.New("ticker.price", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) SymbolBookTicker() *request.Request {
	return request.New("ticker.book", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) Ping() *request.Request {
	return request.New("ping", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) Time() *request.Request {
	return request.New("time", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) ExchangeInfo() *request.Request {
	return request.New("exchangeInfo", wa.waHost, wa.waPath, wa.sign)
}
