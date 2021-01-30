package CurrencyContract



type  ContractIndexInfo struct {

	Status string `json:"status"`
	
	Ts int64 `json:"ts"`
	Data []indexData `json:"data"`
}

type indexData struct {

	IndexPrice float64 `json:"index_price"`
	IndexTs int64 `json:"index_ts"`
	ContractCode string `json:"contract_code"`
}
