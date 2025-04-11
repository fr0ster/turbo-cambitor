package futures_web_stream_test

import (
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
	"/ws/btcusdt@kline_1m":                    `{"stream":"kline_1m","data":"mock_kline"}`,
	"/ws/btcusdt@aggTrade":                    `{"stream":"aggTrade","data":"mock_agg"}`,
	"/ws/btcusdt@trade":                       `{"stream":"trade","data":"mock_trade"}`,
	"/ws/btcusdt@miniTicker":                  `{"stream":"miniTicker","data":"mock_mini"}`,
	"/ws/btcusdt@ticker":                      `{"stream":"ticker","data":"mock_ticker"}`,
	"/ws/btcusdt@bookTicker":                  `{"stream":"bookTicker","data":"mock_book"}`,
	"/ws/btcusdt@depth5@100ms":                `{"stream":"depth5","data":"mock_depth5"}`,
	"/ws/btcusdt@depth@100ms":                 `{"stream":"depth","data":"mock_depth"}`,
	"/ws/btcusdt@markPrice":                   `{"stream":"markPrice","data":"mock_mark"}`,
	"/ws/btcusdt@forceOrder":                  `{"stream":"liquidation","data":"mock_liq"}`,
	"/ws/btcusdtperpetual@continuousKline_1m": `{"stream":"ckline","data":"mock_cont_kline"}`,
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

		path := r.URL.Path[len("/ws"):]
		response := ResponseByPath["/ws"+path]
		time.Sleep(100 * time.Millisecond)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(response))
	})

	server := &http.Server{Addr: addr, Handler: mux}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Mock WS server error: %v", err)
		}
	}()
	return server
}

func checkStreamMessage(t *testing.T, stream streamer.StreamInterface, expect string) {
	msgChan := make(chan string, 1)
	err := stream.Connect()
	stream.SetMessageLogger(func(msg web_socket.LogRecord) {
		select {
		case msgChan <- string(msg.Body):
		default:
		}
	})
	require.NoError(t, err, "stream.Connect() failed")
	defer stream.Disconnect()

	select {
	case body := <-msgChan:
		require.Contains(t, body, expect)
	case <-time.After(1 * time.Second):
		t.Fatalf("Timeout: no message received")
	}
}

// func TestKlines is replaced by TestKlines_Debug for single-purpose testability
func TestKlines_Debug(t *testing.T) {
	StartMockWSServer(":9090")
	ws := web_stream.New("localhost:9090", "/ws", "ws")
	checkStreamMessage(t, ws.Klines("1m").SetSymbol("btcusdt"), "mock_kline")
}

func TestAggTrades_Debug(t *testing.T) {
	StartMockWSServer(":9090")
	ws := web_stream.New("localhost:9090", "/ws", "ws")
	checkStreamMessage(t, ws.AggTrades().SetSymbol("btcusdt"), "mock_agg")
}

func TestTrades_Debug(t *testing.T) {
	StartMockWSServer(":9090")
	ws := web_stream.New("localhost:9090", "/ws", "ws")
	checkStreamMessage(t, ws.Trades().SetSymbol("btcusdt"), "mock_trade")
}

func TestMiniTickers_Debug(t *testing.T) {
	StartMockWSServer(":9090")
	ws := web_stream.New("localhost:9090", "/ws", "ws")
	checkStreamMessage(t, ws.MiniTickers().SetSymbol("btcusdt"), "mock_mini")
}

func TestTickers_Debug(t *testing.T) {
	StartMockWSServer(":9090")
	ws := web_stream.New("localhost:9090", "/ws", "ws")
	checkStreamMessage(t, ws.Tickers().SetSymbol("btcusdt"), "mock_ticker")
}

func TestBookTickers_Debug(t *testing.T) {
	StartMockWSServer(":9090")
	ws := web_stream.New("localhost:9090", "/ws", "ws")
	checkStreamMessage(t, ws.BookTickers().SetSymbol("btcusdt"), "mock_book")
}

func TestPartialBookDepths_Debug(t *testing.T) {
	StartMockWSServer(":9090")
	ws := web_stream.New("localhost:9090", "/ws", "ws")
	checkStreamMessage(t, ws.PartialBookDepths(5, 100).SetSymbol("btcusdt"), "mock_depth5")
}

func TestDiffBookDepths_Debug(t *testing.T) {
	StartMockWSServer(":9090")
	ws := web_stream.New("localhost:9090", "/ws", "ws")
	checkStreamMessage(t, ws.DiffBookDepths(100).SetSymbol("btcusdt"), "mock_depth")
}

func TestMarkPrice_Debug(t *testing.T) {
	StartMockWSServer(":9090")
	ws := web_stream.New("localhost:9090", "/ws", "ws")
	checkStreamMessage(t, ws.MarkPrice().SetSymbol("btcusdt"), "mock_mark")
}

func TestLiquidationOrder_Debug(t *testing.T) {
	StartMockWSServer(":9090")
	ws := web_stream.New("localhost:9090", "/ws", "ws")
	checkStreamMessage(t, ws.LiquidationOrder().SetSymbol("btcusdt"), "mock_liq")
}

func TestContinuousKlines_Debug(t *testing.T) {
	StartMockWSServer(":9090")
	ws := web_stream.New("localhost:9090", "/ws", "ws")
	checkStreamMessage(t, ws.ContinuousKlines("1m", "perpetual").SetSymbol("btcusdt"), "mock_cont_kline")
}
