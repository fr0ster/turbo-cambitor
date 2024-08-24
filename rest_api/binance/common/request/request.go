package order

import (
	"net/http"
	"time"

	"github.com/bitly/go-simplejson"

	common_rest_api "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/rest_api"
	rest_api "github.com/fr0ster/turbo-restler/rest_api"
	signature "github.com/fr0ster/turbo-signer/signature"
)

type (
	RequestBuilder struct {
		sign       signature.Sign
		apiBaseUrl rest_api.ApiBaseUrl
		endPoint   rest_api.EndPoint
		method     rest_api.HttpMethod
		params     *simplejson.Json
	}
)

func (rq *RequestBuilder) Set(name string, value interface{}) *RequestBuilder {
	if rq.params == nil {
		rq.params = simplejson.New()
	}
	rq.params.Set(name, value)
	return rq
}

func (rq *RequestBuilder) SetAPIKey() *RequestBuilder {
	if rq.params == nil {
		rq.params = simplejson.New()
	}
	rq.params.Set("apiKey", rq.sign.GetAPIKey())
	return rq
}

func (rq *RequestBuilder) SetTimestamp() *RequestBuilder {
	if rq.params == nil {
		rq.params = simplejson.New()
	}
	rq.params.Set("timestamp", int64(time.Nanosecond)*time.Now().UnixNano()/int64(time.Millisecond))
	return rq
}

func (rq *RequestBuilder) SetSignature() *RequestBuilder {
	if rq.params == nil {
		rq.params = simplejson.New()
	}
	rq.params, _ = rq.sign.SignParameters(rq.params)
	return rq
}

func (rq *RequestBuilder) Do() (req *http.Request, err error) {
	req, err = common_rest_api.Params2Request(
		rq.params,
		rq.apiBaseUrl,
		rq.endPoint,
		rq.method,
		rq.sign.GetAPIKey())
	return
}

func New(method rest_api.HttpMethod, baseUrl rest_api.ApiBaseUrl, endPoint rest_api.EndPoint, params *simplejson.Json, sign signature.Sign) *RequestBuilder {
	return &RequestBuilder{
		sign:       sign,
		apiBaseUrl: baseUrl,
		endPoint:   endPoint,
		method:     method,
		params:     params,
	}
}
