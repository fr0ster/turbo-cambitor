package spot_rest_api

import (
	"net/http"

	"github.com/bitly/go-simplejson"
	request "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/request"
)

func (ra *RestApi) NewOrder() *request.Request {
	params := simplejson.New()
	return request.New(http.MethodPost, ra.apiBaseUrl, "/api/v3/Request", params, ra.sign)
}

func (ra *RestApi) TestOrder() *request.Request {
	params := simplejson.New()
	return request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/Request/test", params, ra.sign)
}

func (ra *RestApi) QueryOrder() *request.Request {
	params := simplejson.New()
	return request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/Request", params, ra.sign)
}

func (ra *RestApi) CancelOrder() *request.Request {
	params := simplejson.New()
	return request.New(http.MethodDelete, ra.apiBaseUrl, "/api/v3/Request", params, ra.sign)
}

func (ra *RestApi) CancelAllOrder() *request.Request {
	params := simplejson.New()
	return request.New(http.MethodDelete, ra.apiBaseUrl, "/api/v3/openOrders", params, ra.sign)
}

func (ra *RestApi) CancelReplaceOrder() *request.Request {
	params := simplejson.New()
	return request.New(http.MethodPost, ra.apiBaseUrl, "/api/v3/Request/cancelReplace", params, ra.sign)
}

func (ra *RestApi) QueryOpenOrders() *request.Request {
	params := simplejson.New()
	return request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/openOrders", params, ra.sign)
}

func (ra *RestApi) QueryAllOrders() *request.Request {
	params := simplejson.New()
	return request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/allOrders", params, ra.sign)
}
