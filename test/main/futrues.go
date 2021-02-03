package main

import (
	"BitCoinProfitStrategy/huobi/CashTradeClient"
	"BitCoinProfitStrategy/huobi/CurrencyPersistantContract"
	"BitCoinProfitStrategy/huobi/getRequest/model"
	"BitCoinProfitStrategy/huobi/response/CurrencyContract"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FutrueProfitRate struct {
	CashContractCode   string `json:"cash_contract_code"`
	CurrencyFutureCode string `json:"currency_future_code"`

	CashCodePrice float64 `json:"cash_code_price"`
	FutruePrice   float64 `json:"futrue_price"`

	FundingRate     float64 `json:"funding_rate"`
	NextFundingRate float64 `json:"next_funding_rate"`

	// 费率结算时间
	FundingTime string `json:"funding_time"`

	// 1个币对应合约张数
	Sheets float64 `json:"sheets"`
}

// 计算没张，不同币的张对应usd个数不一样 差价利润 TODO

// 合约考虑张，合约赚的该币价差， 现货转usdt价差，合约的价差换算成usdt，计算最终价差
type ProfitResult struct {
	ContractCode string `json:"contract_code"`

	// 操作方向
	Operation string `json:"operation"`

	// 现货价格
	CashPrice float64 `json:"cash_price"`
	// 合约价格
	FuturePrice float64 `json:"future_price"`

	// 当前费率
	FundingRate float64 `json:"funding_rate"`

	//  买卖单考虑立刻成交，只是taker 手续费 TODO
	BuyCashCoinRate float64 `json:"buy_cash_coin_rate"`

	SellCashCoinRate float64 `json:"sell_cash_coin_rate"`

	// 费率结算时间
	FundingTime string `json:"funding_time"`

	// 默认1倍
	BuyContractRate  float64 `json:"buy_contract_rate"`
	SellContractRate float64 `json:"sell_contract_rate"`

	// 预估 平仓时合约和现货价格， 计算最终得到的值有价差利润  TODO
	CloseCash     float64 `json:"close_cash"`
	CloseContract float64 `json:"close_contract"`

	// 预估得到的利润
	Profit float64 `json:"profit"`
}

// 现货 合约差价
type FutureGoodsSpread struct {
	Operation    string  `json:"operation"`
	Margin       float64 `json:"margin"`
	ContractCode string  `json:"contract_code"`

	GoodsPrice  float64 `json:"goods_price"`
	FutruePrice float64 `json:"futrue_price"`

	FundingRate     float64 `json:"funding_rate"`
	NextFundingRate float64 `json:"next_funding_rate"`
}

type allSpreadList []FutureGoodsSpread

func (a allSpreadList) Len() int {
	return len(a)
}
func (a allSpreadList) Less(i, j int) bool {
	return a[i].Margin < a[j].Margin
}

func (a allSpreadList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// test cash client

func testCashClient() {

	client := new(CashTradeClient.CashMarketInfo).Init("api.huobi.pro")

	data, err := client.MarcketTrade("btcussdt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", data.Tick)

}

func testAllContract() {

	client := new(CurrencyPersistantContract.ContractMarketInfo).Init("api.btcgateway.pro")
	data, err := client.ContractInfo(nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", data.Data)
	var codes []string
	for _, d := range data.Data {
		codes = append(codes, d.ContractCode)
	}
	var wait sync.WaitGroup

	for _, c := range codes {

		wait.Add(1)
		go func(code string) {
			defer func() {
				wait.Done()
			}()

			data, err := client.CurrentFundingRate(code)
			if err != nil {
				panic(err)
			}
			fmt.Printf("\n%+v", data.Data)
		}(c)
	}

	wait.Wait()

}

// 必须输入当前买入现货价格，买入和合约价格,  以当前的合约价格作为交割结算计算利润， 不等到交割时间计算利润
// 以当前价格买入现货和合约，假设交割价格一样，计算交割利润 TODO

func ComputeSpread(buyContractRate float64, sellContractRate float64,
	buyGoodRate float64, sellGoodRate float64, coinIntrestDay float64,
	buyContractPrice float64, buyGoodPrice float64, coinCount int64) []ProfitResult {

	var res []ProfitResult

	// 带入 费率计算差价（预估）  给出卖出平仓价位
	source := ContractCodePairePrices(allContractCodeAndRate())
	for _, s := range source {
		tmp := ProfitResult{}

		if s.FundingRate > 0 {
			// 买入现货， 开空合约
			tmp.Operation = "买入现货-开空1倍合约"
			tmp.BuyCashCoinRate = buyGoodRate
			tmp.BuyContractRate = buyContractRate

			//  计算profit
			// 参考火币 https://support.hbfile.net/hc/zh-cn/articles/900000106903-%E8%B5%84%E9%87%91%E8%B4%B9%E7%94%A8%E8%AF%B4%E6%98%8E
			// 假设净持仓量为负的 TODO

			// 做空 会有币 兑换usd 减值风险

			// 需要搞明白
			// 1 taker 和maker  maker 挂单，taker 主动吃单
			//
			// 2 买卖现货 得到的利润， 计算方式
			// 利润换算成u， (卖出总价-买入总价)- 买入总价*手续费率 - 卖出总价*手续费
			// 3 买卖合约  得到的利润， 计算方式，
			//  合约保证金 担保金计算  持仓担保资金=（合约面值*持仓合约数量）/ 最新成交价 / 倍数
			//  可用担保资产 怎么计算？
			//  永续账户权益 = 账户余额 + 本期已实现盈亏 + 本期未实现盈亏
			//  账户余额 不一定等于100%成交的委托量，成交量 会除去可用担保资产和手续费
			//  已实现盈亏：手续费等  待补充---
			//  未实现盈亏计算：
			/*

				多仓未实现盈亏 =（1/持仓均价-1/最新成交价）* 多仓合约张数 * 合约面值

				空仓未实现盈亏 =（1/最新成交价-1/持仓均价）* 空仓合约张数 * 合约面值
				// 平仓盈亏
				多仓已实现盈亏 =（1/持仓均价-1/平仓成交均价）* 平多仓合约张数 * 合约面值

				空仓已实现盈亏 =（1/平仓成交均价-1/持仓均价）* 平空仓合约张数 * 合约面值

			*/
			// 4 合约费率计算方式 得到的利润
			// 资金费率 净持仓量 * 合约面值 / 结算价 * 资金费率；
			// 结算价格 结算价格 = sum ( 每笔成交量(张) ) / sum ( 每笔成交量(张) / 每笔成交价 )
			// 合约保证币 随市场价格的变动

			//var profit  = tmp.CashPrice*(1-sellGoodRate) +  tmp.FuturePrice*(1-sellContractRate) +
			//	(buyContractPrice *s.Sheets/s.FutruePrice)*s.FundingRate - (buyGoodPrice*(1+buyGoodRate) +
			//		buyContractPrice*(1+buyContractRate))
			//

			// 1 计算现货 利润

		} else {
			// 卖出现货,  开多合约 (借币卖出，考虑利息)
			tmp.Operation = "卖出现货-开多1倍合约"
		}

		tmp.FundingRate = s.FundingRate
		tmp.ContractCode = s.CurrencyFutureCode
		tmp.CashPrice = s.CashCodePrice
		tmp.FuturePrice = s.FutruePrice
		tmp.FundingTime = s.FundingTime

		res = append(res, tmp)
	}

	return res

}

// 4 现货买入 卖出利润，以usdt为本位
func CointTradeProfit(coinCount float64, buyUsdtPrice float64, sellUsdtPrice float64, rate float64) float64 {

	return math.Pow(1-rate, 2)*coinCount*sellUsdtPrice - coinCount*buyUsdtPrice
}

// 接币 卖出

// 5 合约做空利润, 以币为单位, 1倍做空，得到该币的利润， 币换成u，用现货的值计算
func CurrencyContractShortProfit(danbao float64, buyContractPrice, sellContractPrice, coinCount, buyRate, sellRate float64) float64 {

	//return (1-sellRate)*((1/sellContractPrice - 1/buyContractPrice)*(coinCount*(1-buyRate)-danbao)*buyContractPrice+
	//	(coinCount*(1-buyRate)-danbao)) - (coinCount*(1-buyRate)-danbao)

	return (coinCount*(1-buyRate) - danbao) * (((1/sellContractPrice-1/buyContractPrice)*buyContractPrice+1)*(1-sellRate) - 1)

}

// 6 合约做多利润

func CurrencyContractBuyProfit(danbao float64, buyContractPrice, closeContractPrice, coinCount, buyRate, sellRate float64) float64 {

	return (coinCount*(1-buyRate) - danbao) * (((1/buyContractPrice-1/closeContractPrice)*buyContractPrice+1)*(1-sellRate) - 1)
}

// 7 费率利润
func ContracteMoneyRate(coinCount float64, buyRate float64, danbao float64,
	buyContractPrice float64, ClosePrice float64, rate float64) float64 {

	return (coinCount*(1-buyRate) - danbao) * buyContractPrice / ClosePrice * rate

}

// 计算费率 差价
func demotest() {

	// 现货 买卖利润
	c1 := CointTradeProfit(100, 101, 101.50, 0.002)

	// 合约做空 利润
	c2 := CurrencyContractShortProfit(0.04, 102, 101.60,
		100, 0.0005, 0.0005)

	// 合约做多 利润
	c3 := CurrencyContractBuyProfit(0, 97.70, 99.50, 65.139,
		0.0005, 0.0005)

	// 扩大倍数
	r := ContracteMoneyRate(120.3, 0.0005, 0.04, 100.60, 100.500, 0.0016)

	// 杠杆
	fmt.Printf("%f %f %f %f", c1, c2*102.50, c3, r)

}

// 1 获取所有可用币本位代码， 以及费率
func allContractCodeAndRate() []CurrencyContract.LatestFundingRatedata {

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

// 2 获取所有交易对的 币币交易价格  和  币本位永续当前价格  （每秒刷新）

func ContractCodePairePrices(codes []CurrencyContract.LatestFundingRatedata) []FutrueProfitRate {

	var res []FutrueProfitRate

	if len(codes) == 0 {
		return nil
	}
	client := new(CurrencyPersistantContract.ContractMarketInfo).Init("api.btcgateway.pro")
	for _, c := range codes {
		var tmp FutrueProfitRate

		trade, err := client.LatestTradeRecord(c.ContractCode)
		if err != nil {
			fmt.Printf("%s --->", err)
			continue
		}

		if len(trade.Tick.Data) >= 1 {
			// 取第一个数据 为合约价格
			cPrice, err := strconv.ParseFloat(trade.Tick.Data[0].Price, 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			tmp.FutruePrice = cPrice
			tmp.CurrencyFutureCode = strings.ToLower(c.ContractCode)
		} else {
			continue
		}

		// 币币交易 code 名字改变
		goodClient := new(CashTradeClient.CashMarketInfo).Init("api.huobi.pro")
		cashCode := change2CoinPaire(c.ContractCode)

		data, err := goodClient.MarcketTrade(cashCode)
		if err != nil {
			fmt.Printf("%s <-----", err)
			continue
		}
		if len(data.Tick.Data) >= 1 {
			// 取第一个数据 为现货交易数据
			gPrice := data.Tick.Data[0].Price
			tmp.CashCodePrice = gPrice
			tmp.CashContractCode = cashCode

		} else {
			continue
		}

		frate, err := strconv.ParseFloat(c.FundingRate, 64)
		if err != nil {
			fmt.Println(err)
			continue
		}

		erate, err := strconv.ParseFloat(c.EstimatedRate, 64)

		if err != nil {
			fmt.Println(err)
			continue
		}

		tmp.FundingRate = frate
		tmp.NextFundingRate = erate
		tmp.FundingTime, _ = msToTime(c.FundingTime)
		res = append(res, tmp)

	}

	return res

}

// 时间戳转换为 可读时间

func change2CoinPaire(code string) string {
	c := strings.ToLower(code)
	m := strings.Split(c, "-")
	m[1] += "t"
	return strings.Join(m, "")
}

func testMarket() {

	client := new(CurrencyPersistantContract.ContractMarketInfo).Init("api.btcgateway.pro")

	data, err := client.CurrentFundingRate("eos-usdt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", data.Data)

}

func testGetCurrencyIndexs() {

	client := new(CurrencyPersistantContract.ContractMarketInfo).Init("api.btcgateway.pro")

	var req = model.GetContractCodeRequest{
		ContractCode: "storj-usd",
	}
	data, err := client.ContractIndex(&req)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", data.Data)

	d, err := client.LatestTradeRecord("storj-usd")
	fmt.Printf("%+v", d.Tick)
}

func msToTime(ms string) (string, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return "", err
	}

	tm := time.Unix(0, msInt*int64(time.Millisecond))

	return tm.Format("2006-02-01 15:04:05.000"), nil

}

func main() {

	//testMarket()
	//testAllContract()
	//testGetCurrencyIndexs()

	//  速度优化，改成http 接口



	fmt.Printf("%v", allContractCodeAndRate())
	//demotest()

}
