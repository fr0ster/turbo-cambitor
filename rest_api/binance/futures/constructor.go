package futures_rest_api

import (
	"sync"

	signature "github.com/fr0ster/turbo-signer/signature"
)

func New(sign signature.Sign, useTestNet ...bool) (api *RestApi) {
	const (
		BaseAPIMainUrl    = "https://fapi.binance.com"
		BaseAPITestnetUrl = "https://testnet.binancefuture.com"
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
