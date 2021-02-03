package handlers

import (
	"BitCoinProfitStrategy/huobi/CashTradeClient"
	"BitCoinProfitStrategy/huobi/CurrencyPersistantContract"
	"BitCoinProfitStrategy/huobi/response/CurrencyContract"
	"fmt"
	"github.com/gin-gonic/gin"
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
	Operation    string  `json:"操作方式"`
	Margin       float64 `json:"差价比率"`
	ContractCode string  `json:"现货名称"`

	GoodsPrice  float64 `json:"现货价格"`
	FutruePrice float64 `json:"合约价格"`

	FundingRate     float64 `json:"合约资金费率"`
	NextFundingRate float64 `json:"下次合约资金费率"`
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



func FgProfit(c *gin.Context) {

	source := ContractCodePairePrices(allContractCodeAndRate())

	var datas allSpreadList

	for _, s := range source {
		var tmp FutureGoodsSpread

		if s.FutruePrice >= s.CashCodePrice {

			c := (s.FutruePrice - s.CashCodePrice) / s.CashCodePrice
			//fmt.Printf("%s 差价幅度%f, 做空合约 买入现货", s.CashContractCode, c)
			tmp.Margin = c
			tmp.Operation = "合约空现货买"

		} else {
			c := (s.CashCodePrice - s.FutruePrice) / s.FutruePrice
			//fmt.Printf( "%s 差价幅度%f, 做多合约   卖出现货（借币）",s.CashContractCode, c)
			tmp.Margin = c
			tmp.Operation = "合约多现货卖"

		}

		tmp.FutruePrice = s.FutruePrice
		tmp.ContractCode = s.CashContractCode
		tmp.GoodsPrice = s.CashCodePrice
		tmp.FundingRate = s.FundingRate
		tmp.NextFundingRate = s.NextFundingRate
		datas = append(datas, tmp)

	}

	sort.Sort(sort.Reverse(datas))

	c.JSON(200, datas)

}



func allContractCodeAndRate() []CurrencyContract.LatestFundingRatedata {

	var res []CurrencyContract.LatestFundingRatedata
	client := new(CurrencyPersistantContract.ContractMarketInfo).Init("api.btcgateway.pro")
	contracts, err := client.ContractInfo(nil)
	if err != nil {
		panic(err)
	}



	var wait sync.WaitGroup


	for _, c := range contracts.Data {

		tmp := c
		if c.ContractStatus == 1 || c.ContractStatus == 5 || c.ContractStatus == 7 {
			// 获取费率  次数多，频率限制
			wait.Add(1)
			go func() {
				defer func() {
					wait.Done()
				}()
				data, err := client.CurrentFundingRate(tmp.ContractCode)
				if err != nil {
					fmt.Printf("%s ---->", err)
					return
				}
				res = append(res, *data.Data)
			}()
		}
	}

	wait.Wait()

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


func ContractCodePairePrices(codes []CurrencyContract.LatestFundingRatedata) []FutrueProfitRate {

	var res []FutrueProfitRate

	if len(codes) == 0 {
		return nil
	}
	client := new(CurrencyPersistantContract.ContractMarketInfo).Init("api.btcgateway.pro")

	var wait sync.WaitGroup

	for _, c := range codes {
		wait.Add(1)

		go func(m *CurrencyContract.LatestFundingRatedata) {

			defer func() {
				wait.Done()
			}()
			var tmp FutrueProfitRate


			trade, err := client.LatestTradeRecord(c.ContractCode)
			if err != nil {
				fmt.Printf("%s --->", err)
				return
			}

			if len(trade.Tick.Data) >= 1 {
				// 取第一个数据 为合约价格
				cPrice, err := strconv.ParseFloat(trade.Tick.Data[0].Price, 64)
				if err != nil {
					fmt.Println(err)
					return
				}
				tmp.FutruePrice = cPrice
				tmp.CurrencyFutureCode = strings.ToLower(c.ContractCode)
			} else {
				return
			}

			// 币币交易 code 名字改变
			goodClient := new(CashTradeClient.CashMarketInfo).Init("api.huobi.pro")
			cashCode := change2CoinPaire(c.ContractCode)

			data, err := goodClient.MarcketTrade(cashCode)
			if err != nil {
				fmt.Printf("%s <-----", err)
				return
			}
			if len(data.Tick.Data) >= 1 {
				// 取第一个数据 为现货交易数据
				gPrice := data.Tick.Data[0].Price
				tmp.CashCodePrice = gPrice
				tmp.CashContractCode = cashCode

			} else {
				return
			}

			frate, err := strconv.ParseFloat(c.FundingRate, 64)
			if err != nil {
				fmt.Println(err)
				return
			}

			erate, err := strconv.ParseFloat(c.EstimatedRate, 64)

			if err != nil {
				fmt.Println(err)
				return
			}

			tmp.FundingRate = frate
			tmp.NextFundingRate = erate
			tmp.FundingTime, _ = msToTime(c.FundingTime)
			res = append(res, tmp)

		}(&c)
	}
	return res

}


func change2CoinPaire(code string) string {
	c := strings.ToLower(code)
	m := strings.Split(c, "-")
	m[1] += "t"
	return strings.Join(m, "")
}


func msToTime(ms string) (string, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return "", err
	}

	tm := time.Unix(0, msInt*int64(time.Millisecond))

	return tm.Format("2006-02-01 15:04:05.000"), nil

}
