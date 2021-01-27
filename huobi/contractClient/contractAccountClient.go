package contractClient

import (
	"BitCoinProfitStrategy/huobi/getRequest/model"
	"BitCoinProfitStrategy/huobi/requestBuilder"
	"BitCoinProfitStrategy/huobi/response/contract"
	"BitCoinProfitStrategy/utils"
	"encoding/json"
	"github.com/pkg/errors"
)

type ContractAccountClient struct {
	privateUrl *requestBuilder.PrivateUrlBuilder
}


func (c *ContractAccountClient) Init(accessKey, secretKey, host string) *ContractAccountClient {
	c.privateUrl = new(requestBuilder.PrivateUrlBuilder).Init(accessKey, secretKey, host)
	return c
}


func (c *ContractAccountClient) AccountInfo(request *model.GetContractAccountRequest) (*contract.ContractAccountInfo ,error) {


	url := c.privateUrl.Build("POST", "/api/v1/contract_account_info", nil)

	body, err := utils.ToJson(request)
	if err != nil{
		return nil, err
	}

	resp, err := utils.HttpPost(url, body)

	var info = contract.ContractAccountInfo{}

	jErr := json.Unmarshal([]byte(resp), &info)
	if jErr != nil{
		return nil, jErr
	}

	if info.Status == "ok" && info.Data != nil{
		return &info, nil
	}

	return nil, errors.New(resp)



}
