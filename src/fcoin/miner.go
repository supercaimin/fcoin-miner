package fcoin

import (
	"strconv"
	"sync"
)

const (
	MinerNormalLevel       = 0 // 同样价格买入卖出
	MinerConservatismLevel = 1 // 低买、高卖
	MinerRadicalLevel      = 2 // 高买、低卖
)

type Miner struct {
	wg         sync.WaitGroup
	Level      int
	inBalance  float64
	outBalance float64
	price      float64
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
	data, _ := ApiInstance.GetTicker(TARGET_SYMBOL_IN + TARGET_SYMBOL_OUT)
	mdata := data.(map[string]interface{})
	mmdata := mdata["data"].(map[string]interface{})
	ticker := mmdata["ticker"]
	values := ticker.([]interface{})
	m.price, _ = values[0].(float64)
}

func (m *Miner) Work() {
	m.updateBalance()
	m.updatePrice()
}

var FCoinMiner *Miner

func init() {
	FCoinMiner = new(Miner)
}
