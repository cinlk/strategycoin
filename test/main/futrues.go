package main

import (
	"BitCoinProfitStrategy/huobi/CashTradeClient"
	"BitCoinProfitStrategy/huobi/CurrencyPersistantContract"
	"BitCoinProfitStrategy/huobi/getRequest/model"
	"BitCoinProfitStrategy/huobi/response/CurrencyContract"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)



type FutrueProfitRate struct {
	CashContractCode string `json:"cash_contract_code"`
	CurrencyFutureCode string `json:"currency_future_code"`
	
	CashCodePrice float64 `json:"cash_code_price"`
	FutruePrice  float64 `json:"futrue_price"`

	FundingRate  float64 `json:"funding_rate"`
	NextFundingRate float64 `json:"next_funding_rate"`
	
	
	
}


// test cash client

func testCashClient() {

	client := new(CashTradeClient.CashMarketInfo).Init("api.huobi.pro")

	data, err := client.MarcketTrade("btcussdt")
	if err != nil{
		panic(err)
	}
	fmt.Printf("%+v", data.Tick)


}

func testAllContract(){

	client := new(CurrencyPersistantContract.ContractMarketInfo).Init("api.btcgateway.pro")
	data ,err := client.ContractInfo(nil)
	if err != nil{
		panic(err)
	}
	fmt.Printf("%+v\n", data.Data)
	var codes []string
	for _, d := range data.Data{
		 codes = append(codes, d.ContractCode)
	}
	var wait sync.WaitGroup

	for _, c := range codes{

		wait.Add(1)
		go func(code string) {
			defer func() {
				wait.Done()
			}()

			data, err := client.CurrentFundingRate(code)
			if err != nil{
				panic(err)
			}
			fmt.Printf("\n%+v", data.Data)
		}(c)
	}

	wait.Wait()


}



func ComputeSpread(MakerfutureRate float64, TakeFutrueRate float64,
	MakerGoodRate float64, TakeGoodRate float64,  CoinIntrestDay float64){

		//1  获取所有可用币本位代码， 以及费率，安装费率从大到小排序


		// 获取所有交易对的 币币交易价格  和  币本位永续价格


		// 带入 费率计算差价（预估）  给出卖出平仓价位



}


// 1 获取所有可用币本位代码， 以及费率
func allContractCodeAndRate() []CurrencyContract.LatestFundingRatedata {

	var res  []CurrencyContract.LatestFundingRatedata
	client := new(CurrencyPersistantContract.ContractMarketInfo).Init("api.btcgateway.pro")
	contracts, err := client.ContractInfo(nil)
	if err != nil{
		panic(err)
	}

	for _, c := range contracts.Data{
		if c.ContractStatus == 1 || c.ContractStatus == 5 || c.ContractStatus == 7{
			// 获取费率  次数多，频率限制
			data, err := client.CurrentFundingRate(c.ContractCode)
			if err !=nil{
				fmt.Printf("%s ---->",err)
				continue
			}
			res = append(res,  *data.Data)

		}
	}


	sort.Slice(res, func(i, j int) bool {
		r1, err  := strconv.ParseFloat(res[i].FundingRate,64)
		if err != nil{
			panic(err)
		}

		r2, err  := strconv.ParseFloat(res[j].FundingRate,64)
		if err != nil{
			panic(err)
		}
		return r1 >= r2


	})
	return res
}



// 2 获取所有交易对的 币币交易价格  和  币本位永续当前价格  （每秒刷新）

func ContractCodePairePrices(codes []CurrencyContract.LatestFundingRatedata) []FutrueProfitRate{

	var res []FutrueProfitRate

	if len(codes) == 0{
		return nil
	}
	client := new(CurrencyPersistantContract.ContractMarketInfo).Init("api.btcgateway.pro")
	for _, c := range  codes{
		var tmp FutrueProfitRate

		trade, err := client.LatestTradeRecord(c.ContractCode)
		if err != nil{
			fmt.Printf("%s --->", err)
			continue
		}

		if len(trade.Tick.Data) >= 1 {
			// 取第一个数据 为合约价格
			cPrice, err :=  strconv.ParseFloat(trade.Tick.Data[0].Price,64)
			if err != nil{
				fmt.Println(err)
				continue
			}
			tmp.FutruePrice = cPrice
			tmp.CurrencyFutureCode = c.ContractCode
		}else{
			continue
		}

		// 币币交易 code 名字改变
		goodClient := new(CashTradeClient.CashMarketInfo).Init("api.huobi.pro")
		cashCode := change2CoinPaire(c.ContractCode)
		data, err := goodClient.MarcketTrade(cashCode)
		if err != nil{
			fmt.Printf("%s <-----",err)
			continue
		}
		if len(data.Tick.Data) >= 1 {
			// 取第一个数据 为现货交易数据
			gPrice :=  data.Tick.Data[0].Price
			fmt.Printf(" %f", gPrice)
			tmp.CashCodePrice = gPrice
			tmp.CashContractCode = cashCode

		}else {
			continue
		}

		frate, err := strconv.ParseFloat(c.FundingRate,64)
		if err != nil{
			fmt.Println(err)
			continue
		}

		erate, err := strconv.ParseFloat(c.EstimatedRate, 64)

		if err != nil{
			fmt.Println(err)
			continue
		}

		tmp.FundingRate = frate
		tmp.NextFundingRate = erate
		res = append(res, tmp)

	}


	return res


}

func change2CoinPaire(code string) string{
	c := strings.ToLower(code)
	m :=strings.Split(c, "-")
	m[1] += "d"
	return strings.Join(m,"")
}


//  根据费率值，和 手续费预估 合理的差价
func estimateRate(){
	//rate :=  allContractCodeAndRate()


}



func testMarket(){

	client := new(CurrencyPersistantContract.ContractMarketInfo).Init("api.btcgateway.pro")

	data, err := client.CurrentFundingRate("eos-usdt")
	if err != nil{
		panic(err)
	}
	fmt.Printf("%+v", data.Data)

}

func testGetCurrencyIndexs(){

	client := new(CurrencyPersistantContract.ContractMarketInfo).Init("api.btcgateway.pro")

	var req = model.GetContractCodeRequest{
		ContractCode: "storj-usd",
	}
	data, err := client.ContractIndex(&req)
	if err != nil{
		panic(err)
	}
	fmt.Printf("%+v\n", data.Data)


	d, err :=  client.LatestTradeRecord("storj-usd")
	fmt.Printf("%+v", d.Tick)
}


func main(){

	//testMarket()
	//testAllContract()
	//testGetCurrencyIndexs()

	fmt.Printf("%+v", ContractCodePairePrices(allContractCodeAndRate()))

}

