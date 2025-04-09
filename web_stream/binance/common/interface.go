package common_web_stream

import (
	stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream"
)

// DepthStreamLevel represents depth stream precision (e.g. 5, 10, 20).
type DepthStreamLevel int

// DepthStreamRate represents update frequency (e.g. 1000ms or 100ms).
type DepthStreamRate int

// StreamBuilderInterface defines a builder interface for creating various stream wrappers.
type StreamBuilderInterface interface {
	// Klines returns a kline (candlestick) stream for the given interval.
	Klines(interval string) *stream.StreamWrapper

	// ContinuousKlines returns continuous contract kline stream.
	ContinuousKlines(interval string, contractType string) *stream.StreamWrapper

	// PartialBookDepths returns a partial order book depth stream.
	PartialBookDepths(level DepthStreamLevel, rates ...DepthStreamRate) *stream.StreamWrapper

	// DiffBookDepths returns a diff depth stream.
	DiffBookDepths(rates ...DepthStreamRate) *stream.StreamWrapper

	// AggTrades returns a stream of aggregated trades.
	AggTrades() *stream.StreamWrapper

	// Trades returns a stream of raw trades.
	Trades() *stream.StreamWrapper

	// BookTickers returns a stream of best bid/ask prices.
	BookTickers() *stream.StreamWrapper

	// Tickers returns a stream of full ticker updates.
	Tickers() *stream.StreamWrapper

	// MiniTickers returns a stream of mini ticker updates.
	MiniTickers() *stream.StreamWrapper

	// UserData returns a stream bound to a listenKey (user-specific data).
	UserData(listenKey string) *stream.StreamWrapper

	// MarkPrice returns a stream for mark price updates.
	MarkPrice() *stream.StreamWrapper

	// LiquidationOrder returns a stream for liquidation updates.
	LiquidationOrder() *stream.StreamWrapper

	// ContractInfo returns a stream for perpetual contract metadata.
	ContractInfo() *stream.StreamWrapper

	// Stream returns a generic custom stream with no predefined path.
	Stream() *stream.StreamWrapper
}
