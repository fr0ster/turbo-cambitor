package common_web_stream

import (
	"strconv"
	"strings"

	stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream"
	"github.com/fr0ster/turbo-restler/web_socket"
	"github.com/gorilla/websocket"
)

type WebStream struct {
	scheme stream.WsScheme
	waHost stream.WsHost
	symbol string
	wsPath stream.WsPath
}

func (wa *WebStream) makeStream(path string) stream.StreamInterface {
	wa.wsPath = stream.WsPath(path)
	factory := func() (web_socket.WebSocketInterface, error) {
		url := string(wa.scheme) + "://" + string(wa.waHost) + string(wa.wsPath)
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			return nil, err
		}
		return web_socket.NewWebSocketWrapper(conn), nil
	}
	return stream.NewStreamWrapper(factory, wa.wsPath).SetSymbol(wa.symbol)
}

func (wa *WebStream) Klines(interval string) stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@kline_" + interval)
}

func (wa *WebStream) ContinuousKlines(interval, contractType string) stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + strings.ToLower(contractType) + "@continuousKline_" + interval)
}

func (wa *WebStream) PartialBookDepths(level DepthStreamLevel, rates ...DepthStreamRate) stream.StreamInterface {
	if len(rates) > 0 {
		return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@depth" + strconv.Itoa(int(level)) + "@" + strconv.Itoa(int(rates[0])) + "ms")
	}
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@depth" + strconv.Itoa(int(level)))
}

func (wa *WebStream) DiffBookDepths(rates ...DepthStreamRate) stream.StreamInterface {
	if len(rates) > 0 {
		return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@depth@" + strconv.Itoa(int(rates[0])) + "ms")
	}
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@depth")
}

func (wa *WebStream) AggTrades() stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@aggTrade")
}

func (wa *WebStream) Trades() stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@trade")
}

func (wa *WebStream) BookTickers() stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@bookTicker")
}

func (wa *WebStream) Tickers() stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@ticker")
}

func (wa *WebStream) MiniTickers() stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@miniTicker")
}

func (wa *WebStream) UserData(listenKey string) stream.StreamInterface {
	return wa.makeStream("/" + listenKey)
}

func (wa *WebStream) MarkPrice() stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@markPrice")
}

func (wa *WebStream) LiquidationOrder() stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@forceOrder")
}

func (wa *WebStream) ContractInfo() stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "!contractInfo")
}

func (wa *WebStream) Stream() stream.StreamInterface {
	return wa.makeStream("")
}
