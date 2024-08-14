package request

import (
	"fmt"
	"time"

	"github.com/bitly/go-simplejson"

	web_api "github.com/fr0ster/turbo-restler/web_api"
	signature "github.com/fr0ster/turbo-signer/signature"
	"github.com/google/uuid"
)

type (
	Method  string
	Request struct {
		sign       signature.Sign
		waHost     web_api.WsHost
		waPath     web_api.WsPath
		method     Method
		params     *simplejson.Json
		connection *web_api.WebApi
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
	rq.params, _ = signature.SignParameters(rq.params, rq.sign)
	return rq
}

func (rq *Request) Do() (result *simplejson.Json, err error) {
	request := simplejson.New()
	request.Set("id", uuid.New().String())
	request.Set("method", rq.method)
	if rq.params != nil {
		request.Set("params", rq.params)
	}
	response, err := rq.connection.Call(request)
	if err != nil {
		return
	}

	if response.Get("status").MustInt() != 200 {
		err = fmt.Errorf("error request: %v", response.Get("error").Interface())
		return
	}

	bytes, err := response.Get("result").MarshalJSON()
	if err != nil {
		return
	}
	result, err = simplejson.NewJson(bytes)
	return
}

func New(method Method, waHost web_api.WsHost, waPath web_api.WsPath, sign signature.Sign) *Request {
	wa, err := web_api.New(waHost, waPath)
	if err != nil {
		return nil
	}
	return &Request{
		sign:       sign,
		waHost:     waHost,
		waPath:     waPath,
		params:     nil,
		method:     method,
		connection: wa,
	}
}
