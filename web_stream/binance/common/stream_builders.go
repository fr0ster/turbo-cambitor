package common_web_stream

import (
	"fmt"
	"strconv"
	"strings"

	common "github.com/fr0ster/turbo-cambitor/common"
	stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream"
	web_socket "github.com/fr0ster/turbo-restler/web_socket"

	"github.com/gorilla/websocket"
)

// NewWebStream створює новий екземпляр StreamBuilder з переданими параметрами.
func NewStreamBuilder(scheme common.WsScheme, host common.WsHost, endpoint common.WsEndpoint, symbol ...string) *StreamBuilder {
	if len(symbol) == 0 {
		symbol = append(symbol, "")
	}
	return &StreamBuilder{
		wsScheme:   scheme,
		wsHost:     common.WsHost(string(host) + string(endpoint)),
		wsEndpoint: "",
		symbol:     symbol[0],
	}
}

type StreamBuilder struct {
	wsScheme   common.WsScheme
	wsHost     common.WsHost
	wsEndpoint common.WsEndpoint
	symbol     string
}

func (wa *StreamBuilder) makeStream() stream.StreamInterface {
	factory := func() (web_socket.WebSocketInterface, error) {
		url := ""
		if wa.wsEndpoint != "" {
			url = string(wa.wsScheme) + "://" + string(wa.wsHost) + "/" + string(wa.wsEndpoint)
		} else {
			url = string(wa.wsScheme) + "://" + string(wa.wsHost)
		}
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			return nil, err
		}
		return web_socket.NewWebSocketWrapper(conn), nil
	}

	return stream.NewStreamWrapper(
		factory,
		wa.wsScheme,
		wa.wsHost,
		wa.wsEndpoint,
	)
}

func (wa *StreamBuilder) Klines(interval string) *StreamBuilder {
	wa.wsEndpoint = common.WsEndpoint("@kline_" + interval)
	return wa
}

func (wa *StreamBuilder) ContinuousKlines(interval, contractType string) *StreamBuilder {
	wa.wsEndpoint = common.WsEndpoint(strings.ToLower(contractType) + "@continuousKline_" + interval)
	return wa
}

func (wa *StreamBuilder) PartialBookDepths(level DepthStreamLevel, rates ...DepthStreamRate) *StreamBuilder {
	if len(rates) > 0 {
		wa.wsEndpoint = common.WsEndpoint(strings.ToLower(wa.symbol) + "@depth" + strconv.Itoa(int(level)) + "@" + strconv.Itoa(int(rates[0])) + "ms")
		return wa
	}
	wa.wsEndpoint = common.WsEndpoint(strings.ToLower(wa.symbol) + "@depth" + strconv.Itoa(int(level)))
	return wa
}

func (wa *StreamBuilder) DiffBookDepths(rates ...DepthStreamRate) *StreamBuilder {
	if len(rates) > 0 {
		wa.wsEndpoint = common.WsEndpoint(strings.ToLower(wa.symbol) + "@depth@" + strconv.Itoa(int(rates[0])) + "ms")
		return wa
	}
	wa.wsEndpoint = common.WsEndpoint(strings.ToLower(wa.symbol) + "@depth")
	return wa
}

func (wa *StreamBuilder) AggTrades() *StreamBuilder {
	wa.wsEndpoint = common.WsEndpoint(strings.ToLower(wa.symbol) + "@aggTrade")
	return wa
}

func (wa *StreamBuilder) Trades() *StreamBuilder {
	wa.wsEndpoint = common.WsEndpoint(strings.ToLower(wa.symbol) + "@trade")
	return wa
}

func (wa *StreamBuilder) BookTickers() *StreamBuilder {
	wa.wsEndpoint = common.WsEndpoint(strings.ToLower(wa.symbol) + "@bookTicker")
	return wa
}

func (wa *StreamBuilder) Tickers() *StreamBuilder {
	wa.wsEndpoint = common.WsEndpoint(strings.ToLower(wa.symbol) + "@ticker")
	return wa
}

func (wa *StreamBuilder) MiniTickers() *StreamBuilder {
	wa.wsEndpoint = common.WsEndpoint(strings.ToLower(wa.symbol) + "@miniTicker")
	return wa
}

func (wa *StreamBuilder) UserData(listenKey string) *StreamBuilder {
	wa.wsEndpoint = common.WsEndpoint("/" + listenKey)
	return wa
}

func (wa *StreamBuilder) MarkPrice() *StreamBuilder {
	wa.wsEndpoint = common.WsEndpoint(strings.ToLower(wa.symbol) + "@markPrice")
	return wa
}

func (wa *StreamBuilder) LiquidationOrder() *StreamBuilder {
	wa.wsEndpoint = common.WsEndpoint(strings.ToLower(wa.symbol) + "@forceOrder")
	return wa
}

func (wa *StreamBuilder) ContractInfo() *StreamBuilder {
	wa.wsEndpoint = common.WsEndpoint(strings.ToLower(wa.symbol) + "@markPrice")
	return wa
}

func (wa *StreamBuilder) SetSymbol(symbol string) stream.StreamInterface {
	wa.symbol = symbol
	if symbol != "" {
		wa.wsEndpoint = common.WsEndpoint(fmt.Sprintf("%s%s", strings.ToLower(symbol), wa.wsEndpoint))
	}
	return wa.makeStream()
}

func (wa *StreamBuilder) Stream() stream.StreamInterface {
	return wa.makeStream()
}
