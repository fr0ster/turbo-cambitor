package order

import (
	"fmt"

	"github.com/bitly/go-simplejson"

	signature "github.com/fr0ster/turbo-restler/utils/signature"
	web_api "github.com/fr0ster/turbo-restler/web_api"
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
	rq.params.Set(name, value)
	return rq
}

func (rq *Request) Do() (order *simplejson.Json, err error) {
	request := simplejson.New()
	request.Set("method", rq.method)
	request.Set("params", rq.params)
	rq.connection.SignParameters(request, rq.sign)
	response, err := rq.connection.Call(request)
	if err != nil {
		return
	}

	if response.Get("Status").MustInt() != 200 {
		err = fmt.Errorf("error request: %v", response.Get("Error").MustString())
		return
	}

	bytes, err := response.Get("Result").MarshalJSON()
	if err != nil {
		return
	}
	order, err = simplejson.NewJson(bytes)
	return
}

func New(apiKey, symbol string, method Method, waHost web_api.WsHost, waPath web_api.WsPath, sign signature.Sign) *Request {
	simpleJson := simplejson.New()
	simpleJson.Set("apiKey", apiKey)
	simpleJson.Set("symbol", symbol)
	wa, err := web_api.New(waHost, waPath)
	if err != nil {
		return nil
	}
	return &Request{
		sign:       sign,
		waHost:     waHost,
		waPath:     waPath,
		params:     simpleJson,
		method:     method,
		connection: wa,
	}
}
