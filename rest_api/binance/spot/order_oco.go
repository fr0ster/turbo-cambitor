package spot_rest_api

import (
	"net/http"

	request "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/request"
)

func (ra *RestApi) OrderListOCO() *request.Request {
	return request.New(ra.apiKey, ra.symbol, http.MethodGet, ra.apiBaseUrl, "/api/v3/orderList/oco", ra.parameters, ra.sign)
}

func (ra *RestApi) OrderListOTO() *request.Request {
	return request.New(ra.apiKey, ra.symbol, http.MethodGet, ra.apiBaseUrl, "/api/v3/orderList/oto", ra.parameters, ra.sign)
}

func (ra *RestApi) NewOrderOCO() *request.Request {
	return request.New(ra.apiKey, ra.symbol, http.MethodPost, ra.apiBaseUrl, "/api/v3/order/oco", ra.parameters, ra.sign)
}

func (ra *RestApi) CancelOrderOCO() *request.Request {
	return request.New(ra.apiKey, ra.symbol, http.MethodDelete, ra.apiBaseUrl, "/api/v3/orderList/oco", ra.parameters, ra.sign)
}

func (ra *RestApi) QueryOrderOCO() *request.Request {
	return request.New(ra.apiKey, ra.symbol, http.MethodGet, ra.apiBaseUrl, "/api/v3/orderList ", ra.parameters, ra.sign)
}

func (ra *RestApi) QueryAllOrderOCO() *request.Request {
	return request.New(ra.apiKey, ra.symbol, http.MethodGet, ra.apiBaseUrl, "/api/v3/allOrderList", ra.parameters, ra.sign)
}

func (ra *RestApi) QueryOpenOrderOCO() *request.Request {
	return request.New(ra.apiKey, ra.symbol, http.MethodGet, ra.apiBaseUrl, "/api/v3/openOrderList", ra.parameters, ra.sign)
}

func (ra *RestApi) NewOrderSOR() *request.Request {
	return request.New(ra.apiKey, ra.symbol, http.MethodPost, ra.apiBaseUrl, "/api/v3/order/sor", ra.parameters, ra.sign)
}
