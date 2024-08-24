package spot_rest_api

import (
	"net/http"

	request "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/request"
)

func (ra *RestApiWrapper) Lock() {
	ra.GetMutex().Lock()
}

func (ra *RestApiWrapper) Unlock() {
	ra.GetMutex().Unlock()
}

func (ra *RestApiWrapper) ExchangeInfo() *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/exchangeInfo", nil, ra.GetSignature())
	return rq
}

func (ra *RestApiWrapper) Ping() *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/ping", nil, ra.GetSignature())
	return rq
}

func (ra *RestApiWrapper) OrderBook(symbol string) *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/depth", nil, ra.GetSignature())
	return rq.Set("symbol", symbol)
}

func (ra *RestApiWrapper) RecentTrades(symbol string) *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/trades", nil, ra.GetSignature())
	return rq.Set("symbol", symbol)
}

func (ra *RestApiWrapper) OldTradesLookup(symbol string) *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/historicalTrades", nil, ra.GetSignature())
	return rq.Set("symbol", symbol)
}

func (ra *RestApiWrapper) AggregateTrades(symbol string) *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/aggTrades", nil, ra.GetSignature())
	return rq.Set("symbol", symbol)
}

func (ra *RestApiWrapper) Klines(symbol string, interval string) *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/klines", nil, ra.GetSignature())
	return rq.Set("symbol", symbol).Set("interval", interval)
}

func (ra *RestApiWrapper) UIKlines(symbol string, interval string) *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/uiKlines", nil, ra.GetSignature())
	return rq.Set("symbol", symbol).Set("interval", interval)
}

func (ra *RestApiWrapper) CurrentAvgPrice(symbol string) *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/avgPrice", nil, ra.GetSignature())
	return rq.Set("symbol", symbol)
}

func (ra *RestApiWrapper) Ticker24hr(symbol string) *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/ticker/24hr", nil, ra.GetSignature())
	return rq.Set("symbol", symbol)
}

func (ra *RestApiWrapper) TradingDay() *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/ticker/tradingDay", nil, ra.GetSignature())
	return rq
}

func (ra *RestApiWrapper) TickerPrice() *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/ticker/price", nil, ra.GetSignature())
	return rq
}

func (ra *RestApiWrapper) TickerBookTicker() *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/ticker/bookTicker", nil, ra.GetSignature())
	return rq
}

func (ra *RestApiWrapper) TickerV3() *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/ticker/price", nil, ra.GetSignature())
	return rq
}
