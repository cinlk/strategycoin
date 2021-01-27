package contract


type ContractAccountInfo struct {
	Status string `json:"status"`
	Ts int64 `json:"ts"`
	Data  []contractAccountInfoData `json:"data"`
}

type contractAccountInfoData struct {
	Symbol string `json:"symbol"`
	ProfitReal  float64 `json:"profit_real"`
}
