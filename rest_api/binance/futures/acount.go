package futures_rest_api

import (
	"net/http"

	rest_api "github.com/fr0ster/turbo-restler/rest_api"
)

func (ra *RestApi) AccountV3() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v3/account", ra.sign)
	return
}

func (ra *RestApi) Account() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v2/account", ra.sign)
	return
}

func (ra *RestApi) BalanceV3() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v3/balance", ra.sign)
	return
}

func (ra *RestApi) Balance() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v2/balance", ra.sign)
	return
}

func (ra *RestApi) UserCommissionRate() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/commissionRate", ra.sign)
	return
}

func (ra *RestApi) AccountConfiguration() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/accountConfig", ra.sign)
	return
}

func (ra *RestApi) SymbolConfiguration() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/symbolConfig", ra.sign)
	return
}

func (ra *RestApi) ForceOrders() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/rateLimit/order", ra.sign)
	return
}

func (ra *RestApi) NotionalLeverageBrackets() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/leverageBracket", ra.sign)
	return
}

func (ra *RestApi) MultiAssetsMargin() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/multiAssetsMargin", ra.sign)
	return
}

func (ra *RestApi) PositionMode() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/positionSide/dual", ra.sign)
	return
}

func (ra *RestApi) Income() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/income", ra.sign)
	return
}

func (ra *RestApi) ApiTradingStatus() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/apiTradingStatus", ra.sign)
	return
}

func (ra *RestApi) TransactionHistoryDownloadId() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/income/asyn", ra.sign)
	return
}

func (ra *RestApi) TransactionHistoryDownloadLinkById() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/income/asyn/id", ra.sign)
	return
}

func (ra *RestApi) OrderHistoryDownloadId() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/order/asyn", ra.sign)
	return
}

func (ra *RestApi) OrderHistoryDownloadLinkById() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/order/asyn/id", ra.sign)
	return
}

func (ra *RestApi) TradeHistoryDownloadId() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/trade/asyn", ra.sign)
	return
}

func (ra *RestApi) TradeHistoryDownloadLinkById() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/trade/asyn/id", ra.sign)
	return
}

func (ra *RestApi) ToggleBNBBurn() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodPost, nil, "/fapi/v1/feeBurn", ra.sign)
	return
}

func (ra *RestApi) BNBBurnStatus() (err error) {
	_, err = rest_api.CallRestAPI(ra.apiBaseUrl, http.MethodGet, nil, "/fapi/v1/feeBurn", ra.sign)
	return
}
