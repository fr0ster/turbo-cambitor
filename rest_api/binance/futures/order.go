package futures_rest_api

import (
	"net/http"

	request "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/request"
)

func (ra *RestApi) NewOrder() *request.Request {
	return request.New(ra.apiKey, ra.symbol, http.MethodPost, ra.apiBaseUrl, "/fapi/v1/Request", ra.parameters, ra.sign)
}

func (ra *RestApi) QueryOrder() *request.Request {
	return request.New(ra.apiKey, ra.symbol, http.MethodGet, ra.apiBaseUrl, "/fapi/v1/Request", ra.parameters, ra.sign)
}

func (ra *RestApi) CancelOrder() *request.Request {
	return request.New(ra.apiKey, ra.symbol, http.MethodDelete, ra.apiBaseUrl, "/fapi/v1/Request", ra.parameters, ra.sign)
}

func (ra *RestApi) CancelAllOrder() *request.Request {
	return request.New(ra.apiKey, ra.symbol, http.MethodDelete, ra.apiBaseUrl, "/fapi/v1/openOrders", ra.parameters, ra.sign)
}

func (ra *RestApi) CancelReplaceOrder() *request.Request {
	return request.New(ra.apiKey, ra.symbol, http.MethodPost, ra.apiBaseUrl, "/fapi/v1/Request/cancelReplace", ra.parameters, ra.sign)
}

func (ra *RestApi) QueryOpenOrders() *request.Request {
	return request.New(ra.apiKey, ra.symbol, http.MethodGet, ra.apiBaseUrl, "/fapi/v1/openOrders", ra.parameters, ra.sign)
}

func (ra *RestApi) QueryAllOrders() *request.Request {
	return request.New(ra.apiKey, ra.symbol, http.MethodGet, ra.apiBaseUrl, "/fapi/v1/allOrders", ra.parameters, ra.sign)
}
