package common_web_stream

import (
	"strconv"
	"strings"
	"sync"

	stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream"

	web_api "github.com/fr0ster/turbo-restler/web_api"
)

func (wa *WebStream) Klines(interval string) *stream.Stream {
	// wsPath := web_api.WsPath(fmt.Sprintf("%s@kline_%s", strings.ToLower(wa.symbol), interval))
	wsPath := web_api.WsPath("/" + strings.ToLower(wa.symbol) + "@kline_" + interval)
	return stream.New(wa.waHost, wsPath)
}

func (wa *WebStream) Depths(level DepthStreamLevel) *stream.Stream {
	wsPath := web_api.WsPath("/" + strings.ToLower(wa.symbol) + "@depth" + strconv.Itoa(int(level)))
	return stream.New(wa.waHost, wsPath)
}

func (wa *WebStream) Depths100ms(level DepthStreamLevel) *stream.Stream {
	wsPath := web_api.WsPath("/" + strings.ToLower(wa.symbol) + "@depth%v@100ms" + strconv.Itoa(int(level)))
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

func New(host web_api.WsHost, symbol string) *WebStream {
	return &WebStream{
		symbol: symbol,
		waHost: host,
		mutex:  &sync.Mutex{},
	}
}
