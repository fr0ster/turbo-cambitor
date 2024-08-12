package turbo_cambitor

type (
	UserDataEventReasonType string
	UserDataEventType       string
	OrderSide               string
	OrderType               string
	SideType                string
	TimeInForceType         string
	QuantityType            string
	OrderExecutionType      string
	OrderStatusType         string
	WorkingType             string
	PositionSideType        string
	DepthSide               string

	AccountType string
	MarginType  string
)

const (
	DepthSideAsk DepthSide = "ASK"
	DepthSideBid DepthSide = "BID"
	SideTypeBuy  OrderSide = "BUY"
	SideTypeSell OrderSide = "SELL"
	SideTypeNone OrderSide = "NONE"
	// SpotAccountType is a constant for spot account type.
	// SPOT/MARGIN/ISOLATED_MARGIN/USDT_FUTURE/COIN_FUTURE
	SpotAccountType           AccountType = "SPOT"
	MarginAccountType         AccountType = "MARGIN"
	IsolatedMarginAccountType AccountType = "ISOLATED_MARGIN"
	USDTFutureType            AccountType = "USDT_FUTURE"
	CoinFutureType            AccountType = "COIN_FUTURE"

	// Для USDT_FUTURE/COIN_FUTURE
	CrossMarginType    MarginType = "CROSS"
	IsolatedMarginType MarginType = "ISOLATED"
)
