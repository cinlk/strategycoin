package futureProfit



// 现货价格
type Goods struct {
	Direction  int  // -1 卖出，  1 买入
	Price  float64   // 现货开仓的价格（买入时的价格，卖出时的价格）
	Coin string  // 币种 交易对
}

// 合约价格
type Contract struct {
	Direction int // -1 卖出开空   1 买入开多
	Price float64 // 合约开仓的价格
	Coin  string // 币种  （币本位永续  usdt本位永续）
	Leverage int // 杠杠倍数
	Rate float64 // 当前费率值
	NRate float64 // 下期费率值
	times int  // 每天3次结算

}





var marketContractPrice float64 // 当前市场该币种合约价格

var marketGoodProice float64 // 当前市场该币种的现货价格

var goodFee float64 // 开仓现货的手续费

var contractFee float64 // 开仓合约的手续费




func fetchData() {

	//client := client.MarketClient{}
}


// 按照当前资金费率计算的差价
func caculateProfit() {


	// 获取 币本位 当前合约价格
	// 买卖手续费率的计算
	//  如何预估利润 TODO ？



	// 计算不同币以及合约 价格

	// 准备计算数据

	// 计算利润



}


func watchProfit(goods float64,  contract float64){

	// 每秒刷新


}



