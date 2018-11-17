# fcoin-miner
由于Google相关的包无法直接下载，可以使用如果方式下载下来，然后创建软连接或者修改文件夹名称来解决：

第一步：

sudo git clone https://github.com/golang/net.git $GOPATH/src/github.com/golang/net

sudo git clone https://github.com/golang/sys.git $GOPATH/src/github.com/golang/sys

sudo git clone https://github.com/golang/tools.git $GOPATH/src/github.com/golang/tools

第二步：
sudo mkdir -p $GOPATH/src/golang.org/x

将net、sys、tools三个文件夹放到$GOPATH/src/golang.org/x目录下。

或者Linux下可以创建软连接：

sudo ln -s $GOPATH/src/github.com/golang $GOPATH/src/golang.org/x
