package spot_rest_api

import (
	"sync"

	"github.com/bitly/go-simplejson"
	signature "github.com/fr0ster/turbo-signer/signature"
)

func New(apiKey, apiSecret string, symbol string, sign signature.Sign, useTestNet ...bool) (api *RestApi) {
	const (
		BaseAPIMainUrl    = "https://api.binance.com"
		BaseAPITestnetUrl = "https://testnet.binance.vision"
	)
	simpleJson := simplejson.New()
	simpleJson.Set("apiKey", apiKey)
	simpleJson.Set("symbol", symbol)
	api = &RestApi{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		symbol:     symbol,
		mutex:      &sync.Mutex{},
		parameters: simpleJson,
		sign:       sign,
	}
	if len(useTestNet) > 0 && useTestNet[0] {
		api.apiBaseUrl = BaseAPITestnetUrl
	} else {
		api.apiBaseUrl = BaseAPIMainUrl
	}
	return
}
