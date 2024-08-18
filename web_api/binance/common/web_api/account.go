package common_web_api

import request "github.com/fr0ster/turbo-cambitor/web_api/binance/common/request"

func (wa *WebApiWrapper) AccountBalanceV2() *request.Request {
	return request.New("v2/account.balance", wa.waHost, wa.waPath, wa.sign)
}
func (wa *WebApiWrapper) AccountBalance() *request.Request {
	return request.New("account.balance", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) AccountInformationV2() *request.Request {
	return request.New("v2/account.status", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) AccountInformation() *request.Request {
	return request.New("account.status", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) AccountPositionsV2() *request.Request {
	return request.New("v2/account.position", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApiWrapper) AccountPositions() *request.Request {
	return request.New("account.position", wa.waHost, wa.waPath, wa.sign)
}
