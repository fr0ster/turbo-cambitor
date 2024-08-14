package spot_web_stream

import (
	common "github.com/fr0ster/turbo-cambitor/web_stream/binance/common"
	stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream"

	web_api "github.com/fr0ster/turbo-restler/web_api"
)

type WebStream interface {
	Klines(interval string) *stream.Stream
	ContinuousKlines(interval string, contractType string) *stream.Stream
	PartialBookDepths(level common.DepthStreamLevel, rates ...common.DepthStreamRate) *stream.Stream
	DiffBookDepths(rates ...common.DepthStreamRate) *stream.Stream
	AggTrades() *stream.Stream
	Trades() *stream.Stream
	BookTickers() *stream.Stream
	Tickers() *stream.Stream
	MiniTickers() *stream.Stream
	UserData(listenKey string) *stream.Stream
	MarkPrice() *stream.Stream
	LiquidationOrder() *stream.Stream
	ContractInfo() *stream.Stream
}

func New(symbol string, useTestNet ...bool) WebStream {
	var (
		wsEndpoint string
	)
	if len(useTestNet) == 0 {
		useTestNet = append(useTestNet, false)
	}
	if useTestNet[0] {
		wsEndpoint = "testnet.binance.vision/ws"
	} else {
		wsEndpoint = "stream.binance.com:9443/ws"
	}
	return common.New(web_api.WsHost(wsEndpoint), symbol)
}
