package common_web_api

import (
	"sync"
	"time"

	common "github.com/fr0ster/turbo-cambitor/common"
	web_socket "github.com/fr0ster/turbo-restler/web_socket"
	signature "github.com/fr0ster/turbo-signer/v2/signature"
)

const (
	DepthAPILimit5    DepthAPILimit = 5
	DepthAPILimit10   DepthAPILimit = 10
	DepthAPILimit20   DepthAPILimit = 20
	DepthAPILimit50   DepthAPILimit = 50
	DepthAPILimit100  DepthAPILimit = 100
	DepthAPILimit500  DepthAPILimit = 500
	DepthAPILimit1000 DepthAPILimit = 1000
)

type (
	DepthAPILimit int
	WebApiWrapper struct {
		waHost     common.WsHost
		waEndpoint common.WsEndpoint
		waScheme   common.WsScheme
		mutex      *sync.Mutex
		sign       signature.Sign
		connection web_socket.WebSocketClientInterface
		timeout    *time.Duration
	}
)
