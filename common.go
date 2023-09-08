package qgnet

import (
	"time"
)

// 动态共享代理（按时）中国
func NewIpManagerByDyChinaIpByTime(username, password string, ipsize int) (*IpManager, error) {
	return NewIpManager(ipsize, NewDyChinaIpByTimeClient(username, password))
}

// 动态共享代理（按时）世界
func NewIpManagerByDyWorldIpByTime(username, password string, ipsize int) (*IpManager, error) {
	return NewIpManager(ipsize, NewDyWorldIpByTimeClient(username, password))
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
	ExpAt() time.Time               // 代理IP有效期
	Serve() string                  // 代理IP的服务
	Proxy(user, pass string) string // 代理IP的连接
	String() string                 // 代理IP的信息
}
