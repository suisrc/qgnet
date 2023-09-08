# qgnet

国内动态按时代理客户端： dy_china_ip.go
```go  例子

func TestChinaIP(t *testing.T) {
	// ip_manager 比 ip_handler 增加了预提取ip的功能， 更适合于大量使用ip的场景
	// ipm, err := qgnet.NewIpManagerByDyChinaIpByTime("username", "password", 1)
	ipc := qgnet.NewDyChinaIpByTimeClient("username", "")
	ipp, err := ipc.Get1(false)
	assert.Nil(t, err)

	// qgnet.SetProxyToGout(...)
	pxy := ipp.Proxy("", "username", "")
	t.Log(pxy)
	body := ""
	err = gout.GET("https://ipinfo.io/ip").SetProxy(pxy).BindBody(&body).Do()
	assert.Nil(t, err)
	t.Log(body)
	ipc.Free(ipp)
}

```

国际动态按时代理客户端： dy_world_ip.go
```go 例子

func TestWorldIP(t *testing.T) {
	// ip_manager 比 ip_handler 增加了预提取ip的功能， 更适合于大量使用ip的场景
	// ipm, err := qgnet.NewIpManagerByDyChinaIpByTime("username", "password", 1)
	ipc := qgnet.NewDyWorldIpByTimeClient("username", "")
	ipp, err := ipc.Get1(false)
	assert.Nil(t, err)

	// qgnet.SetProxyToGout(...)
	pxy := ipp.Proxy("", "username", "")
	t.Log(pxy)
	body := ""
	err = gout.GET("https://ipinfo.io/ip").SetProxy(pxy).BindBody(&body).Do()
	assert.Nil(t, err)
	t.Log(body)
	ipc.Free(ipp)
}

```

一个使用按时代理来实现每次请求更换一次IP的例子, 不推荐用于生产，如果用于生产，请选用 隧道代理产品
```go
package main

import (
	"io"
	"net"
	"time"
	"github.com/suisrc/qgnet"
	"github.com/sirupsen/logrus"
)

var ipmng *qgnet.IpManager

func main() {
	// curl -x xxx:zzz@127.0.0.1:9090 ipinfo.io/ip
	ipmng, _ = qgnet.NewIpManagerByDyChinaIpByTime(
		"", // username
		"", // password
		1,  // ipsize
	)

	addr := ":9090"
	// 创建一个 tcp 服务
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Panic(err)
	}
	defer ln.Close()

	logrus.Infof("server is running at %s", addr)

	for {
		// 等待客户端连接
		conn, err := ln.Accept()
		if err != nil {
			logrus.Error(err)
			continue
		}
		// 处理客户端请求
		go handleConn(conn)
	}

}

var ipc qgnet.IpProxy

func handleConn(ssc net.Conn) {
	defer ssc.Close()

	ip1, err := ipmng.H.Get1(false)
	// ip1, err := ipmng.Get(time.Second * 2)
	if err != nil {
		logrus.Errorf("[GET] ip1 error: %s", err.Error())
		if ipc == nil {
			return // 没有可用的代理IP
		}
		ip1 = ipc
	} else {
		ipc = ip1
	}
	tip := ip1.Serve()
	// 将流量转发到 tip 上
	ttc, err := net.DialTimeout("tcp", tip, time.Second*10)
	if err != nil {
		logrus.Errorf("[DIAL] to %s error: %s", tip, err.Error())
		return
	}
	go func() {
		if _, err := io.Copy(ttc, ssc); err != nil {
			logrus.Errorf("[COPY] to %s error: %s", tip, err.Error())
		}
	}()
	if _, err := io.Copy(ssc, ttc); err != nil {
		logrus.Errorf("[COPY] form %s error: %s", tip, err.Error())
	}

}
```