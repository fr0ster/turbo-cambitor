package futures_rest_api

import (
	"net/http"
	"strconv"

	"github.com/bitly/go-simplejson"
	request "github.com/fr0ster/turbo-cambitor/rest_api/binance/common/request"
)

func (ra *RestApiWrapper) Account() *request.RequestBuilder {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v3/account", params, ra.GetSignature())
	return rq.SetAPIKey()
}

func (ra *RestApiWrapper) Balance() *request.RequestBuilder {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v3/balance", params, ra.GetSignature())
	return rq.SetAPIKey()
}

func (ra *RestApiWrapper) PositionRisk() *request.RequestBuilder {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v3/positionRisk", params, ra.GetSignature())
	return rq.SetAPIKey()
}

func (ra *RestApiWrapper) UserCommissionRate() *request.RequestBuilder {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/commissionRate", params, ra.GetSignature())
	return rq.SetAPIKey()
}

func (ra *RestApiWrapper) AccountConfiguration() *request.RequestBuilder {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/accountConfig", params, ra.GetSignature())
	return rq.SetAPIKey()
}

func (ra *RestApiWrapper) SymbolConfiguration() *request.RequestBuilder {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/symbolConfig", params, ra.GetSignature())
	return rq.SetAPIKey()
}

func (ra *RestApiWrapper) ForceOrders() *request.RequestBuilder {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/rateLimit/order", params, ra.GetSignature())
	return rq.SetAPIKey()
}

func (ra *RestApiWrapper) NotionalLeverageBrackets() *request.RequestBuilder {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/leverageBracket", params, ra.GetSignature())
	return rq.SetAPIKey()
}

func (ra *RestApiWrapper) MultiAssetsMargin() *request.RequestBuilder {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/multiAssetsMargin", params, ra.GetSignature())
	return rq.SetAPIKey()
}

func (ra *RestApiWrapper) PositionMode() *request.RequestBuilder {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/positionSide/dual", params, ra.GetSignature())
	return rq.SetAPIKey()
}

func (ra *RestApiWrapper) Income() *request.RequestBuilder {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/income", params, ra.GetSignature())
	return rq.SetAPIKey()
}

func (ra *RestApiWrapper) ToggleBNBBurn(feeBurn bool) *request.RequestBuilder {
	params := simplejson.New()
	rq := request.New(http.MethodPost, ra.GetApiBaseUrl(), "/fapi/v1/feeBurn", params, ra.GetSignature())
	return rq.SetAPIKey().Set("feeBurn", strconv.FormatBool(feeBurn))
}

func (ra *RestApiWrapper) BNBBurnStatus() *request.RequestBuilder {
	params := simplejson.New()
	rq := request.New(http.MethodGet, ra.GetApiBaseUrl(), "/fapi/v1/feeBurn", params, ra.GetSignature())
	return rq.SetAPIKey()
}
