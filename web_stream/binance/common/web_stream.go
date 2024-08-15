package common_web_stream

import (
	"strconv"
	"strings"
	"sync"

	"github.com/bitly/go-simplejson"
	stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream"

	web_api "github.com/fr0ster/turbo-restler/web_api"
)

func (wa *WebStream) Klines(interval string) *stream.Stream {
	wsPath := web_api.WsPath("/" + strings.ToLower(wa.symbol) + "@kline_" + interval)
	return stream.New(wa.waHost, wsPath)
}

func (wa *WebStream) ContinuousKlines(interval string, contractType string) *stream.Stream {
	wsPath := web_api.WsPath("/" + strings.ToLower(wa.symbol) + strings.ToLower(contractType) + "@continuousKline_" + interval)
	return stream.New(wa.waHost, wsPath)
}

func (wa *WebStream) PartialBookDepths(level DepthStreamLevel, rates ...DepthStreamRate) *stream.Stream {
	var wsPath web_api.WsPath
	if len(rates) == 0 {
		wsPath = web_api.WsPath("/" + strings.ToLower(wa.symbol) + "@depth" + strconv.Itoa(int(level)) + "@" + strconv.Itoa(int(rates[0])) + "ms")
	} else {
		wsPath = web_api.WsPath("/" + strings.ToLower(wa.symbol) + "@depth" + strconv.Itoa(int(level)))
	}
	return stream.New(wa.waHost, wsPath)
}

func (wa *WebStream) DiffBookDepths(rates ...DepthStreamRate) *stream.Stream {
	var wsPath web_api.WsPath
	if len(rates) == 0 {
		wsPath = web_api.WsPath("/" + strings.ToLower(wa.symbol) + "@depth" + "@" + strconv.Itoa(int(rates[0])) + "ms")
	} else {
		wsPath = web_api.WsPath("/" + strings.ToLower(wa.symbol) + "@depth")
	}
	return stream.New(wa.waHost, wsPath)
}

func (wa *WebStream) AggTrades() *stream.Stream {
	wsPath := web_api.WsPath("/" + strings.ToLower(wa.symbol) + "@aggTrade")
	return stream.New(wa.waHost, wsPath)
}

func (wa *WebStream) Trades() *stream.Stream {
	wsPath := web_api.WsPath("/" + strings.ToLower(wa.symbol) + "@Trade")
	return stream.New(wa.waHost, wsPath)
}

func (wa *WebStream) BookTickers() *stream.Stream {
	wsPath := web_api.WsPath("/" + strings.ToLower(wa.symbol) + "@bookTicker")
	return stream.New(wa.waHost, wsPath)
}

func (wa *WebStream) Tickers() *stream.Stream {
	wsPath := web_api.WsPath("/" + strings.ToLower(wa.symbol) + "@ticker")
	return stream.New(wa.waHost, wsPath)
}

func (wa *WebStream) MiniTickers() *stream.Stream {
	wsPath := web_api.WsPath("/" + strings.ToLower(wa.symbol) + "@miniTicker")
	return stream.New(wa.waHost, wsPath)
}

func (wa *WebStream) UserData(listenKey string) *stream.Stream {
	return stream.New(wa.waHost, web_api.WsPath("/"+listenKey))
}

func (wa *WebStream) MarkPrice() *stream.Stream {
	wsPath := web_api.WsPath("/" + strings.ToLower(wa.symbol) + "@markPrice")
	return stream.New(wa.waHost, wsPath)
}

func (wa *WebStream) LiquidationOrder() *stream.Stream {
	wsPath := web_api.WsPath("/" + strings.ToLower(wa.symbol) + "@forceOrder")
	return stream.New(wa.waHost, wsPath)
}

func (wa *WebStream) ContractInfo() *stream.Stream {
	wsPath := web_api.WsPath("/" + strings.ToLower(wa.symbol) + "!contractInfo")
	return stream.New(wa.waHost, wsPath)
}

func (wa *WebStream) Subscribe(streams []string) (response *simplejson.Json, err error) {
	wsPath := web_api.WsPath("")
	socket, err := web_api.New(wa.waHost, wsPath)
	if err != nil {
		return nil, err
	}
	rq := simplejson.New()
	rq.Set("method", "SUBSCRIBE")
	rq.Set("id", 1)
	params := make([]string, 0)
	for _, stream := range streams {
		params = append(params, strings.ToLower(wa.symbol)+"@"+stream)
	}
	rq.Set("params", params)
	return socket.Call(rq)
}

func (wa *WebStream) ListOfSubscriptions() (response *simplejson.Json, err error) {
	wsPath := web_api.WsPath("")
	socket, err := web_api.New(wa.waHost, wsPath)
	if err != nil {
		return nil, err
	}
	rq := simplejson.New()
	rq.Set("method", "LIST_SUBSCRIPTIONS")
	rq.Set("id", 1)
	return socket.Call(rq)
}

func (wa *WebStream) Unsubscribe(streams []string) (response *simplejson.Json, err error) {
	wsPath := web_api.WsPath("")
	socket, err := web_api.New(wa.waHost, wsPath)
	if err != nil {
		return nil, err
	}
	rq := simplejson.New()
	rq.Set("method", "SUBSCRIBE")
	rq.Set("id", 1)
	params := make([]string, 0)
	for _, stream := range streams {
		params = append(params, strings.ToLower(wa.symbol)+"@"+stream)
	}
	rq.Set("params", params)
	return socket.Call(rq)
}

func New(host web_api.WsHost, symbol string) *WebStream {
	return &WebStream{
		symbol: symbol,
		waHost: host,
		mutex:  &sync.Mutex{},
	}
}
