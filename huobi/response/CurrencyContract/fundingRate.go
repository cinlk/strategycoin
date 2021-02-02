package CurrencyContract

type FundingRateInfo struct {
	Status string                 `json:"status"`
	Ts     int64                  `json:"ts"`
	Data   *LatestFundingRatedata `json:"data"`
}

type LatestFundingRatedata struct {
	EstimatedRate   string `json:"estimated_rate"`
	FundingRate     string `json:"funding_rate"`
	ContractCode    string `json:"contract_code"`
	Symbol          string `json:"symbol"`
	FeeAsset        string `json:"fee_asset"`
	FundingTime     string `json:"funding_time"`
	NextFundingTime string `json:"next_funding_time"`
}
