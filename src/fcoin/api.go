package fcoin

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type HttpHelper struct {
	params    map[string]string
	uri       string
	timestamp string
}

//
func NewHttpHelper(uri string, params map[string]string) *HttpHelper {
	t := time.Now()
	//fmt.Println(t.Unix())
	return &HttpHelper{uri: uri, params: params, timestamp: strconv.FormatInt(t.Unix()*1000, 10)}
}

func (h *HttpHelper) signature(method string) string {
	var signedBuffer bytes.Buffer
	signedBuffer.WriteString(method)
	signedBuffer.WriteString(API_URL)
	signedBuffer.WriteString(h.uri)

	if method == http.MethodGet {
		if h.params != nil {
			signedBuffer.WriteString("?")
			signedBuffer.WriteString(h.generateSortedParamsString())
		}
		signedBuffer.WriteString(h.timestamp)
	} else if method == http.MethodPost {
		signedBuffer.WriteString(h.timestamp)
		signedBuffer.WriteString(h.generateSortedParamsString())
	}
	//log.Println(signedBuffer.String())

	sign := base64.StdEncoding.EncodeToString(signedBuffer.Bytes())
	mac := hmac.New(sha1.New, []byte(API_SECRET))
	mac.Write([]byte(sign))
	sum := mac.Sum(nil)

	s := base64.StdEncoding.EncodeToString(sum)
	//log.Println(s)
	return s
}

func (h *HttpHelper) generateSortedParamsString() string {
	if h.params == nil || len(h.params) == 0 {
		return ""
	}
	var keys []string
	for k := range h.params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buffer bytes.Buffer
	for _, k := range keys {
		buffer.WriteString(k)
		buffer.WriteString("=")
		buffer.WriteString(h.params[k])
		buffer.WriteString("&")
	}
	str := buffer.String()
	str = str[:len(str)-1]
	fmt.Println(str)
	return str
}

// get请求
func (h *HttpHelper) Get() (interface{}, error) {
	var urlBuffer bytes.Buffer
	urlBuffer.WriteString(API_URL)
	urlBuffer.WriteString(h.uri)
	if h.params != nil {
		urlBuffer.WriteString("?")
		urlBuffer.WriteString(h.generateSortedParamsString())
	}
	fmt.Println(urlBuffer.String())

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, urlBuffer.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("FC-ACCESS-KEY", API_KEY)
	req.Header.Add("FC-ACCESS-SIGNATURE", h.signature(http.MethodGet))
	req.Header.Add("FC-ACCESS-TIMESTAMP", h.timestamp)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var data interface{}
	fmt.Println(string(body))
	e := json.Unmarshal(body, &data)
	if e != nil {
		return nil, e
	}
	return data, nil
}

// Post 请求json传递
func (h *HttpHelper) Post() (interface{}, error) {
	var urlBuffer bytes.Buffer
	urlBuffer.WriteString(API_URL)
	urlBuffer.WriteString(h.uri)
	fmt.Println(urlBuffer.String())
	client := &http.Client{}
	postData, _ := json.Marshal(h.params)
	req, err := http.NewRequest(http.MethodPost, urlBuffer.String(), strings.NewReader(string(postData)))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("FC-ACCESS-KEY", API_KEY)
	req.Header.Add("FC-ACCESS-SIGNATURE", h.signature(http.MethodPost))
	req.Header.Add("FC-ACCESS-TIMESTAMP", h.timestamp)
	//fmt.Println(h.timestamp)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	var data interface{}
	e := json.Unmarshal(body, &data)
	if e != nil {
		return nil, e
	}
	return data, nil
}

type FCoinApi struct {
}

//此 API 用于获取服务器时间。
func (api *FCoinApi) GetServerTime() (interface{}, error) {
	helper := NewHttpHelper("public/server-time", nil)
	return helper.Get()
}

//此 API 用于获取可用币种。
func (api *FCoinApi) GetCurrencies() (interface{}, error) {
	helper := NewHttpHelper("public/currencies", nil)
	return helper.Get()
}

//此 API 用于获取可用交易对。
func (api *FCoinApi) GetSymbols() (interface{}, error) {
	helper := NewHttpHelper("public/symbols", nil)
	return helper.Get()
}

//此 API 用于查询用户的资产列表。
func (api *FCoinApi) GetBalance() (interface{}, error) {
	helper := NewHttpHelper("accounts/balance", nil)
	return helper.Get()
}

//此 API 用于创建新的订单。
func (api *FCoinApi) CreateOrder(symbol, side, etype, price, amount, exchange string) (interface{}, error) {
	params := map[string]string{
		"symbol":   symbol,
		"side":     side,
		"type":     etype,
		"price":    price,
		"amount":   amount,
		"exchange": exchange,
	}
	helper := NewHttpHelper("orders", params)
	return helper.Post()
}

//此 API 用于查询订单列表。
func (api *FCoinApi) QueryOrders(symbol string, states string, before int64, after int64, limt int) (interface{}, error) {
	params := map[string]string{
		"symbol": symbol,
		"states": states,
	}
	if before != 0 {
		params["before"] = strconv.FormatInt(before, 10)
	}
	if after != 0 {
		params["after"] = strconv.FormatInt(after, 10)
	}
	if limt != 0 {
		params["limt"] = strconv.Itoa(limt)
	}
	helper := NewHttpHelper("orders", params)
	return helper.Get()
}

//此 API 用于返回指定的订单详情。
func (api *FCoinApi) GetOrder(orderId string) (interface{}, error) {
	helper := NewHttpHelper("orders/"+orderId, nil)
	return helper.Get()
}

//此 API 用于撤销指定订单，订单撤销过程是异步的，即此 API 的调用成功代表着订单已经进入撤销申请的过程，
//需要等待撮合的进一步处理，才能进行订单的撤销确认。
func (api *FCoinApi) CancelOrder(orderId string) (interface{}, error) {
	helper := NewHttpHelper("orders/"+orderId+"/submit-cancel", nil)
	return helper.Post()
}

//此 API 用于获取指定订单的成交记录
func (api *FCoinApi) GetOrderMatchResults(orderId string) (interface{}, error) {
	helper := NewHttpHelper("orders/"+orderId+"/match-results", nil)
	return helper.Get()
}

//此 API 用于获取指定订单的成交记录
func (api *FCoinApi) GetTicker(symbol string) (interface{}, error) {
	helper := NewHttpHelper("market/ticker/"+symbol, nil)
	return helper.Get()
}

var ApiInstance *FCoinApi

func init() {
	ApiInstance = new(FCoinApi)
}
