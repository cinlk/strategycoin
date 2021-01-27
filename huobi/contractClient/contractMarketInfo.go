package contractClient

import (
	"BitCoinProfitStrategy/huobi/getRequest/model"
	"BitCoinProfitStrategy/huobi/requestBuilder"
	"BitCoinProfitStrategy/huobi/getRequest"
	"BitCoinProfitStrategy/utils"
	"BitCoinProfitStrategy/huobi/response/contract"
	"encoding/json"
	"errors"
	"strconv"
)



type ContractMarketClient struct {

	publicUrlBuilder  *requestBuilder.PublicUrlBuilder

}



func (c *ContractMarketClient) Init(host string) *ContractMarketClient{
	c.publicUrlBuilder = new(requestBuilder.PublicUrlBuilder).Init(host)
	return c
}


// apis

func (c *ContractMarketClient) GetBasis(symbol, period string, size int,
	optionalRequest *model.GetContractBasisRequest )  (*contract.BasisResponse ,error) {

	request := new(getRequest.GetRequest).Init()
	request.AddParam("symbol", symbol)
	request.AddParam("size", strconv.Itoa(size))
	request.AddParam("period", period)

	if optionalRequest != nil && optionalRequest.BasisPriceType != "" {
		request.AddParam("basis_price_type", optionalRequest.BasisPriceType)
	}

	url := c.publicUrlBuilder.Build("/index/market/history/basis", request)

	resp, err := utils.HttpGet(url)
	if err != nil{
		return nil,  err

	}

	result := contract.BasisResponse{}
	jsonErr := json.Unmarshal([]byte(resp), &result)
	if jsonErr != nil{
		return nil, jsonErr
	}
	if result.Status == "ok" && result.Data != nil{
		return &result, nil
	}

	return nil, errors.New(resp)


}




