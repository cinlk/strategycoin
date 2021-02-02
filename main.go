package main

import "BitCoinProfitStrategy/cmd"

func main() {

	if err := cmd.Run(); err != nil {
		panic(err)
	}

}
