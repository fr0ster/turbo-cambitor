package futures_web_stream

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

	// Lock()
	// Unlock()
}

func New(useTestNet ...bool) WebStream {
	var (
		waScheme = common.WsSchemeWSS
		// waHost   common.WsHost
		// waPath   common.WsPath
		// wsEndpoint = common.WsEndpoint
		waHost     common.WsHost
		waPath     common.WsPath
		wsEndpoint string
	)
	if len(useTestNet) == 0 {
		useTestNet = append(useTestNet, false)
	}
	if useTestNet[0] {
		wsEndpoint = "fstream.binancefuture.com/ws"
	} else {
		waHost = "fstream.binance.com"
		waPath = "/ws"
		// wsEndpoint = "fstream.binance.com/ws"
	}
	return builder.NewStreamBuilder(waScheme, waHost)
}
