package spot_rest_api

import (
	"net/http"

	"github.com/bitly/go-simplejson"
	request "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/request"
)

func (ra *RestApi) OrderListOCO() *request.Request {
	params := simplejson.New()
	return request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/orderList/oco", params, ra.sign)
}

func (ra *RestApi) OrderListOTO() *request.Request {
	params := simplejson.New()
	return request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/orderList/oto", params, ra.sign)
}

func (ra *RestApi) NewOrderOCO() *request.Request {
	params := simplejson.New()
	return request.New(http.MethodPost, ra.apiBaseUrl, "/api/v3/order/oco", params, ra.sign)
}

func (ra *RestApi) CancelOrderOCO() *request.Request {
	params := simplejson.New()
	return request.New(http.MethodDelete, ra.apiBaseUrl, "/api/v3/orderList/oco", params, ra.sign)
}

func (ra *RestApi) QueryOrderOCO() *request.Request {
	params := simplejson.New()
	return request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/orderList ", params, ra.sign)
}

func (ra *RestApi) QueryAllOrderOCO() *request.Request {
	params := simplejson.New()
	return request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/allOrderList", params, ra.sign)
}

func (ra *RestApi) QueryOpenOrderOCO() *request.Request {
	params := simplejson.New()
	return request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/openOrderList", params, ra.sign)
}

func (ra *RestApi) NewOrderSOR() *request.Request {
	params := simplejson.New()
	return request.New(http.MethodPost, ra.apiBaseUrl, "/api/v3/order/sor", params, ra.sign)
}
