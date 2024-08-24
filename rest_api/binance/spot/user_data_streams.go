package spot_rest_api

import (
	"net/http"

	"github.com/bitly/go-simplejson"
	request "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/request"
)

func (ra *RestApiWrapper) ListenKey() (listenKey string, err error) {
	rq, err := request.New(
		http.MethodPost,
		ra.GetApiBaseUrl(),
		"/api/v3/userDataStream",
		nil, ra.GetSignature()).Do()
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

func (ra *RestApiWrapper) KeepAliveListenKey(listenKey string) (err error) {
	params := simplejson.New()
	params.Set("listenKey", listenKey)
	rq, err := request.New(http.MethodPut, ra.GetApiBaseUrl(), "/api/v3/userDataStream", params, ra.GetSignature()).Do()
	if err != nil {
		return
	}
	_, err = ra.Call(rq)
	return
}

func (ra *RestApiWrapper) CloseListenKey(listenKey string) (err error) {
	params := simplejson.New()
	params.Set("listenKey", listenKey)
	rq, err := request.New(http.MethodDelete, ra.GetApiBaseUrl(), "/api/v3/userDataStream", params, ra.GetSignature()).Do()
	if err != nil {
		return
	}
	_, err = ra.Call(rq)
	return
}
