package fcoin

// API地址，FCOIN为 "https://api.fcoin.com/v2/" FCoinJP为 "https://api.fcoinjp.com/v2/"
const API_URL = "https://api.fcoinjp.com/v2/"

// API KEY SECRET获取方式参考文档https://developer.fcoin.com/zh.html
const API_KEY = "234400c9944443299cbe5c63a756b356"
const API_SECRET = "2705e1b70adb04711a13db6dd765ce50ff"

const WSS_URL = "wss://api.fcoin.com/v2/ws"

// 要挖矿的交易对
const TARGET_SYMBOL_IN = "fj"
const TARGET_SYMBOL_OUT = "usdt"

// 挖矿模式
//const (
//	MinerNormalMode       = 0 // 同样价格买入卖出，无矿损，挖坑效率中性
//	MinerConservatismMode = 1 // 低买、高卖，无矿损，相同本金挖矿最多，但效率低下
//	MinerRadicalMode      = 2 // 高买、低卖，有矿损，挖矿效率最高
//	MinerFastMode         = 3 // 同时下买单、卖单无矿损，挖矿效率极高，推荐模式
//)

const MINER_MODE = MinerFastMode

// 差价0.001522
const PRICE_DEFF = 0.000001

// 单笔订单购买数量
const PER_ORDER_AMOUNT = 10
