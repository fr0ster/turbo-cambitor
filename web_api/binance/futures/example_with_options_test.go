package futures_web_api_test

import (
	"time"

	web_api "github.com/fr0ster/turbo-cambitor/web_api/binance/common/web_api"
	futures "github.com/fr0ster/turbo-cambitor/web_api/binance/futures"
	web_socket "github.com/fr0ster/turbo-restler/web_socket"
	signature "github.com/fr0ster/turbo-signer/v2/signature"
)

// Example showing how to pass custom timeouts via WithWebSocketConfig.
func ExampleNewDefaultWithOptions() {
	var sign signature.Sign = nil // provide your implementation
	wa := futures.NewDefaultWithOptions(sign, true,
		web_api.WithWebSocketConfig(web_socket.WebSocketConfig{
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		}),
	)
	_ = wa
	// Output:
}
