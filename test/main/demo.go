package main

import (
	"BitCoinProfitStrategy/huobi/CurrencyPersistantContract"
	"BitCoinProfitStrategy/huobi/response/CurrencyContract"
	"fmt"
	"sort"
	"strconv"
)

func main() {

	// 资金费率排名

	fmt.Printf("%v", profit())
}

func profit() []CurrencyContract.LatestFundingRatedata {

	var res []CurrencyContract.LatestFundingRatedata
	client := new(CurrencyPersistantContract.ContractMarketInfo).Init("api.btcgateway.pro")
	contracts, err := client.ContractInfo(nil)
	if err != nil {
		panic(err)
	}

	for _, c := range contracts.Data {
		if c.ContractStatus == 1 || c.ContractStatus == 5 || c.ContractStatus == 7 {
			// 获取费率  次数多，频率限制
			data, err := client.CurrentFundingRate(c.ContractCode)
			if err != nil {
				fmt.Printf("%s ---->", err)
				continue
			}
			res = append(res, *data.Data)

		}
	}

	sort.Slice(res, func(i, j int) bool {
		r1, err := strconv.ParseFloat(res[i].FundingRate, 64)
		if err != nil {
			panic(err)
		}

		r2, err := strconv.ParseFloat(res[j].FundingRate, 64)
		if err != nil {
			panic(err)
		}
		return r1 >= r2

	})
	return res

}
