package spot_rest_api

import (
	"net/http"

	request "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/request"
)

func (ra *RestApi) Lock() {
	ra.mutex.Lock()
}

func (ra *RestApi) Unlock() {
	ra.mutex.Unlock()
}

func (ra *RestApi) ExchangeInfo() *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/exchangeInfo", nil, ra.sign)
	return rq
}

func (ra *RestApi) Ping() *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/ping", nil, ra.sign)
	return rq
}

func (ra *RestApi) OrderBook(symbol string) *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/depth", nil, ra.sign)
	return rq.Set("symbol", symbol)
}

func (ra *RestApi) RecentTrades(symbol string) *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/trades", nil, ra.sign)
	return rq.Set("symbol", symbol)
}

func (ra *RestApi) OldTradesLookup(symbol string) *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/historicalTrades", nil, ra.sign)
	return rq.Set("symbol", symbol)
}

func (ra *RestApi) AggregateTrades(symbol string) *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/aggTrades", nil, ra.sign)
	return rq.Set("symbol", symbol)
}

func (ra *RestApi) Klines(symbol string, interval string) *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/klines", nil, ra.sign)
	return rq.Set("symbol", symbol).Set("interval", interval)
}

func (ra *RestApi) UIKlines(symbol string, interval string) *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/uiKlines", nil, ra.sign)
	return rq.Set("symbol", symbol).Set("interval", interval)
}

func (ra *RestApi) CurrentAvgPrice(symbol string) *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/avgPrice", nil, ra.sign)
	return rq.Set("symbol", symbol)
}

func (ra *RestApi) Ticker24hr(symbol string) *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/ticker/24hr", nil, ra.sign)
	return rq.Set("symbol", symbol)
}

func (ra *RestApi) TradingDay() *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/ticker/tradingDay", nil, ra.sign)
	return rq
}

func (ra *RestApi) TickerPrice() *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/ticker/price", nil, ra.sign)
	return rq
}

func (ra *RestApi) TickerBookTicker() *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/ticker/bookTicker", nil, ra.sign)
	return rq
}

func (ra *RestApi) TickerV3() *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/ticker/price", nil, ra.sign)
	return rq
}
