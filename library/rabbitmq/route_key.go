package rabbitmq

type RouteKey string

var (
	TradeCreate = RouteKey("tradeCreate")
	TradeChange = RouteKey("tradeChange")
)
