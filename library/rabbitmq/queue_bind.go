package rabbitmq

type queueBind struct {
	name string
	keys []RouteKey
}

var (
	OrderUserRelCreateUpdate = &queueBind{
		name: "Order.userRel.generate",
		keys: []RouteKey{TradeCreate, TradeChange},
	}
	TradeOrderCreateNotice = &queueBind{
		name: "Trade.orderCreateNotice",
		keys: []RouteKey{TradeCreate},
	}
)
