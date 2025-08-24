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
		wsScheme:       scheme,
		wsHost:         host,
		EndpointPrefix: endpoint,
		EndpointSuffix: "",
		symbol:         symbol[0],
	}
}

type StreamBuilder struct {
	wsScheme       common.WsScheme
	wsHost         common.WsHost
	EndpointPrefix common.WsEndpoint
	EndpointSuffix common.WsEndpoint
	symbol         string
	// optional overrides
	dialer        *websocket.Dialer
	socketFactory func(url string) (web_socket.WebSocketCommonInterface, error)
}

func (wa *StreamBuilder) makeStream() stream.StreamInterface {
	factory := func() (web_socket.WebSocketCommonInterface, error) {
		url := string(wa.wsScheme) + "://" + string(wa.wsHost)
		if wa.EndpointPrefix != "" {
			url = url + "/" + string(wa.EndpointPrefix)
		}
		if wa.EndpointSuffix != "" {
			url = url + "/" + string(wa.EndpointSuffix)
		}
		// Prefer a user-provided socket factory if set
		if wa.socketFactory != nil {
			return wa.socketFactory(url)
		}
		// Otherwise use a custom dialer if provided, falling back to DefaultDialer
		d := wa.dialer
		if d == nil {
			d = websocket.DefaultDialer
		}
		ws, err := web_socket.NewWebSocketWrapper(d, url)
		if err != nil {
			return nil, fmt.Errorf("websocket dial failed url=%s: %w", url, err)
		}
		return ws, nil
	}

	return stream.NewStreamWrapper(
		factory,
		wa.wsScheme,
		wa.wsHost,
		wa.EndpointPrefix,
	)
}

func (wa *StreamBuilder) Klines(interval string) *StreamBuilder {
	wa.EndpointSuffix = common.WsEndpoint("@kline_" + interval)
	return wa
}

func (wa *StreamBuilder) ContinuousKlines(interval, contractType string) *StreamBuilder {
	wa.EndpointSuffix = common.WsEndpoint(strings.ToLower(contractType) + "@continuousKline_" + interval)
	return wa
}

func (wa *StreamBuilder) PartialBookDepths(level DepthStreamLevel, rates ...DepthStreamRate) *StreamBuilder {
	if len(rates) > 0 {
		wa.EndpointSuffix = common.WsEndpoint(strings.ToLower(wa.symbol) + "@depth" + strconv.Itoa(int(level)) + "@" + strconv.Itoa(int(rates[0])) + "ms")
		return wa
	}
	wa.EndpointSuffix = common.WsEndpoint(strings.ToLower(wa.symbol) + "@depth" + strconv.Itoa(int(level)))
	return wa
}

func (wa *StreamBuilder) DiffBookDepths(rates ...DepthStreamRate) *StreamBuilder {
	if len(rates) > 0 {
		wa.EndpointSuffix = common.WsEndpoint(strings.ToLower(wa.symbol) + "@depth@" + strconv.Itoa(int(rates[0])) + "ms")
		return wa
	}
	wa.EndpointSuffix = common.WsEndpoint(strings.ToLower(wa.symbol) + "@depth")
	return wa
}

func (wa *StreamBuilder) AggTrades() *StreamBuilder {
	wa.EndpointSuffix = common.WsEndpoint(strings.ToLower(wa.symbol) + "@aggTrade")
	return wa
}

func (wa *StreamBuilder) Trades() *StreamBuilder {
	wa.EndpointSuffix = common.WsEndpoint(strings.ToLower(wa.symbol) + "@trade")
	return wa
}

func (wa *StreamBuilder) BookTickers() *StreamBuilder {
	wa.EndpointSuffix = common.WsEndpoint(strings.ToLower(wa.symbol) + "@bookTicker")
	return wa
}

func (wa *StreamBuilder) Tickers() *StreamBuilder {
	wa.EndpointSuffix = common.WsEndpoint(strings.ToLower(wa.symbol) + "@ticker")
	return wa
}

func (wa *StreamBuilder) MiniTickers() *StreamBuilder {
	wa.EndpointSuffix = common.WsEndpoint(strings.ToLower(wa.symbol) + "@miniTicker")
	return wa
}

func (wa *StreamBuilder) UserData(listenKey string) *StreamBuilder {
	wa.EndpointSuffix = common.WsEndpoint("/" + listenKey)
	return wa
}

func (wa *StreamBuilder) MarkPrice() *StreamBuilder {
	wa.EndpointSuffix = common.WsEndpoint(strings.ToLower(wa.symbol) + "@markPrice")
	return wa
}

func (wa *StreamBuilder) LiquidationOrder() *StreamBuilder {
	wa.EndpointSuffix = common.WsEndpoint(strings.ToLower(wa.symbol) + "@forceOrder")
	return wa
}

func (wa *StreamBuilder) ContractInfo() *StreamBuilder {
	wa.EndpointSuffix = common.WsEndpoint(strings.ToLower(wa.symbol) + "@markPrice")
	return wa
}

func (wa *StreamBuilder) SetSymbol(symbol string) stream.StreamInterface {
	wa.symbol = symbol
	if symbol != "" {
		wa.EndpointSuffix = common.WsEndpoint(fmt.Sprintf("%s%s", strings.ToLower(symbol), wa.EndpointSuffix))
	}
	return wa.makeStream()
}

func (wa *StreamBuilder) Stream() stream.StreamInterface {
	return wa.makeStream()
}

// WithDialer allows to set a custom websocket dialer (proxy/TLS/etc.).
func (wa *StreamBuilder) WithDialer(d *websocket.Dialer) *StreamBuilder {
	wa.dialer = d
	return wa
}

// WithFactory allows to override socket creation completely. The function
// receives the fully assembled websocket URL and must return a ready wrapper.
func (wa *StreamBuilder) WithFactory(f func(url string) (web_socket.WebSocketCommonInterface, error)) *StreamBuilder {
	wa.socketFactory = f
	return wa
}
