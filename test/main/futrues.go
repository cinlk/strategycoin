package main

import (
	"BitCoinProfitStrategy/huobi/CashTradeClient"
	"BitCoinProfitStrategy/huobi/contractClient"
	"fmt"
)



// test cash client

func testCashClient() {

	client := new(CashTradeClient.CashMarketInfo).Init("api.huobi.pro")

	data, err := client.MarcketTrade("btcussdt")
	if err != nil{
		panic(err)
	}
	fmt.Printf("%+v", data.Tick)


}

func testMarket(){
	client := new(contractClient.CheckServerClient).Init("api.btcgateway.pro")

	data, err := client.CheckServer()
	if err != nil{
		panic(err)
	}
	fmt.Printf("%+v", data.Data)
}

func main(){

	testCashClient()

}

