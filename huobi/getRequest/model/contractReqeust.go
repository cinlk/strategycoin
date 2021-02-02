package model

type GetContractBasisRequest struct {
	BasisPriceType string
}

type GetContractAccountRequest struct {
	Symbol string `json:"symbol"`
}

type GetContractCodeRequest struct {
	ContractCode string `json:"contract_code"`
}

type CancerOrderReq struct {
	ClientOrderId string `json:"client_order_id"`
	OrderId       string `json:"order_id"`
}

type CreateOrderReq struct {
	ContractCode     string  `json:"contract_code"`
	ClientOrderId    int64   `json:"client_order_id,omitempty"`
	Price            float64 `json:"price"`
	Volume           int64   `json:"volume"`
	Direction        string  `json:"direction"`
	Offset           string  `json:"offset"`
	LeverRate        int     `json:"lever_rate"`
	OrderPriceType   string  `json:"order_price_type"`
	TpTriggerPrice   float64 `json:"tp_trigger_price,omitempty"`
	TpOrderPrice     float64 `json:"tp_order_price,omitempty"`
	TpOrderPriceType string  `json:"tp_order_price_type,omitempty"`
	SlTriggerPrice   float64 `json:"sl_trigger_price,omitempty"`
	SlOrderPrice     float64 `json:"sl_order_price,omitempty"`
	SlOrderPriceType string  `json:"sl_order_price_type,omitempty"`
}
