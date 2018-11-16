package fcoin

import (
	"fmt"
	"testing"
)

func TestGetCurrencies(t *testing.T) {
	data, err := ApiInstance.GetCurrencies()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}

func TestGetServerTime(t *testing.T) {
	data, err := ApiInstance.GetServerTime()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}

func TestGetSymbols(t *testing.T) {
	data, err := ApiInstance.GetSymbols()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}

func TestGetBalance(t *testing.T) {
	data, err := ApiInstance.GetBalance()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}

func TestCreateOrder(t *testing.T) {
	data, err := ApiInstance.CreateOrder("fjusdt", "buy", "limit", "0.001252", "10", "main")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}

func TestQueryOrders(t *testing.T) {
	data, err := ApiInstance.QueryOrders("fjusdt", "submitted", 11, 0, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}
