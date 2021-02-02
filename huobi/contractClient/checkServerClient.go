package contractClient

import (
	"BitCoinProfitStrategy/huobi/requestBuilder"
	"BitCoinProfitStrategy/huobi/response/contract"
	"BitCoinProfitStrategy/utils"
	"encoding/json"
	"github.com/pkg/errors"
)

type CheckServerClient struct {
	PublicUrlBuild *requestBuilder.PublicUrlBuilder
}

func (c *CheckServerClient) Init(host string) *CheckServerClient {
	c.PublicUrlBuild = new(requestBuilder.PublicUrlBuilder).Init(host)
	return c
}

func (c *CheckServerClient) CheckServer() (*contract.ServerHeatBeatInfo, error) {

	url := c.PublicUrlBuild.Build("/heartbeat/", nil)

	resp, err := utils.HttpGet(url)
	if err != nil {
		return nil, errors.Wrap(err, "http get data faile")
	}

	var info = contract.ServerHeatBeatInfo{}
	jErr := json.Unmarshal([]byte(resp), &info)
	if jErr != nil {
		return nil, errors.Wrap(jErr, "decode json failed")
	}
	if info.Status == "ok" && info.Data != nil {
		return &info, nil
	}

	return nil, errors.New(resp)
}
