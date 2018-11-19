# FCoin挖矿机器人
用Golang实现的挖矿机器人，适用于FCoin、FCoinJP，支持Linux、Mac、Windows。
下面是Linux、Mac下的使用方法：
## 使用方法
### 1.配置golang运行环境
第一步：
```shell
sudo apt-get install golang
```
第二步：
下载依赖包，由于Google相关的包无法直接下载，可以使用如果方式下载下来，然后创建软连接或者修改文件夹名称来解决：
```shell
sudo git clone https://github.com/golang/net.git $GOPATH/src/github.com/golang/net

sudo git clone https://github.com/golang/sys.git $GOPATH/src/github.com/golang/sys

sudo git clone https://github.com/golang/tools.git $GOPATH/src/github.com/golang/tools
```
第三步：
```shell
sudo mkdir -p $GOPATH/src/golang.org/x
```
将net、sys、tools三个文件夹放到$GOPATH/src/golang.org/x目录下。

或者Linux下可以创建软连接：
```shell
sudo ln -s $GOPATH/src/github.com/golang $GOPATH/src/golang.org/x
```
### 2.使用方法
第一步：
修改配置文件./src/fcoin/config.go
```golang
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
// 特别说明 MinerFastMode 模式下交易对象初始余额各占50%
const MINER_MODE = MinerFastMode

// 差价0.001522
const PRICE_DEFF = 0.000001

// 单笔订单购买数量
const PER_ORDER_AMOUNT = 10
```
第二步：
编译程序，运行install.sh文件
```shell
$ bash install.sh
finished
```
第三步：
运行编译生成的程序./bin/app
```shell
$ ./bin/app
https://api.fcoinjp.com/v2/accounts/balance
https://api.fcoinjp.com/v2/accounts/balance
{"status":0,"data":[{"currency":"usdt","available":"0.20330
..................................
```

挖矿开始

若有疑问请加微信：261029366
