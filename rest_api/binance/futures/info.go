package futures_rest_api

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

func (ra *RestApi) Ping() *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/ping", nil, ra.sign)
	return rq
}

func (ra *RestApi) ExchangeInfo() *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/exchangeInfo", nil, ra.sign)
	return rq
}

func (ra *RestApi) OrderBook(symbol string) *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/depth", nil, ra.sign)
	return rq.Set("symbol", symbol)
}

func (ra *RestApi) RecentTrades(symbol string) *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/trades", nil, ra.sign)
	return rq.Set("symbol", symbol)
}

func (ra *RestApi) OldTradesLookup(symbol string) *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/historicalTrades", nil, ra.sign)
	return rq.Set("symbol", symbol)
}

func (ra *RestApi) AggregateTrades(symbol string) *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/aggTrades", nil, ra.sign)
	return rq.Set("symbol", symbol)
}

func (ra *RestApi) Klines(symbol string, internal string) *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/klines", nil, ra.sign)
	return rq.Set("symbol", symbol).Set("interval", internal)
}

func (ra *RestApi) ContinuousKlines(symbol string, interval string, contractType string) *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/continuousKlines", nil, ra.sign)
	return rq.Set("pair", symbol).Set("interval", interval).Set("contractType", contractType)
}

func (ra *RestApi) IndexPriceKlines(symbol string, interval string) *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/indexPriceKlines", nil, ra.sign)
	return rq.Set("pair", symbol).Set("interval", interval)
}

func (ra *RestApi) MarkPriceKlines(symbol string, interval string) *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/markPriceKlines", nil, ra.sign)
	return rq.Set("symbol", symbol).Set("interval", interval)
}

func (ra *RestApi) FundingRate() *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/fundingRate", nil, ra.sign)
	return rq
}

func (ra *RestApi) FundingInfo() *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/premiumIndex", nil, ra.sign)
	return rq
}

func (ra *RestApi) Ticker24hr() *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/ticker/24hr", nil, ra.sign)
	return rq
}

func (ra *RestApi) TickerPrice() *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/ticker/price", nil, ra.sign)
	return rq
}

func (ra *RestApi) TickerPriceV2() *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v2/ticker/price", nil, ra.sign)
	return rq
}

func (ra *RestApi) BookTicker() *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/ticker/bookTicker", nil, ra.sign)
	return rq
}

func (ra *RestApi) OpenInterest(symbol string) *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/openInterest", nil, ra.sign)
	return rq.Set("symbol", symbol)
}

func (ra *RestApi) CompositeIndexSymbol() *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/indexInfo", nil, ra.sign)
	return rq
}

func (ra *RestApi) MultiAssetsModeAssetIndex() *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/assetIndex", nil, ra.sign)
	return rq
}
