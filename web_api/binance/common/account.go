package common_web_api

import request "github.com/fr0ster/turbo-cambitor/web_api/binance/common/request"

func (wa *WebApi) AccountBalanceV2() *request.Request {
	return request.New("v2/account.balance", wa.waHost, wa.waPath, wa.sign)
}
func (wa *WebApi) AccountBalance() *request.Request {
	return request.New("account.balance", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) AccountInformationV2() *request.Request {
	return request.New("v2/account.status", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) AccountInformation() *request.Request {
	return request.New("account.status", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) AccountPositionsV2() *request.Request {
	return request.New("v2/account.position", wa.waHost, wa.waPath, wa.sign)
}

func (wa *WebApi) AccountPositions() *request.Request {
	return request.New("account.position", wa.waHost, wa.waPath, wa.sign)
}
