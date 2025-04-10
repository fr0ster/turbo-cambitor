package spot_web_stream

import (
	common "github.com/fr0ster/turbo-cambitor/common"
	builder "github.com/fr0ster/turbo-cambitor/web_stream/binance/common"
	stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream"
)

type WebStream interface {
	Klines(interval string) stream.StreamInterface
	ContinuousKlines(interval string, contractType string) stream.StreamInterface
	PartialBookDepths(level builder.DepthStreamLevel, rates ...builder.DepthStreamRate) stream.StreamInterface
	DiffBookDepths(rates ...builder.DepthStreamRate) stream.StreamInterface
	AggTrades() stream.StreamInterface
	Trades() stream.StreamInterface
	BookTickers() stream.StreamInterface
	Tickers() stream.StreamInterface
	MiniTickers() stream.StreamInterface
	UserData(listenKey string) stream.StreamInterface
	MarkPrice() stream.StreamInterface
	LiquidationOrder() stream.StreamInterface
	ContractInfo() stream.StreamInterface
	Stream() stream.StreamInterface
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
		wsHost = "testnet.binance.vision"
		wsEndpoint = "/ws"
	} else {
		wsHost = "stream.binance.com:9443"
		wsEndpoint = "/ws"
	}
	wsScheme = common.WsSchemeWSS
	return builder.NewStreamBuilder(wsScheme, wsHost, wsEndpoint)
}
