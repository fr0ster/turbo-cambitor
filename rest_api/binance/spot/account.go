package spot_rest_api

import (
	"net/http"

	"github.com/bitly/go-simplejson"
	request "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/request"
)

func (ra *RestApi) Account() *request.Request {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/account", params, ra.sign)
	return rq
}

func (ra *RestApi) QueryCurrentOrderCountUsage() *request.Request {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/rateLimit/order", params, ra.sign)
	return rq
}

func (ra *RestApi) QueryPreventedMatches(symbol string) *request.Request {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/myPreventedMatches", params, ra.sign)
	return rq.Set("symbol", symbol)
}

func (ra *RestApi) QueryAllocations(symbol string) *request.Request {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/api/v3/myAllocations", params, ra.sign)
	return rq.Set("symbol", symbol)
}
