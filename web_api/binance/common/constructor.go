package common_web_api

import (
	"sync"

	web_api "github.com/fr0ster/turbo-restler/web_api"
	signature "github.com/fr0ster/turbo-signer/signature"
)

func (wa *WebApi) Lock() {
	wa.mutex.Lock()
}

func (wa *WebApi) Unlock() {
	wa.mutex.Unlock()
}

func New(host web_api.WsHost, path web_api.WsPath, sign signature.Sign) *WebApi {
	return &WebApi{
		waHost: host,
		waPath: path,
		mutex:  &sync.Mutex{},
		sign:   sign,
	}
}
