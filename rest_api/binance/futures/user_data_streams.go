package futures_rest_api

import (
	"fmt"
	"net/http"

	rest_api "github.com/fr0ster/turbo-restler/rest_api"
)

func (ra *RestApi) ListenKey() (listenKey string, err error) {
	response, err := rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodPost, nil, "/fapi/v1/listenKey", ra.sign)
	if err != nil {
		err = fmt.Errorf("error calling API: %v", err)
		return
	}
	listenKey = response.Get("listenKey").MustString()
	return
}

func (ra *RestApi) KeepAliveListenKey() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodPut, nil, "/fapi/v1/listenKey", ra.sign)
	return
}

func (ra *RestApi) CloseListenKey() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodDelete, nil, "/fapi/v1/listenKey", ra.sign)
	return
}
