package request

import (
	"time"

	"github.com/bitly/go-simplejson"

	"github.com/fr0ster/turbo-restler/web_socket"
	signature "github.com/fr0ster/turbo-signer/signature"
	"github.com/google/uuid"
)

type (
	Method         string
	RequestBuilder struct {
		sign   signature.Sign
		waHost web_socket.WsHost
		waPath web_socket.WsPath
		method Method
		params *simplejson.Json
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

func (rq *RequestBuilder) Do() (result *simplejson.Json) {
	result = simplejson.New()
	result.Set("id", uuid.New().String())
	result.Set("method", rq.method)
	if rq.params != nil {
		result.Set("params", rq.params)
	}
	return
}

func New(method Method, waHost web_socket.WsHost, waPath web_socket.WsPath, sign signature.Sign) *RequestBuilder {
	return &RequestBuilder{
		sign:   sign,
		waHost: waHost,
		waPath: waPath,
		params: nil,
		method: method,
	}
}
