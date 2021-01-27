package Cash



type LatestTradeRecord struct {

	Ch string `json:"ch"`
	Status string `json:"status"`
	Ts int64 `json:"ts"`
	Tick *tick `json:"tick"`
}

type tick struct {
	Data []tradeData `json:"data"`
}

type tradeData struct {
 	
	Amount float64 `json:"amount"`
	Direction string `json:"direction"`
	Id int64 `json:"id"`
	Price float64 `json:"price"`
	TradeId int64 `json:"trade-id"`
	Ts int64 `json:"ts"`
}
