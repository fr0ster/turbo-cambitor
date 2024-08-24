package spot_rest_api

import (
	"net/http"

	"github.com/bitly/go-simplejson"
	request "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/request"
)

func (ra *RestApiWrapper) Account() *request.RequestBuilder {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/account", params, ra.GetSignature())
	return rq
}

func (ra *RestApiWrapper) QueryCurrentOrderCountUsage() *request.RequestBuilder {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/rateLimit/order", params, ra.GetSignature())
	return rq
}

func (ra *RestApiWrapper) QueryPreventedMatches(symbol string) *request.RequestBuilder {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/myPreventedMatches", params, ra.GetSignature())
	return rq.Set("symbol", symbol)
}

func (ra *RestApiWrapper) QueryAllocations(symbol string) *request.RequestBuilder {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/api/v3/myAllocations", params, ra.GetSignature())
	return rq.Set("symbol", symbol)
}
