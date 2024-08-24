package spot_rest_api

import (
	"net/http"

	"github.com/bitly/go-simplejson"
	request "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/request"
)

func (ra *RestApiWrapper) NewOrder() *request.RequestBuilder {
	params := simplejson.New()
	return request.New(http.MethodPost, ra.GetApiBaseUrl(), "/api/v3/Request", params, ra.GetSignature())
}

func (ra *RestApiWrapper) TestOrder() *request.RequestBuilder {
	params := simplejson.New()
	return request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/Request/test", params, ra.GetSignature())
}

func (ra *RestApiWrapper) QueryOrder() *request.RequestBuilder {
	params := simplejson.New()
	return request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/Request", params, ra.GetSignature())
}

func (ra *RestApiWrapper) CancelOrder() *request.RequestBuilder {
	params := simplejson.New()
	return request.New(http.MethodDelete, ra.GetApiBaseUrl(), "/api/v3/Request", params, ra.GetSignature())
}

func (ra *RestApiWrapper) CancelAllOrder() *request.RequestBuilder {
	params := simplejson.New()
	return request.New(http.MethodDelete, ra.GetApiBaseUrl(), "/api/v3/openOrders", params, ra.GetSignature())
}

func (ra *RestApiWrapper) CancelReplaceOrder() *request.RequestBuilder {
	params := simplejson.New()
	return request.New(http.MethodPost, ra.GetApiBaseUrl(), "/api/v3/Request/cancelReplace", params, ra.GetSignature())
}

func (ra *RestApiWrapper) QueryOpenOrders() *request.RequestBuilder {
	params := simplejson.New()
	return request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/openOrders", params, ra.GetSignature())
}

func (ra *RestApiWrapper) QueryAllOrders() *request.RequestBuilder {
	params := simplejson.New()
	return request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/allOrders", params, ra.GetSignature())
}
