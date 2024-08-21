package futures_rest_api

import (
	"net/http"
	"strconv"

	"github.com/bitly/go-simplejson"
	request "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/request"
)

func (ra *RestApi) AccountV3() *request.Request {
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v3/account", nil, ra.sign)
	return rq.SetAPIKey()
}

func (ra *RestApi) Account() *request.Request {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v2/account", params, ra.sign)
	return rq.SetAPIKey()
}

func (ra *RestApi) BalanceV3() *request.Request {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v3/balance", params, ra.sign)
	return rq.SetAPIKey()
}

func (ra *RestApi) Balance() *request.Request {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v2/balance", params, ra.sign)
	return rq.SetAPIKey()
}

func (ra *RestApi) UserCommissionRate() *request.Request {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/commissionRate", params, ra.sign)
	return rq.SetAPIKey()
}

func (ra *RestApi) AccountConfiguration() *request.Request {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/accountConfig", params, ra.sign)
	return rq.SetAPIKey()
}

func (ra *RestApi) SymbolConfiguration() *request.Request {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/symbolConfig", params, ra.sign)
	return rq.SetAPIKey()
}

func (ra *RestApi) ForceOrders() *request.Request {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/rateLimit/order", params, ra.sign)
	return rq.SetAPIKey()
}

func (ra *RestApi) NotionalLeverageBrackets() *request.Request {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/leverageBracket", params, ra.sign)
	return rq.SetAPIKey()
}

func (ra *RestApi) MultiAssetsMargin() *request.Request {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/multiAssetsMargin", params, ra.sign)
	return rq.SetAPIKey()
}

func (ra *RestApi) PositionMode() *request.Request {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/positionSide/dual", params, ra.sign)
	return rq.SetAPIKey()
}

func (ra *RestApi) Income() *request.Request {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/income", params, ra.sign)
	return rq.SetAPIKey()
}

func (ra *RestApi) ToggleBNBBurn(feeBurn bool) *request.Request {
	params := simplejson.New()
	rq := request.New(http.MethodPost, ra.apiBaseUrl, "/fapi/v1/feeBurn", params, ra.sign)
	return rq.SetAPIKey().Set("feeBurn", strconv.FormatBool(feeBurn))
}

func (ra *RestApi) BNBBurnStatus() *request.Request {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.apiBaseUrl, "/fapi/v1/feeBurn", params, ra.sign)
	return rq.SetAPIKey()
}
