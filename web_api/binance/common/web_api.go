package common_web_api

import (
	"fmt"

	"github.com/bitly/go-simplejson"
	request "github.com/fr0ster/turbo-cambitor/web_api/binance/common/request"
	web_api "github.com/fr0ster/turbo-restler/web_api"
	"github.com/google/uuid"
)

func (wa *WebApi) ListOfSubscriptions() *request.Request {
	return request.New("LIST_SUBSCRIPTIONS", wa.waHost, wa.waPath, wa.sign)
}

// Функція для логіну
func (wa *WebApi) Logon() *request.Request {
	return request.New("session.logon", wa.waHost, wa.waPath, wa.sign)
}

// Функція для логіну
func (wa *WebApi) Logout() *request.Request {
	return request.New("session.logout", wa.waHost, wa.waPath, wa.sign)
}

// Функція для перевірки статусу сесії
func (wa *WebApi) Status() *request.Request {
	return request.New("session.status", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) OrderBook() *request.Request {
	return request.New("depth", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) SymbolPriceTicker() *request.Request {
	return request.New("ticker.price", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) SymbolBookTicker() *request.Request {
	return request.New("ticker.book", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) Ping() *request.Request {
	return request.New("ping", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) Time() *request.Request {
	return request.New("time", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) ExchangeInfo() *request.Request {
	return request.New("exchangeInfo", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) UserDataStreamStart() (listenKey string, err error) {
	web_api, err := web_api.New(wa.waHost, wa.waPath)
	if err != nil {
		return
	}
	req := simplejson.New()
	req.Set("id", uuid.New().String())
	req.Set("method", "userDataStream.start")
	params := simplejson.New()
	params.Set("apiKey", wa.sign.GetAPIKey())
	req.Set("params", params)
	response, err := web_api.Call(req)
	if err != nil {
		return
	}
	if response.Get("status").MustInt() != 200 {
		err = fmt.Errorf("error: %s", response.Get("error").Get("msg").MustString())
		return
	}
	result := response.Get("result")
	if result != nil {
		listenKey = result.Get("listenKey").MustString()
	} else {
		return
	}
	return
}

func (wa *WebApi) UserDataStreamPing(listenKey string) (newListenKey string, err error) {
	web_api, err := web_api.New(wa.waHost, wa.waPath)
	if err != nil {
		return
	}
	req := simplejson.New()
	req.Set("id", uuid.New().String())
	req.Set("method", "userDataStream.ping")
	params := simplejson.New()
	params.Set("apiKey", wa.sign.GetAPIKey())
	params.Set("listenKey", listenKey)
	req.Set("params", params)
	response, err := web_api.Call(req)
	if err != nil {
		return
	}
	if response.Get("status").MustInt() != 200 {
		err = fmt.Errorf("error: %s", response.Get("error").Get("msg").MustString())
		return
	}
	result := response.Get("result")
	if result != nil {
		newListenKey = result.Get("listenKey").MustString()
	} else {
		return
	}
	return
}

func (wa *WebApi) UserDataStreamStop(listenKey string) (err error) {
	web_api, err := web_api.New(wa.waHost, wa.waPath)
	if err != nil {
		return
	}
	req := simplejson.New()
	req.Set("id", uuid.New().String())
	req.Set("method", "userDataStream.stop")
	params := simplejson.New()
	params.Set("apiKey", wa.sign.GetAPIKey())
	params.Set("listenKey", listenKey)
	req.Set("params", params)
	response, err := web_api.Call(req)
	if err != nil {
		return
	}
	if response.Get("status").MustInt() != 200 {
		err = fmt.Errorf("error: %s", response.Get("error").Get("msg").MustString())
		return
	}
	return
}
