package fcoin

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

const (
	MinerNormalMode       = 0 // 同样价格买入卖出
	MinerConservatismMode = 1 // 低买、高卖
	MinerRadicalMode      = 2 // 高买、低卖
	MinerFastMode         = 3 // 同时小买单、卖单无矿损

)

type Miner struct {
	wg                 sync.WaitGroup
	Mode               int
	symbol             string
	inBalance          float64
	outBalance         float64
	price              float64
	checkStateRetryCnt int
}

func NewMiner() *Miner {
	m := new(Miner)
	m.Mode = MINER_MODE
	m.symbol = TARGET_SYMBOL_IN + TARGET_SYMBOL_OUT
	return m
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

// 计算当前成交价格
func (m *Miner) calculatePrice(otype string) {
	data, _ := ApiInstance.GetTicker(m.symbol)
	mdata := data.(map[string]interface{})
	mmdata := mdata["data"].(map[string]interface{})
	ticker := mmdata["ticker"]
	values := ticker.([]interface{})
	m.price, _ = values[0].(float64)
	fmt.Println("Origin Price:", m.price)
	var price float64 = m.price
	if otype == "buy" {
		if m.Mode == MinerRadicalMode {
			price = m.price + PRICE_DEFF
		}
		if m.Mode == MinerConservatismMode {
			price = m.price - PRICE_DEFF
		}
		if m.Mode == MinerFastMode {
			price = m.price + PRICE_DEFF
		}
		if m.outBalance < PER_ORDER_AMOUNT*price {
			panic(TARGET_SYMBOL_OUT + "余额不足")
		}

	}

	if otype == "sell" {
		if m.Mode == MinerRadicalMode {
			price = m.price - PRICE_DEFF
		}
		if m.Mode == MinerConservatismMode {
			price = m.price + PRICE_DEFF
		}
		if m.Mode == MinerFastMode {
			price = m.price + PRICE_DEFF
		}
		if m.inBalance < PER_ORDER_AMOUNT*price {
			panic(TARGET_SYMBOL_IN + "余额不足")
		}
	}
	m.price = price
}

func (m *Miner) order(otype string) string {
	fmt.Println("下" + otype + "单！")
	fmt.Println("价格:", m.price)
	data, err := ApiInstance.CreateOrder(m.symbol, otype, "limit", strconv.FormatFloat(m.price, 'f', 6, 64), strconv.FormatInt(PER_ORDER_AMOUNT, 10), "main")
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
	m.calculatePrice(atype)
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

// 快速挖矿模式
func (m *Miner) fastWorker(atype string) {
	m.updateBalance()
	m.calculatePrice(atype)
	orderId := m.order(atype)
	m.checkStateRetryCnt = 0
	m.checkOrderState(orderId)
	m.wg.Done()
}
func (m *Miner) Start() {
	if MINER_MODE == MinerFastMode {
		for {
			m.wg.Add(2)
			go m.fastWorker("buy")
			go m.fastWorker("sell")
			m.wg.Wait()
		}
	} else {
		m.goWorker("buy")
	}
}

var FCoinMiner *Miner

func init() {
	FCoinMiner = NewMiner()
}
