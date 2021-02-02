package CurrencyPersistantContract

import (
	"BitCoinProfitStrategy/huobi/getRequest"
	"BitCoinProfitStrategy/huobi/getRequest/model"
	"BitCoinProfitStrategy/huobi/requestBuilder"
	"BitCoinProfitStrategy/huobi/response/CurrencyContract"
	"BitCoinProfitStrategy/utils"
	"encoding/json"
	"github.com/pkg/errors"
)

type ContractMarketInfo struct {
	publicUrlBuilder *requestBuilder.PublicUrlBuilder
}

func (c *ContractMarketInfo) Init(host string) *ContractMarketInfo {
	c.publicUrlBuilder = new(requestBuilder.PublicUrlBuilder).Init(host)
	return c
}

func (c *ContractMarketInfo) LatestTradeRecord(contractCode string) (*CurrencyContract.LatestTradeRecord, error) {

	req := new(getRequest.GetRequest).Init()
	req.AddParam("contract_code", contractCode)
	url := c.publicUrlBuilder.Build("/swap-ex/market/trade", req)

	resp, err := utils.HttpGet(url)
	if err != nil {
		return nil, err
	}

	var data = CurrencyContract.LatestTradeRecord{}
	jErr := json.Unmarshal([]byte(resp), &data)

	if jErr != nil {
		return nil, jErr
	}

	if data.Status == "ok" && data.Tick != nil {
		return &data, nil
	}

	return nil, errors.New(resp)

}

func (c *ContractMarketInfo) CurrentFundingRate(contractCode string) (*CurrencyContract.FundingRateInfo, error) {

	req := new(getRequest.GetRequest).Init()
	req.AddParam("contract_code", contractCode)

	url := c.publicUrlBuilder.Build("/swap-api/v1/swap_funding_rate", req)

	resp, err := utils.HttpGet(url)
	if err != nil {
		return nil, err
	}

	var data = CurrencyContract.FundingRateInfo{}
	jErr := json.Unmarshal([]byte(resp), &data)

	if jErr != nil {
		return nil, jErr
	}

	if data.Status == "ok" && data.Data != nil {
		return &data, nil
	}

	return nil, errors.New(resp)

}

func (c *ContractMarketInfo) ContractInfo(optionReq *model.GetContractCodeRequest) (
	*CurrencyContract.ContractInfo, error) {

	var req *getRequest.GetRequest = nil
	if optionReq != nil && optionReq.ContractCode != "" {
		req = new(getRequest.GetRequest).Init()
		req.AddParam("contract_code", optionReq.ContractCode)
	}

	url := c.publicUrlBuilder.Build("/swap-api/v1/swap_contract_info", req)
	resp, err := utils.HttpGet(url)

	if err != nil {
		return nil, err
	}

	var data = CurrencyContract.ContractInfo{}
	jErr := json.Unmarshal([]byte(resp), &data)

	if jErr != nil {
		return nil, jErr
	}

	if data.Status == "ok" && data.Data != nil {
		return &data, nil
	}

	return nil, errors.New(resp)

}

func (c *ContractMarketInfo) ContractIndex(optionReq *model.GetContractCodeRequest) (
	*CurrencyContract.ContractIndexInfo, error) {

	var req *getRequest.GetRequest = nil
	if optionReq != nil && optionReq.ContractCode != "" {
		req = new(getRequest.GetRequest).Init()
		req.AddParam("contract_code", optionReq.ContractCode)
	}
	url := c.publicUrlBuilder.Build("/swap-api/v1/swap_index", req)
	resp, err := utils.HttpGet(url)
	if err != nil {
		return nil, errors.Wrap(err, "http get data failed")
	}
	var data = CurrencyContract.ContractIndexInfo{}
	jErr := json.Unmarshal([]byte(resp), &data)
	if jErr != nil {
		return nil, errors.Wrap(jErr, "decode data failed")
	}
	if data.Status == "ok" && data.Data != nil {
		return &data, nil
	}

	return nil, errors.New(resp)

}
