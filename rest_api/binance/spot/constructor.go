package spot_rest_api

import (
	"sync"

	signature "github.com/fr0ster/turbo-signer/signature"
)

func New(sign signature.Sign, useTestNet ...bool) (api *RestApi) {
	const (
		BaseAPIMainUrl    = "https://api.binance.com"
		BaseAPITestnetUrl = "https://testnet.binance.vision"
	)
	api = &RestApi{
		mutex: &sync.Mutex{},
		sign:  sign,
	}
	if len(useTestNet) > 0 && useTestNet[0] {
		api.apiBaseUrl = BaseAPITestnetUrl
	} else {
		api.apiBaseUrl = BaseAPIMainUrl
	}
	return
}
