package spot_rest_api

import (
	"net/http"

	"github.com/bitly/go-simplejson"
	request "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/request"
)

func (ra *RestApiWrapper) OrderListOCO() *request.RequestBuilder {
	params := simplejson.New()
	return request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/orderList/oco", params, ra.GetSignature())
}

func (ra *RestApiWrapper) OrderListOTO() *request.RequestBuilder {
	params := simplejson.New()
	return request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/orderList/oto", params, ra.GetSignature())
}

func (ra *RestApiWrapper) NewOrderOCO() *request.RequestBuilder {
	params := simplejson.New()
	return request.New(http.MethodPost, ra.GetApiBaseUrl(), "/api/v3/order/oco", params, ra.GetSignature())
}

func (ra *RestApiWrapper) CancelOrderOCO() *request.RequestBuilder {
	params := simplejson.New()
	return request.New(http.MethodDelete, ra.GetApiBaseUrl(), "/api/v3/orderList/oco", params, ra.GetSignature())
}

func (ra *RestApiWrapper) QueryOrderOCO() *request.RequestBuilder {
	params := simplejson.New()
	return request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/orderList ", params, ra.GetSignature())
}

func (ra *RestApiWrapper) QueryAllOrderOCO() *request.RequestBuilder {
	params := simplejson.New()
	return request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/allOrderList", params, ra.GetSignature())
}

func (ra *RestApiWrapper) QueryOpenOrderOCO() *request.RequestBuilder {
	params := simplejson.New()
	return request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/openOrderList", params, ra.GetSignature())
}

func (ra *RestApiWrapper) NewOrderSOR() *request.RequestBuilder {
	params := simplejson.New()
	return request.New(http.MethodPost, ra.GetApiBaseUrl(), "/api/v3/order/sor", params, ra.GetSignature())
}
