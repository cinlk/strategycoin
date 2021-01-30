package CurrencyContract


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

	Amount string `json:"amount"`
	Direction string `json:"direction"`
	Id float64 `json:"id"`
	Price string `json:"price"`
	Ts int64 `json:"ts"`
}