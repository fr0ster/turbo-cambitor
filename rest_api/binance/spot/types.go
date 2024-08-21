package spot_rest_api

import (
	"sync"

	rest_api "github.com/fr0ster/turbo-restler/rest_api"
	signature "github.com/fr0ster/turbo-signer/signature"
)

type (
	RestApi struct {
		apiBaseUrl rest_api.ApiBaseUrl
		mutex      *sync.Mutex
		sign       signature.Sign
	}
)
