package spot_web_stream

import (
	common "github.com/fr0ster/turbo-cambitor/web_stream/binance/common"
	stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream"

	"github.com/fr0ster/turbo-restler/web_socket"
)

type WebStream interface {
	Klines(interval string) *stream.StreamWrapper
	ContinuousKlines(interval string, contractType string) *stream.StreamWrapper
	PartialBookDepths(level common.DepthStreamLevel, rates ...common.DepthStreamRate) *stream.StreamWrapper
	DiffBookDepths(rates ...common.DepthStreamRate) *stream.StreamWrapper
	AggTrades() *stream.StreamWrapper
	Trades() *stream.StreamWrapper
	BookTickers() *stream.StreamWrapper
	Tickers() *stream.StreamWrapper
	MiniTickers() *stream.StreamWrapper
	UserData(listenKey string) *stream.StreamWrapper
	MarkPrice() *stream.StreamWrapper
	LiquidationOrder() *stream.StreamWrapper
	ContractInfo() *stream.StreamWrapper
	Stream() *stream.StreamWrapper
}

func New(useTestNet ...bool) WebStream {
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
	return common.New(web_socket.WsHost(wsEndpoint))
}
