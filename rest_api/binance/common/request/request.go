package order

import (
	"time"

	"github.com/bitly/go-simplejson"

	common_rest_api "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/rest_api"
	rest_api "github.com/fr0ster/turbo-restler/rest_api"
	signature "github.com/fr0ster/turbo-signer/signature"
)

type (
	Request struct {
		sign       signature.Sign
		apiBaseUrl rest_api.ApiBaseUrl
		endPoint   rest_api.EndPoint
		method     rest_api.HttpMethod
		params     *simplejson.Json
	}
)

func (rq *Request) Set(name string, value interface{}) *Request {
	if rq.params == nil {
		rq.params = simplejson.New()
	}
	rq.params.Set(name, value)
	return rq
}

func (rq *Request) SetAPIKey() *Request {
	if rq.params == nil {
		rq.params = simplejson.New()
	}
	rq.params.Set("apiKey", rq.sign.GetAPIKey())
	return rq
}

func (rq *Request) SetTimestamp() *Request {
	if rq.params == nil {
		rq.params = simplejson.New()
	}
	rq.params.Set("timestamp", int64(time.Nanosecond)*time.Now().UnixNano()/int64(time.Millisecond))
	return rq
}

func (rq *Request) SetSignature() *Request {
	if rq.params == nil {
		rq.params = simplejson.New()
	}
	rq.params, _ = rq.sign.SignParameters(rq.params)
	return rq
}

func (rq *Request) Do() (response *simplejson.Json, err error) {
	req, err := common_rest_api.Params2Request(
		rq.params,
		rq.apiBaseUrl,
		rq.endPoint,
		rq.method,
		rq.sign.GetAPIKey())
	if err != nil {
		return
	}
	response, err = rest_api.CallRestAPI(req)
	return
}

func New(method rest_api.HttpMethod, baseUrl rest_api.ApiBaseUrl, endPoint rest_api.EndPoint, params *simplejson.Json, sign signature.Sign) *Request {
	return &Request{
		sign:       sign,
		apiBaseUrl: baseUrl,
		endPoint:   endPoint,
		method:     method,
		params:     params,
	}
}
