package CurrencyContract

type ContractInfo struct {
	Status string `json:"status"`
	Ts     int64  `json:"ts"`
	Data   []ContractInfoData `json:"data"`
}

type ContractInfoData struct {
	Symbol         string  `json:"symbol"`
	ContractCode   string  `json:"contract_code"`
	ContractSize   float64 `json:"contract_size"`
	PriceTick      float64 `json:"price_tick"`
	DeliveryTime   string  `json:"delivery_time"`
	CreateDate     string  `json:"create_date"`
	ContractStatus int     `json:"contract_status"`
	SettlementDate string  `json:"settlement_date"`
}
