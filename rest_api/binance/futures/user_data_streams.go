package futures_rest_api

import (
	"net/http"

	request "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/request"
)

func (ra *RestApi) ListenKey() (listenKey string, err error) {
	rq := request.New(http.MethodPost, ra.apiBaseUrl, "/fapi/v1/listenKey", nil, ra.sign)
	response, err := rq.SetAPIKey().SetTimestamp().SetSignature().Do()
	if err != nil {
		return
	}
	listenKey = response.Get("listenKey").MustString()
	return
}

func (ra *RestApi) KeepAliveListenKey() (listenKey string, err error) {
	rq := request.New(http.MethodPut, ra.apiBaseUrl, "/fapi/v1/listenKey", nil, ra.sign)
	response, err := rq.SetAPIKey().SetTimestamp().SetSignature().Do()
	if err != nil {
		return
	}
	listenKey = response.Get("listenKey").MustString()
	return
}

func (ra *RestApi) CloseListenKey() (err error) {
	rq := request.New(http.MethodDelete, ra.apiBaseUrl, "/fapi/v1/listenKey", nil, ra.sign)
	_, err = rq.SetAPIKey().SetTimestamp().SetSignature().Do()
	return
}
