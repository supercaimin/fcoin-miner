package main

import (
	. "fcoin"
	"fmt"
)

func main() {
	data, err := ApiInstance.GetServerTime()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}
