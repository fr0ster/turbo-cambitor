package futures_web_stream_test

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	streamer "github.com/fr0ster/turbo-cambitor/web_stream/binance/common/stream"
	web_stream "github.com/fr0ster/turbo-cambitor/web_stream/binance/futures"
	"github.com/fr0ster/turbo-restler/web_socket"
)

var upgrader = websocket.Upgrader{}

var ResponseByPath = map[string]string{
	"/ws/btcusdt@kline_1m":                     `{"stream":"kline_1m","data":"mock_kline"}`,
	"/ws/btcusdt@aggTrade":                     `{"stream":"aggTrade","data":"mock_agg"}`,
	"/ws/btcusdt@trade":                        `{"stream":"trade","data":"mock_trade"}`,
	"/ws/btcusdt@miniTicker":                   `{"stream":"miniTicker","data":"mock_mini"}`,
	"/ws/btcusdt@ticker":                       `{"stream":"ticker","data":"mock_ticker"}`,
	"/ws/btcusdt@bookTicker":                   `{"stream":"bookTicker","data":"mock_book"}`,
	"/ws/btcusdt@depth5@100ms":                 `{"stream":"depth5","data":"mock_depth5"}`,
	"/ws/btcusdt@depth@100ms":                  `{"stream":"depth","data":"mock_depth"}`,
	"/ws/!markPrice@arr":                       `{"stream":"markPrice","data":"mock_mark"}`,
	"/ws/!forceOrder@arr":                      `{"stream":"liquidation","data":"mock_liq"}`,
	"/ws/btcusdt_perpetual@continuousKline_1m": `{"stream":"ckline","data":"mock_cont_kline"}`,
}

func StartMockWSServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Upgrade error: %v", err)
			return
		}
		defer conn.Close()

		path := r.URL.Path[len("/ws"):] // e.g. /btcusdt@kline_1m
		response := ResponseByPath[path]

		// Відправляємо мок-повідомлення через 100ms
		time.Sleep(100 * time.Millisecond)
		err = conn.WriteMessage(websocket.TextMessage, []byte(response))
		if err != nil {
			log.Printf("Write error: %v", err)
		}
	})

	return &http.Server{Addr: addr, Handler: mux}
}

type testCase struct {
	name      string
	getStream func(web_stream.WebStream) streamer.StreamInterface
	expect    string
}

func TestAllStreamsWithMockResponses(t *testing.T) {
	server := StartMockWSServer(":9090")
	go server.ListenAndServe()
	defer server.Shutdown(context.Background())

	time.Sleep(200 * time.Millisecond) // Дати серверу стартанути

	ws := web_stream.New("localhost:9090", "/ws", "ws")

	tests := []testCase{
		{"Klines", func(ws web_stream.WebStream) streamer.StreamInterface { return ws.Klines("1m") }, `mock_kline`},
		{"AggTrades", func(ws web_stream.WebStream) streamer.StreamInterface { return ws.AggTrades() }, `mock_agg`},
		{"Trades", func(ws web_stream.WebStream) streamer.StreamInterface { return ws.Trades() }, `mock_trade`},
		{"MiniTickers", func(ws web_stream.WebStream) streamer.StreamInterface { return ws.MiniTickers() }, `mock_mini`},
		{"Tickers", func(ws web_stream.WebStream) streamer.StreamInterface { return ws.Tickers() }, `mock_ticker`},
		{"BookTickers", func(ws web_stream.WebStream) streamer.StreamInterface { return ws.BookTickers() }, `mock_book`},
		{"PartialBookDepths", func(ws web_stream.WebStream) streamer.StreamInterface {
			return ws.PartialBookDepths(5, 100)
		}, `mock_depth5`},
		{"DiffBookDepths", func(ws web_stream.WebStream) streamer.StreamInterface {
			return ws.DiffBookDepths(100)
		}, `mock_depth`},
		{"MarkPrice", func(ws web_stream.WebStream) streamer.StreamInterface { return ws.MarkPrice() }, `mock_mark`},
		{"LiquidationOrder", func(ws web_stream.WebStream) streamer.StreamInterface { return ws.LiquidationOrder() }, `mock_liq`},
		{"ContinuousKlines", func(ws web_stream.WebStream) streamer.StreamInterface {
			return ws.ContinuousKlines("1m", "perpetual")
		}, `mock_cont_kline`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msgChan := make(chan string, 1)

			stream := tc.getStream(ws).
				SetSymbol("btcusdt").
				SetMessageLogger(func(msg web_socket.LogRecord) {
					select {
					case msgChan <- string(msg.Body):
					default:
						// Don't block if the channel is full
					}
				})

			require.NotNil(t, stream)

			err := stream.Connect()
			require.NoError(t, err)
			defer stream.Disconnect()

			select {
			case body := <-msgChan:
				require.Contains(t, body, tc.expect)
			case <-time.After(1 * time.Second):
				t.Fatalf("Timeout: no message received for %s", tc.name)
			}
		})
	}
}
