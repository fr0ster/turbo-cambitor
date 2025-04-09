package common_web_stream

import (
	"strconv"
	"strings"

	stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream"
	"github.com/fr0ster/turbo-restler/web_socket"
	"github.com/gorilla/websocket"
)

// NewWebStream створює новий екземпляр StreamBuilder з переданими параметрами.
func NewStreamBuilder(scheme stream.WsScheme, host stream.WsHost, symbol string) *StreamBuilder {
	return &StreamBuilder{
		scheme: scheme,
		waHost: host,
		symbol: symbol,
	}
}

type StreamBuilder struct {
	scheme stream.WsScheme
	waHost stream.WsHost
	symbol string
	wsPath stream.WsPath
}

func (wa *StreamBuilder) makeStream(path string) stream.StreamInterface {
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

func (wa *StreamBuilder) Klines(interval string) stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@kline_" + interval)
}

func (wa *StreamBuilder) ContinuousKlines(interval, contractType string) stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + strings.ToLower(contractType) + "@continuousKline_" + interval)
}

func (wa *StreamBuilder) PartialBookDepths(level DepthStreamLevel, rates ...DepthStreamRate) stream.StreamInterface {
	if len(rates) > 0 {
		return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@depth" + strconv.Itoa(int(level)) + "@" + strconv.Itoa(int(rates[0])) + "ms")
	}
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@depth" + strconv.Itoa(int(level)))
}

func (wa *StreamBuilder) DiffBookDepths(rates ...DepthStreamRate) stream.StreamInterface {
	if len(rates) > 0 {
		return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@depth@" + strconv.Itoa(int(rates[0])) + "ms")
	}
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@depth")
}

func (wa *StreamBuilder) AggTrades() stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@aggTrade")
}

func (wa *StreamBuilder) Trades() stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@trade")
}

func (wa *StreamBuilder) BookTickers() stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@bookTicker")
}

func (wa *StreamBuilder) Tickers() stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@ticker")
}

func (wa *StreamBuilder) MiniTickers() stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@miniTicker")
}

func (wa *StreamBuilder) UserData(listenKey string) stream.StreamInterface {
	return wa.makeStream("/" + listenKey)
}

func (wa *StreamBuilder) MarkPrice() stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@markPrice")
}

func (wa *StreamBuilder) LiquidationOrder() stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "@forceOrder")
}

func (wa *StreamBuilder) ContractInfo() stream.StreamInterface {
	return wa.makeStream("/" + strings.ToLower(wa.symbol) + "!contractInfo")
}

func (wa *StreamBuilder) Stream() stream.StreamInterface {
	return wa.makeStream("")
}
