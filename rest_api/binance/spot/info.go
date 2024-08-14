package spot_rest_api

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

func (ra *RestApi) ExchangeInfo() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/exchangeInfo", ra.sign)
	return
}

func (ra *RestApi) Ping() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/ping", ra.sign)
	return
}

func (ra *RestApi) OrderBook() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/depth", ra.sign)
	return
}

func (ra *RestApi) RecentTrades() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/trades", ra.sign)
	return
}

func (ra *RestApi) OldTradesLookup() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/historicalTrades", ra.sign)
	return
}

func (ra *RestApi) AggregateTrades() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/aggTrades", ra.sign)
	return
}

func (ra *RestApi) Klines() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/klines", ra.sign)
	return
}

func (ra *RestApi) UIKlines() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/uiKlines", ra.sign)
	return
}

func (ra *RestApi) CurrentAvgPrice() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/avgPrice", ra.sign)
	return
}

func (ra *RestApi) Ticker24hr() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/ticker/24hr", ra.sign)
	return
}

func (ra *RestApi) TradingDay() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/ticker/tradingDay", ra.sign)
	return
}

func (ra *RestApi) TickerPrice() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/ticker/price", ra.sign)
	return
}

func (ra *RestApi) TickerBookTicker() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/ticker/bookTicker", ra.sign)
	return
}

func (ra *RestApi) TickerV3() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/ticker", ra.sign)
	return
}
