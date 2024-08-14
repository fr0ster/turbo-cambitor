package common_web_api

import request "github.com/fr0ster/turbo-cambitor/web_api/binance/common/request"

func (wa *WebApi) PlaceOrder() *request.Request {
	return request.New("order.place", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) CancelOrder() *request.Request {
	return request.New("order.cancel", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) QueryOrder() *request.Request {
	return request.New("order.status", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) CancelReplaceOrder() *request.Request {
	return request.New("order.cancelReplace", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) ModifyOrder() *request.Request {
	return request.New("order.modify", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) QueryOpenOrders() *request.Request {
	return request.New("openOrders.status", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) QueryAllOrders() *request.Request {
	return request.New("orderList.status", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) QueryPositionV2() *request.Request {
	return request.New("v2/account.position", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) QueryPosition() *request.Request {
	return request.New("account.position", wa.waHost, wa.waPath, wa.sign)
}
