package qgnet

import (
	"net/http"
	"net/url"
	"time"

	"github.com/guonaihong/gout/dataflow"
)

// 动态共享代理（按时）中国
func NewIpManagerByDyChinaIpByTime(username, password string, ipsize int) (*IpManager, error) {
	return NewIpManager(ipsize, NewDyChinaIpByTimeClient(username, password))
}

// 动态共享代理（按时）世界
func NewIpManagerByDyWorldIpByTime(username, password string, ipsize int) (*IpManager, error) {
	return NewIpManager(ipsize, NewDyWorldIpByTimeClient(username, password))
}

// =====================================================================================
// 配置代理

func SetProxyToGout(flow *dataflow.DataFlow, ipp IpProxy, prof, user, pass string) *dataflow.DataFlow {
	if prof == "socks5" {
		return flow.SetSOCKS5(ipp.Proxy(prof, user, pass))
	}
	return flow.SetProxy(ipp.Proxy(prof, user, pass))
}

func SetProxyToHttp(cli *http.Client, ipp IpProxy, prof, user, pass string) error {
	pxy, err := url.Parse(ipp.Proxy(prof, user, pass))
	if err != nil {
		return err
	}
	cli.Transport = &http.Transport{
		Proxy: http.ProxyURL(pxy),
	}
	return nil
}

//=====================================================================================
// IP代理控制器

type IpHandler interface {
	Gets(int) ([]IpProxy, error) // 获取一批代理IPs
	Get1(bool) (IpProxy, error)  // 获取一个代理IP
	Free(IpProxy) error          // 释放一个代理IP
	WaitForFirst()               // 等待第一次获取
}

type IpProxy interface {
	ExpAt() time.Time // 代理IP有效期
	Serve() string    // 代理IP的服务
	String() string   // 代理IP的信息

	Proxy(prof, user, pass string) string // 代理IP的连接

}
