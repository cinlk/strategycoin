package CurrencyContract

type CancerContractOrderResp struct {

	Status string `json:"status"`
	Ts int64 `json:"ts"`
	Data orderRespData `json:"data"`
	
}

type orderRespData struct {
	
	Errors []errors `json:"errors"`
	Successes string `json:"successes"`
	
}

type errors struct {

	OrderId string `json:"order_id"`
	ErrCode int `json:"err_code"`
	ErrMsg string `json:"err_msg"`
}
