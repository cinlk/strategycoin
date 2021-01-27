package CashTradeClient

import (
	"BitCoinProfitStrategy/huobi/requestBuilder"
	"BitCoinProfitStrategy/huobi/response/Cash"
	"BitCoinProfitStrategy/utils"
	"encoding/json"
	"github.com/pkg/errors"
)
import "BitCoinProfitStrategy/huobi/getRequest"


type CashMarketInfo struct {
	publicUrlBuilder *requestBuilder.PublicUrlBuilder
}


func (c *CashMarketInfo) Init(host string) *CashMarketInfo {
	c.publicUrlBuilder = new(requestBuilder.PublicUrlBuilder).Init(host)
	return c
}


func (c *CashMarketInfo) MarcketTrade(symbol string) ( *Cash.LatestTradeRecord , error) {

	req := new(getRequest.GetRequest).Init()
	req.AddParam("symbol", symbol)
	url := c.publicUrlBuilder.Build("/market/trade", req)

	resp, err := utils.HttpGet(url)
	if err != nil{
		return  nil, err
	}
	var data = Cash.LatestTradeRecord{}
	jErr := json.Unmarshal([]byte(resp), &data)

	if jErr != nil{
		return nil, jErr
	}

	if data.Status == "ok" && data.Tick != nil{
		return &data, nil
	}

	return nil, errors.New(resp)
}
