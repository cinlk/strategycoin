package CurrencyContract

type OrderInfoResp struct {
	Status string `json:"status"`
	Ts int64 `json:"ts"`
	Data []order `json:"data"`
	
}



type order struct {

	OrderId int64 `json:"order_id"`
	OrderIdStr string `json:"order_id_str"`
	ContractCode string `json:"contract_code"`
	Symbol string `json:"symbol"`
	LeverRate int `json:"lever_rate"`
	Direction string `json:"direction"`
	Offset string `json:"offset"`
 	Volume int `json:"volume"`
	Price float64 `json:"price"`
	ClientOrderId int64 `json:"client_order_id"`
	CreateAt int64  `json:"create_at"`
	UpdateTime int64 `json:"update_time"`
  	OrderSource string `json:"order_source"`
	OrderPriceType string `json:"order_price_type"`
 	MarginFrozen float64 `json:"margin_frozen"`
	Profit float64 `json:"profit"`
	TradeVolume float64 `json:"trade_volume"`
	TradeTurnover float64 `json:"trade_turnover"`
	Fee float64 `json:"fee"`
	FeeAsset string `json:"fee_asset"`
	TradeAvgPrice float64 `json:"trade_avg_price"`
	Status int `json:"status"`
	OrderType int `json:"order_type"`
	CanceledAt int64 `json:"canceled_at"`
	LiquidationType string `json:"liquidation_type"`
	IsTpsl int `json:"is_tpsl"`
	RealProfit float64 `json:"real_profit"`
	
}