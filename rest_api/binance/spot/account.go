package spot_rest_api

import (
	"net/http"

	rest_api "github.com/fr0ster/turbo-restler/rest_api"
)

func (ra *RestApi) Account() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/account", ra.sign)
	return
}

func (ra *RestApi) QueryCurrentOrderCountUsage() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/rateLimit/order", ra.sign)
	return
}

func (ra *RestApi) QueryPreventedMatches() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/myPreventedMatches", ra.sign)
	return
}

func (ra *RestApi) QueryAllocations() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/myAllocations", ra.sign)
	return
}

func (ra *RestApi) QueryCommissionRates() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/api/v3/account/commission", ra.sign)
	return
}
