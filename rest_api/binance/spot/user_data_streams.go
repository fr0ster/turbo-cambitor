package spot_rest_api

import (
	"fmt"
	"net/http"

	rest_api "github.com/fr0ster/turbo-restler/rest_api"
)

func (ra *RestApi) ListenKey() (listenKey string, err error) {
	response, err := rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodPost, nil, "/api/v3/userDataStream", ra.sign)
	if err != nil {
		err = fmt.Errorf("error calling API: %v", err)
		return
	}
	listenKey = response.Get("listenKey").MustString()
	return
}

func (ra *RestApi) KeepAliveListenKey() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodPut, nil, "/api/v3/userDataStream", ra.sign)
	return
}

func (ra *RestApi) CloseListenKey() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodDelete, nil, "/api/v3/userDataStream", ra.sign)
	return
}
