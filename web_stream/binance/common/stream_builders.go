package common_web_stream

import (
	"strconv"
	"strings"

	stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream"

	"github.com/fr0ster/turbo-restler/web_socket"
)

func (wa *WebStream) Klines(interval string) *stream.StreamWrapper {
	var wsPath web_socket.WsPath
	if wa.symbol != "" {
		wsPath = web_socket.WsPath("/" + strings.ToLower(wa.symbol) + "@kline_" + interval)
	}
	return stream.New(wa.waHost, wsPath, web_socket.SchemeWSS)
}

func (wa *WebStream) ContinuousKlines(interval string, contractType string) *stream.StreamWrapper {
	var wsPath web_socket.WsPath
	if wa.symbol != "" {
		wsPath = web_socket.WsPath("/" + strings.ToLower(wa.symbol) + strings.ToLower(contractType) + "@continuousKline_" + interval)
	}
	return stream.New(wa.waHost, wsPath, web_socket.SchemeWSS)
}

func (wa *WebStream) PartialBookDepths(level DepthStreamLevel, rates ...DepthStreamRate) *stream.StreamWrapper {
	var wsPath web_socket.WsPath
	if wa.symbol != "" {
		if len(rates) == 0 {
			wsPath = web_socket.WsPath("/" + strings.ToLower(wa.symbol) + "@depth" + strconv.Itoa(int(level)) + "@" + strconv.Itoa(int(rates[0])) + "ms")
		} else {
			wsPath = web_socket.WsPath("/" + strings.ToLower(wa.symbol) + "@depth" + strconv.Itoa(int(level)))
		}
	}
	return stream.New(wa.waHost, wsPath, web_socket.SchemeWSS)
}

func (wa *WebStream) DiffBookDepths(rates ...DepthStreamRate) *stream.StreamWrapper {
	var wsPath web_socket.WsPath
	if wa.symbol != "" {
		if len(rates) == 0 {
			wsPath = web_socket.WsPath("/" + strings.ToLower(wa.symbol) + "@depth" + "@" + strconv.Itoa(int(rates[0])) + "ms")
		} else {
			wsPath = web_socket.WsPath("/" + strings.ToLower(wa.symbol) + "@depth")
		}
	}
	return stream.New(wa.waHost, wsPath, web_socket.SchemeWSS)
}

func (wa *WebStream) AggTrades() *stream.StreamWrapper {
	var wsPath web_socket.WsPath
	if wa.symbol != "" {
		wsPath = web_socket.WsPath("/" + strings.ToLower(wa.symbol) + "@aggTrade")
	}
	return stream.New(wa.waHost, wsPath, web_socket.SchemeWSS)
}

func (wa *WebStream) Trades() *stream.StreamWrapper {
	var wsPath web_socket.WsPath
	if wa.symbol != "" {
		wsPath = web_socket.WsPath("/" + strings.ToLower(wa.symbol) + "@Trade")
	}
	return stream.New(wa.waHost, wsPath, web_socket.SchemeWSS)
}

func (wa *WebStream) BookTickers() *stream.StreamWrapper {
	var wsPath web_socket.WsPath
	if wa.symbol != "" {
		wsPath = web_socket.WsPath("/" + strings.ToLower(wa.symbol) + "@bookTicker")
	}
	return stream.New(wa.waHost, wsPath, web_socket.SchemeWSS)
}

func (wa *WebStream) Tickers() *stream.StreamWrapper {
	var wsPath web_socket.WsPath
	if wa.symbol != "" {
		wsPath = web_socket.WsPath("/" + strings.ToLower(wa.symbol) + "@ticker")
	}
	return stream.New(wa.waHost, wsPath, web_socket.SchemeWSS)
}

func (wa *WebStream) MiniTickers() *stream.StreamWrapper {
	var wsPath web_socket.WsPath
	if wa.symbol != "" {
		wsPath = web_socket.WsPath("/" + strings.ToLower(wa.symbol) + "@miniTicker")
	}
	return stream.New(wa.waHost, wsPath, web_socket.SchemeWSS)
}

func (wa *WebStream) UserData(listenKey string) *stream.StreamWrapper {
	return stream.New(wa.waHost, web_socket.WsPath("/"+listenKey), web_socket.SchemeWSS)
}

func (wa *WebStream) MarkPrice() *stream.StreamWrapper {
	var wsPath web_socket.WsPath
	if wa.symbol != "" {
		wsPath = web_socket.WsPath("/" + strings.ToLower(wa.symbol) + "@markPrice")
	}
	return stream.New(wa.waHost, wsPath, web_socket.SchemeWSS)
}

func (wa *WebStream) LiquidationOrder() *stream.StreamWrapper {
	var wsPath web_socket.WsPath
	if wa.symbol != "" {
		wsPath = web_socket.WsPath("/" + strings.ToLower(wa.symbol) + "@forceOrder")
	}
	return stream.New(wa.waHost, wsPath, web_socket.SchemeWSS)
}

func (wa *WebStream) ContractInfo() *stream.StreamWrapper {
	var wsPath web_socket.WsPath
	if wa.symbol != "" {
		wsPath = web_socket.WsPath("/" + strings.ToLower(wa.symbol) + "!contractInfo")
	}
	return stream.New(wa.waHost, wsPath, web_socket.SchemeWSS)
}

func (wa *WebStream) Stream() *stream.StreamWrapper {
	return stream.New(wa.waHost, "", web_socket.SchemeWSS)
}
