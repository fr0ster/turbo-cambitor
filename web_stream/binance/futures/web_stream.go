package futures_web_stream

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

	Lock()
	Unlock()
}

func New(useTestNet ...bool) WebStream {
	var (
		wsEndpoint string
	)
	if len(useTestNet) == 0 {
		useTestNet = append(useTestNet, false)
	}
	if useTestNet[0] {
		wsEndpoint = "fstream.binancefuture.com/ws"
	} else {
		wsEndpoint = "fstream.binance.com/ws"
	}
	return common.New(web_socket.WsHost(wsEndpoint))
}
