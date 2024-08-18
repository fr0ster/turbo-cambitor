package request

import (
	"time"

	"github.com/bitly/go-simplejson"

	"github.com/fr0ster/turbo-restler/web_socket"
	signature "github.com/fr0ster/turbo-signer/signature"
	"github.com/google/uuid"
)

type (
	Method  string
	Request struct {
		sign   signature.Sign
		waHost web_socket.WsHost
		waPath web_socket.WsPath
		method Method
		params *simplejson.Json
		// connection *web_api.WebApi
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

func (rq *Request) Do() (result *simplejson.Json) {
	result = simplejson.New()
	result.Set("id", uuid.New().String())
	result.Set("method", rq.method)
	if rq.params != nil {
		result.Set("params", rq.params)
	}
	return
}

func New(method Method, waHost web_socket.WsHost, waPath web_socket.WsPath, sign signature.Sign) *Request {
	return &Request{
		sign:   sign,
		waHost: waHost,
		waPath: waPath,
		params: nil,
		method: method,
	}
}
