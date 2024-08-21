package common_rest_api

import (
	"net/http"

	"github.com/bitly/go-simplejson"

	"github.com/fr0ster/turbo-restler/rest_api"
	signature "github.com/fr0ster/turbo-signer/signature"
)

func Params2Request(
	params *simplejson.Json,
	baseUrl rest_api.ApiBaseUrl,
	endpoint rest_api.EndPoint,
	method rest_api.HttpMethod,
	apiKey string) (req *http.Request, err error) {
	var (
		fullEndpoint string
		paramsStr    string
	)
	if params != nil {
		paramsStr, err = signature.ConvertSimpleJSONToString(params)
	}
	if err != nil {
		return
	}
	if paramsStr != "" {
		fullEndpoint = string(endpoint) + "?" + string(paramsStr)
	} else {
		fullEndpoint = string(endpoint)
	}
	req, err = http.NewRequest(string(method), string(baseUrl)+string(fullEndpoint), nil)
	if err != nil {
		return
	}
	req.Header.Set("X-MBX-APIKEY", apiKey)
	return
}
