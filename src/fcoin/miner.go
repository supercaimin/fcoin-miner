package fcoin

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

const (
	MinerNormalLevel       = 0 // 同样价格买入卖出
	MinerConservatismLevel = 1 // 低买、高卖
	MinerRadicalLevel      = 2 // 高买、低卖
)

type Miner struct {
	wg                 sync.WaitGroup
	Level              int
	symbol             string
	inBalance          float64
	outBalance         float64
	price              float64
	checkStateRetryCnt int
}

func NewMiner() *Miner {
	m := new(Miner)
	m.Level = LEVEL
	m.symbol = TARGET_SYMBOL_IN + TARGET_SYMBOL_OUT
	return m
}

func (m *Miner) buy() {
	m.wg.Done()
}

func (m *Miner) sell() {
	m.wg.Done()
}

// 更新账号余额
func (m *Miner) updateBalance() {
	data, _ := ApiInstance.GetBalance()
	mdata := data.(map[string]interface{})
	currencys := mdata["data"].([]interface{})
	for _, v := range currencys {
		mv := v.(map[string]interface{})
		if mv["currency"].(string) == TARGET_SYMBOL_IN {
			m.inBalance, _ = strconv.ParseFloat(mv["available"].(string), 64)
		}
		if mv["currency"].(string) == TARGET_SYMBOL_OUT {
			m.outBalance, _ = strconv.ParseFloat(mv["available"].(string), 64)
		}
	}
}

// 更新当前成交价格
func (m *Miner) updatePrice() {
	data, _ := ApiInstance.GetTicker(m.symbol)
	mdata := data.(map[string]interface{})
	mmdata := mdata["data"].(map[string]interface{})
	ticker := mmdata["ticker"]
	values := ticker.([]interface{})
	m.price, _ = values[0].(float64)
	fmt.Println("Origin Price:", m.price)
}

func (m *Miner) order(otype string) string {
	var price float64 = m.price
	if otype == "buy" {
		if m.Level == MinerRadicalLevel {
			price = m.price + PRICE_DEFF
		}
		if m.Level == MinerConservatismLevel {
			price = m.price - PRICE_DEFF
		}
		if m.outBalance < PER_ORDER_AMOUNT*price {
			panic(TARGET_SYMBOL_OUT + "余额不足")
		}

	}

	if otype == "sell" {
		if m.Level == MinerRadicalLevel {
			price = m.price - PRICE_DEFF
		}
		if m.Level == MinerConservatismLevel {
			price = m.price + PRICE_DEFF
		}

		if m.inBalance < PER_ORDER_AMOUNT*price {
			panic(TARGET_SYMBOL_IN + "余额不足")
		}
	}
	fmt.Println("下" + otype + "单！")
	fmt.Println("价格:", price)
	data, err := ApiInstance.CreateOrder(m.symbol, otype, "limit", strconv.FormatFloat(price, 'f', 6, 64), strconv.FormatInt(PER_ORDER_AMOUNT, 10), "main")
	if err != nil {
		fmt.Println("订单提交失败:", err)
	}

	mdata := data.(map[string]interface{})

	orderId := mdata["data"].(string)
	fmt.Println("提交订单成功！订单号：", orderId)

	return orderId
}

func (m *Miner) checkOrderState(orderId string) bool {

	time.Sleep(time.Duration(10) * time.Second)

	data, _ := ApiInstance.GetOrder(orderId)
	mdata := data.(map[string]interface{})
	mmdata := mdata["data"].(map[string]interface{})
	state := mmdata["state"].(string)
	if state != "filled" {
		if m.checkStateRetryCnt == 3 {
			ApiInstance.CancelOrder(orderId)
			fmt.Println("订单" + orderId + "未成交，已取消")
			return false
		}
		m.checkStateRetryCnt++
		m.checkOrderState(orderId)
	}
	return true
}

func (m *Miner) goWorker(atype string) {
	m.updateBalance()
	m.updatePrice()
	orderId := m.order(atype)
	m.checkStateRetryCnt = 0
	if m.checkOrderState(orderId) {
		fmt.Println("订单" + orderId + "已成交！")
		if atype == "buy" {
			m.goWorker("sell")
		} else {
			m.goWorker("buy")
		}
	} else {
		m.goWorker(atype)
	}
}

func (m *Miner) Start() {
	m.goWorker("buy")
}

var FCoinMiner *Miner

func init() {
	FCoinMiner = NewMiner()
}
