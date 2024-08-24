package common_rest_api

import (
	"net/http"
	"sync"

	"github.com/bitly/go-simplejson"
	"github.com/fr0ster/turbo-restler/rest_api"
	signature "github.com/fr0ster/turbo-signer/signature"
)

func (ra *RestApiWrapper) Call(req *http.Request) (response *simplejson.Json, err error) {
	ra.mutex.Lock()
	defer ra.mutex.Unlock()
	return rest_api.CallRestAPI(req)
}

func (ra *RestApiWrapper) GetSignature() signature.Sign {
	return ra.sign
}

func (ra *RestApiWrapper) GetApiBaseUrl() rest_api.ApiBaseUrl {
	return ra.apiBaseUrl
}

func (ra *RestApiWrapper) GetMutex() *sync.Mutex {
	return ra.mutex
}
