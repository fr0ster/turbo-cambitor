package spot_rest_api

import (
	common_rest_api "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/rest_api"
	signature "github.com/fr0ster/turbo-signer/signature"
)

func New(sign signature.Sign, useTestNet ...bool) (api *RestApiWrapper) {
	const (
		BaseAPIMainUrl    = "https://api.binance.com"
		BaseAPITestnetUrl = "https://testnet.binance.vision"
	)
	if len(useTestNet) > 0 && useTestNet[0] {
		api = &RestApiWrapper{*common_rest_api.New(BaseAPITestnetUrl, sign)}
	} else {
		api = &RestApiWrapper{*common_rest_api.New(BaseAPIMainUrl, sign)}
	}
	return
}
