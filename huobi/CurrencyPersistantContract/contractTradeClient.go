package CurrencyPersistantContract

import (
	"BitCoinProfitStrategy/huobi/getRequest/model"
	"BitCoinProfitStrategy/huobi/requestBuilder"
	"BitCoinProfitStrategy/huobi/response/CurrencyContract"
	"BitCoinProfitStrategy/utils"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type ContractTradeClient struct {
	privateBuildUrl *requestBuilder.PrivateUrlBuilder
}


func (c *ContractTradeClient) Init(host string, accessKey string, secretKey string) *ContractTradeClient{
	c.privateBuildUrl = new(requestBuilder.PrivateUrlBuilder).Init(accessKey, secretKey, host)
	return  c

}



// 设置杠杆倍数
func (c *ContractTradeClient) SwapLeverRate(contractCode string, leverRate int)(*CurrencyContract.LeverRate, error){

	body := make(map[string]interface{})
	body["contract_code"] = contractCode
	body["lever_rate"] = leverRate


	url := c.privateBuildUrl.Build("POST", "/swap-api/v1/swap_switch_lever_rate", nil)

	d, _ := json.Marshal(body)
	resp, err := utils.HttpPost(url,string(d))
	if err != nil{
		return nil, err
	}
	var res  CurrencyContract.LeverRate
	jErr := json.Unmarshal([]byte(resp), &res)
	if jErr != nil{

		return nil, errors.Wrap(jErr, "json decode failed")
	}

	return &res, nil


}

// 撤销订单
func (c *ContractTradeClient) CancelTradeOrder(contractCode string, optionRequest *model.CancerOrderReq)(*CurrencyContract.CancerContractOrderResp, error ){

	body := make(map[string]interface{})

	body["contract_code"] = contractCode
	if optionRequest != nil{
		body["order_id"] = optionRequest.OrderId
		body["client_order_id"] = optionRequest.ClientOrderId
	}

	url := c.privateBuildUrl.Build("POST", "/swap-api/v1/swap_cancel", nil)

	d, _ := json.Marshal(body)
	resp, err := utils.HttpPost(url, string(d))
	if err != nil{
		return  nil, err
	}

	var  res CurrencyContract.CancerContractOrderResp

	jErr := json.Unmarshal([]byte(resp), &res)
	if jErr != nil{
		return  nil, errors.Wrap(jErr, "json decode failed")
	}

	return  &res, nil
}


// 查询订单
func (c *ContractTradeClient)GetOrderInfo (contractCode string, optionRequest *model.CancerOrderReq)(
	*CurrencyContract.OrderInfoResp, error){

	body := make(map[string]interface{})
	body["contract_code"] = contractCode
	if optionRequest != nil{
		body["order_id"] = optionRequest.OrderId
		body["client_order_id"] = optionRequest.ClientOrderId
	}


	b, _ := json.Marshal(body)

	url := c.privateBuildUrl.Build("POST","/swap-api/v1/swap_order_info", nil)

	resp, err := utils.HttpPost(url, string(b))
	if err != nil{
		return nil, errors.Wrap(err, "post http failed")
	}

	var res CurrencyContract.OrderInfoResp

	// status 为error 处理 TODO
	//fmt.Printf("%s", resp)
	jErr := json.Unmarshal([]byte(resp), &res)
	if jErr != nil{
		return  nil, errors.Wrap(err, "json decode failed")
	}

	return &res, nil

}


// 合约下单

func (c *ContractTradeClient) CreateOder(req *model.CreateOrderReq)(*CurrencyContract.NewOrderResq, error){


	b, _ := json.Marshal(req)

	url := c.privateBuildUrl.Build("POST","/swap-api/v1/swap_order", nil)

	resp, err := utils.HttpPost(url, string(b))
	if err != nil{
		return nil, errors.Wrap(err, "post http failed")
	}

	var res CurrencyContract.NewOrderResq


	jErr := json.Unmarshal([]byte(resp), &res)
	if jErr != nil{
		return  nil, errors.Wrap(err, "json decode failed")
	}
	fmt.Printf("%s", resp)
	return &res, nil
}

