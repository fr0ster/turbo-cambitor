package common_rest_api

import (
	"sync"

	"github.com/fr0ster/turbo-restler/rest_api"
	signature "github.com/fr0ster/turbo-signer/signature"
)

type (
	RestApiWrapper struct {
		apiBaseUrl rest_api.ApiBaseUrl
		mutex      *sync.Mutex
		sign       signature.Sign
	}
)
