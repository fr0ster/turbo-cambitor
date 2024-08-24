package futures_rest_api

import (
	"net/http"

	request "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/request"
)

func (ra *RestApiWrapper) ListenKey() (listenKey string, err error) {
	rq, err := request.New(
		http.MethodPost,
		ra.GetApiBaseUrl(),
		"/fapi/v1/listenKey",
		nil,
		ra.GetSignature()).SetAPIKey().SetTimestamp().SetSignature().Do()
	if err != nil {
		return
	}
	response, err := ra.Call(rq)
	if err != nil {
		return
	}
	listenKey = response.Get("listenKey").MustString()
	return
}

func (ra *RestApiWrapper) KeepAliveListenKey() (listenKey string, err error) {
	rq, err := request.New(
		http.MethodPut,
		ra.GetApiBaseUrl(),
		"/fapi/v1/listenKey",
		nil,
		ra.GetSignature()).SetAPIKey().SetTimestamp().SetSignature().Do()
	if err != nil {
		return
	}
	response, err := ra.Call(rq)
	if err != nil {
		return
	}
	listenKey = response.Get("listenKey").MustString()
	return
}

func (ra *RestApiWrapper) CloseListenKey() (err error) {
	rq, err := request.New(
		http.MethodDelete,
		ra.GetApiBaseUrl(),
		"/fapi/v1/listenKey",
		nil,
		ra.GetSignature()).SetAPIKey().SetTimestamp().SetSignature().Do()
	if err != nil {
		return
	}
	_, err = ra.Call(rq)
	return
}
