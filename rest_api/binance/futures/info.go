package futures_rest_api

import (
	"net/http"

	rest_api "github.com/fr0ster/turbo-restler/rest_api"
)

func (ra *RestApi) Lock() {
	ra.mutex.Lock()
}

func (ra *RestApi) Unlock() {
	ra.mutex.Unlock()
}

func (ra *RestApi) Ping() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/ping", ra.sign)
	return
}

func (ra *RestApi) ExchangeInfo() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/exchangeInfo", ra.sign)
	return
}

func (ra *RestApi) OrderBook() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/depth", ra.sign)
	return
}

func (ra *RestApi) RecentTrades() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/trades", ra.sign)
	return
}

func (ra *RestApi) OldTradesLookup() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/historicalTrades", ra.sign)
	return
}

func (ra *RestApi) AggregateTrades() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/aggTrades", ra.sign)
	return
}

func (ra *RestApi) Klines() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/klines", ra.sign)
	return
}

func (ra *RestApi) ContinuousKlines() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/continuousKlines", ra.sign)
	return
}

func (ra *RestApi) IndexPriceKlines() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/indexPriceKlines", ra.sign)
	return
}

func (ra *RestApi) MarkPriceKlines() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/markPriceKlines", ra.sign)
	return
}

func (ra *RestApi) PremiumIndexKlines() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/premiumIndexKlines", ra.sign)
	return
}

func (ra *RestApi) MarkPrice() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/markPrice", ra.sign)
	return
}

func (ra *RestApi) FundingRate() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/fundingRate", ra.sign)
	return
}

func (ra *RestApi) FundingInfo() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/premiumIndex", ra.sign)
	return
}

func (ra *RestApi) Ticker24hr() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/ticker/24hr", ra.sign)
	return
}

func (ra *RestApi) TickerPrice() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/ticker/price", ra.sign)
	return
}

func (ra *RestApi) TickerPriceV2() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v2/ticker/price", ra.sign)
	return
}

func (ra *RestApi) BookTicker() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/ticker/bookTicker", ra.sign)
	return
}

func (ra *RestApi) QuarterlyContractSettlementPrice() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/futures/data/delivery-price", ra.sign)
	return
}

func (ra *RestApi) OpenInterest() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/futures/data/openInterest", ra.sign)
	return
}

func (ra *RestApi) OpenInterestHist() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/futures/data/openInterestHist", ra.sign)
	return
}

func (ra *RestApi) TopLongShortPositionRatio() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/futures/data/topLongShortAccountRatio", ra.sign)
	return
}

func (ra *RestApi) TopLongShortAccountRatio() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/futures/data/topLongShortPositionRatio", ra.sign)
	return
}

func (ra *RestApi) GlobalLongShortAccountRatio() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/futures/data/globalLongShortAccountRatio", ra.sign)
	return
}

func (ra *RestApi) TakerBuySellVolume() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/futures/data/takerlongshortRatio", ra.sign)
	return
}

func (ra *RestApi) HistoricalBLVTNavCandlestick() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "//fapi/v1/lvtKlines", ra.sign)
	return
}

func (ra *RestApi) CompositeIndexSymbol() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/indexInfo", ra.sign)
	return
}

func (ra *RestApi) MultiAssetsModeAssetIndex() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/assetIndex", ra.sign)
	return
}
