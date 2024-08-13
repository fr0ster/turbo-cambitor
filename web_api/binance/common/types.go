package common_web_api

import (
	"sync"

	"github.com/bitly/go-simplejson"

	signature "github.com/fr0ster/turbo-restler/utils/signature"
	web_api "github.com/fr0ster/turbo-restler/web_api"
)

const (
	DepthAPILimit5    DepthAPILimit = 5
	DepthAPILimit10   DepthAPILimit = 10
	DepthAPILimit20   DepthAPILimit = 20
	DepthAPILimit50   DepthAPILimit = 50
	DepthAPILimit100  DepthAPILimit = 100
	DepthAPILimit500  DepthAPILimit = 500
	DepthAPILimit1000 DepthAPILimit = 1000
)

type (
	DepthAPILimit int
	WebApi        struct {
		apiKey     string
		apiSecret  string
		symbol     string
		waHost     web_api.WsHost
		waPath     web_api.WsPath
		mutex      *sync.Mutex
		sign       signature.Sign
		connection *web_api.WebApi
	}
	Result struct {
		APIKey           string `json:"apiKey"`
		AuthorizedSince  int64  `json:"authorizedSince"`
		ConnectedSince   int64  `json:"connectedSince"`
		ReturnRateLimits bool   `json:"returnRateLimits"`
		ServerTime       int64  `json:"serverTime"`
	}
	Method  string
	Request struct {
		ID     string           `json:"id"`
		Method Method           `json:"method"`
		Params *simplejson.Json `json:"params"`
	}
	Response struct {
		ID         string      `json:"id"`
		Status     int         `json:"status"`
		Error      ErrorDetail `json:"error"`
		Result     interface{} `json:"result"`
		RateLimits []RateLimit `json:"rateLimits"`
	}

	ErrorDetail struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	RateLimit struct {
		RateLimitType string `json:"rateLimitType"`
		Interval      string `json:"interval"`
		IntervalNum   int    `json:"intervalNum"`
		Limit         int    `json:"limit"`
		Count         int    `json:"count"`
	}
)

func (rq *Request) SetParameter(name string, value string) {
	if rq.Params == nil {
		rq.Params = simplejson.New()
	}
	rq.Params.Set(name, value)
}
