package main

import (
	"BitCoinProfitStrategy/huobi/contractClient"
	"fmt"
)

func main(){


	client := new(contractClient.CheckServerClient).Init("api.btcgateway.pro")

	data, err := client.CheckServer()
	if err != nil{
		panic(err)
	}
	fmt.Printf("%+v", data.Data)
}