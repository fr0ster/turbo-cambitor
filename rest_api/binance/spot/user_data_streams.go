package spot_rest_api

import (
	"net/http"

	"github.com/bitly/go-simplejson"
	request "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/request"
)

func (ra *RestApi) ListenKey() (listenKey string, err error) {
	rq := request.New(http.MethodPost, ra.apiBaseUrl, "/api/v3/userDataStream", nil, ra.sign)
	response, err := rq.Do()
	if err != nil {
		return
	}
	listenKey = response.Get("listenKey").MustString()
	return
}

func (ra *RestApi) KeepAliveListenKey(listenKey string) (err error) {
	params := simplejson.New()
	params.Set("listenKey", listenKey)
	rq := request.New(http.MethodPut, ra.apiBaseUrl, "/api/v3/userDataStream", params, ra.sign)
	_, err = rq.Do()
	return
}

func (ra *RestApi) CloseListenKey(listenKey string) (err error) {
	params := simplejson.New()
	params.Set("listenKey", listenKey)
	rq := request.New(http.MethodDelete, ra.apiBaseUrl, "/api/v3/userDataStream", params, ra.sign)
	_, err = rq.Do()
	return
}
