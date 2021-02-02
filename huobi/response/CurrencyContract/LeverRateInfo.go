package CurrencyContract

type LeverRate struct {
	Status string        `json:"status"`
	Ts     int64         `json:"ts"`
	Data   leverRateData `json:"data"`
}

type leverRateData struct {
	ContractCode string `json:"contract_code"`
	LeverRate    int    `json:"lever_rate"`
}
