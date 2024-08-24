package common_rest_api

import (
	"sync"

	"github.com/fr0ster/turbo-restler/rest_api"
	signature "github.com/fr0ster/turbo-signer/signature"
)

func New(apiBaseUrl rest_api.ApiBaseUrl, sign signature.Sign) *RestApiWrapper {
	return &RestApiWrapper{
		apiBaseUrl: apiBaseUrl,
		mutex:      &sync.Mutex{},
		sign:       sign,
	}
}
