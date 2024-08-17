package common_web_stream

import (
	"sync"

	web_api "github.com/fr0ster/turbo-restler/web_api"
)

func New(host web_api.WsHost) *WebStream {
	return &WebStream{
		waHost: host,
		mutex:  &sync.Mutex{},
	}
}
