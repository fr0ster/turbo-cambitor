package common_web_api

import request "github.com/fr0ster/turbo-cambitor/web_api/binance/common/request"

func (wa *WebApiWrapper) PlaceOrder() *request.RequestBuilder {
	return request.New("order.place", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) CancelOrder() *request.RequestBuilder {
	return request.New("order.cancel", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) QueryOrder() *request.RequestBuilder {
	return request.New("order.status", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) CancelReplaceOrder() *request.RequestBuilder {
	return request.New("order.cancelReplace", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) ModifyOrder() *request.RequestBuilder {
	return request.New("order.modify", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) QueryOpenOrders() *request.RequestBuilder {
	return request.New("openOrders.status", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) QueryAllOrders() *request.RequestBuilder {
	return request.New("orderList.status", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) QueryPositionV2() *request.RequestBuilder {
	return request.New("v2/account.position", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) QueryPosition() *request.RequestBuilder {
	return request.New("account.position", wa.waHost, wa.waPath, wa.sign)
}
