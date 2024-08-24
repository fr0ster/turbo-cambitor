package futures_rest_api

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

func (ra *RestApiWrapper) Ping() *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/ping", nil, ra.GetSignature())
	return rq
}

func (ra *RestApiWrapper) ExchangeInfo() *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/exchangeInfo", nil, ra.GetSignature())
	return rq
}

func (ra *RestApiWrapper) OrderBook(symbol string) *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/depth", nil, ra.GetSignature())
	return rq.Set("symbol", symbol)
}

func (ra *RestApiWrapper) RecentTrades(symbol string) *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/trades", nil, ra.GetSignature())
	return rq.Set("symbol", symbol)
}

func (ra *RestApiWrapper) OldTradesLookup(symbol string) *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/historicalTrades", nil, ra.GetSignature())
	return rq.Set("symbol", symbol)
}

func (ra *RestApiWrapper) AggregateTrades(symbol string) *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/aggTrades", nil, ra.GetSignature())
	return rq.Set("symbol", symbol)
}

func (ra *RestApiWrapper) Klines(symbol string, internal string) *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/klines", nil, ra.GetSignature())
	return rq.Set("symbol", symbol).Set("interval", internal)
}

func (ra *RestApiWrapper) ContinuousKlines(symbol string, interval string, contractType string) *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/continuousKlines", nil, ra.GetSignature())
	return rq.Set("pair", symbol).Set("interval", interval).Set("contractType", contractType)
}

func (ra *RestApiWrapper) IndexPriceKlines(symbol string, interval string) *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/indexPriceKlines", nil, ra.GetSignature())
	return rq.Set("pair", symbol).Set("interval", interval)
}

func (ra *RestApiWrapper) MarkPriceKlines(symbol string, interval string) *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/markPriceKlines", nil, ra.GetSignature())
	return rq.Set("symbol", symbol).Set("interval", interval)
}

func (ra *RestApiWrapper) FundingRate() *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/fundingRate", nil, ra.GetSignature())
	return rq
}

func (ra *RestApiWrapper) FundingInfo() *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/premiumIndex", nil, ra.GetSignature())
	return rq
}

func (ra *RestApiWrapper) Ticker24hr() *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/ticker/24hr", nil, ra.GetSignature())
	return rq
}

func (ra *RestApiWrapper) TickerPrice() *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/ticker/price", nil, ra.GetSignature())
	return rq
}

func (ra *RestApiWrapper) TickerPriceV2() *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v2/ticker/price", nil, ra.GetSignature())
	return rq
}

func (ra *RestApiWrapper) BookTicker() *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/ticker/bookTicker", nil, ra.GetSignature())
	return rq
}

func (ra *RestApiWrapper) OpenInterest(symbol string) *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/openInterest", nil, ra.GetSignature())
	return rq.Set("symbol", symbol)
}

func (ra *RestApiWrapper) CompositeIndexSymbol() *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/indexInfo", nil, ra.GetSignature())
	return rq
}

func (ra *RestApiWrapper) MultiAssetsModeAssetIndex() *request.RequestBuilder {
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/assetIndex", nil, ra.GetSignature())
	return rq
}
