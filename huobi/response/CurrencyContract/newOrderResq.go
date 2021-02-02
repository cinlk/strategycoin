package CurrencyContract

type NewOrderResq struct {
	Status string  `json:"status"`
	Ts     int64   `json:"ts"`
	Data   orderId `json:"data"`
}

type orderId struct {
	OrderId    int64  `json:"order_id"`
	OrderIdStr string `json:"order_id_str"`
}
