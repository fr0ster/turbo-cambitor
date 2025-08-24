package futures_web_stream

import (
	common "github.com/fr0ster/turbo-cambitor/common"
	builder "github.com/fr0ster/turbo-cambitor/web_stream/binance/common"
	stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream"
)

type WebStream interface {
	Klines(interval string) *builder.StreamBuilder
	ContinuousKlines(interval string, contractType string) *builder.StreamBuilder
	PartialBookDepths(level builder.DepthStreamLevel, rates ...builder.DepthStreamRate) *builder.StreamBuilder
	DiffBookDepths(rates ...builder.DepthStreamRate) *builder.StreamBuilder
	AggTrades() *builder.StreamBuilder
	Trades() *builder.StreamBuilder
	BookTickers() *builder.StreamBuilder
	Tickers() *builder.StreamBuilder
	MiniTickers() *builder.StreamBuilder
	UserData(listenKey string) *builder.StreamBuilder
	MarkPrice() *builder.StreamBuilder
	LiquidationOrder() *builder.StreamBuilder
	ContractInfo() *builder.StreamBuilder
	SetSymbol(symbol string) stream.CambitorInterface
	Stream() stream.CambitorInterface
}

func NewDefault(useTestNet ...bool) WebStream {
	var (
		wsScheme   common.WsScheme
		wsHost     common.WsHost
		wsEndpoint common.WsEndpoint
	)
	if len(useTestNet) == 0 {
		useTestNet = append(useTestNet, false)
	}
	if useTestNet[0] {
		wsHost = "fstream.binance.com"
		wsEndpoint = "ws"
	} else {
		wsHost = "fstream.binance.com"
		wsEndpoint = "ws"
	}
	wsScheme = common.WsSchemeWSS
	return builder.NewStreamBuilder(wsScheme, wsHost, wsEndpoint)
}

func New(host string, endpoint string, scheme string) WebStream {
	return builder.NewStreamBuilder(common.WsScheme(scheme), common.WsHost(host), common.WsEndpoint(endpoint))
}
