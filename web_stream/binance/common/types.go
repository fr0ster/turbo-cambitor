package common_web_stream

import (
	"sync"
	"time"

	"github.com/fr0ster/turbo-restler/web_socket"
)

const (
	DepthStreamLevel5    DepthStreamLevel = 5
	DepthStreamLevel10   DepthStreamLevel = 10
	DepthStreamLevel20   DepthStreamLevel = 20
	DepthStreamRate100ms DepthStreamRate  = DepthStreamRate(100 * time.Millisecond)
	DepthStreamRate250ms DepthStreamRate  = DepthStreamRate(250 * time.Millisecond)
	DepthStreamRate500ms DepthStreamRate  = DepthStreamRate(500 * time.Millisecond)
)

type (
	DepthStreamLevel int
	DepthStreamRate  time.Duration
	WebStream        struct {
		symbol string
		waHost web_socket.WsHost
		mutex  *sync.Mutex
	}
)
