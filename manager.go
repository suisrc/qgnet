package qgnet

import (
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
	// github.com/bluele/gcache ? 可过期缓存
)

type IpManager struct {
	ipsize  int          // 管理器种队列的大小
	qrsize  int          // 单次获取的ip数量
	qcache  chan IpProxy // 代理IP队列
	running bool         // 代理IP队列是否正在运行
	H       IpHandler    // 代理IP的控制器,handler
}

// ipsize: 管理器种队列的大小，每次获取的ip的数量为该缓存大小的1.2倍， 最大不超过限制的ip数量
// 这样设计是可异步阻塞，分开处理，更高效, 同时 1.2 倍是为了防止ip没有使用就失效回收
func NewIpManager(ipsize int, handler IpHandler) (*IpManager, error) {
	if ipsize < 1 {
		ipsize = 1
	}
	return &IpManager{
		ipsize:  ipsize,
		qcache:  make(chan IpProxy, ipsize),
		qrsize:  ipsize * 6 / 5, // 1.2倍
		running: false,
		H:       handler,
	}, nil
}

// =========================================================
var _ io.Closer = (*IpManager)(nil)

func (aa *IpManager) Close() error {
	aa.Stop()        // 终止服务
	close(aa.qcache) // 关闭管道
	for ip := range aa.qcache {
		aa.H.Free(ip)
	}
	return nil
}

func (aa *IpManager) Stop() {
	aa.running = false // 终止服务, 提取服务
}

// =========================================================
func (aa *IpManager) RunAsync() {
	go aa.Run()
}

// 执行提取服务
func (aa *IpManager) Run() {
	if aa.running {
		return // 已经运行， 不需要重复运行
	}
	aa.running = true
	defer aa.Stop()
	for aa.running {
		// 获取一批代理IPs
		ips, err := aa.H.Gets(aa.qrsize)
		if err != nil {
			logrus.Errorf("hander gets ip error: %s", err.Error())
			return // 发生错误，退出，重新请求后会重写触发
		} else if len(ips) == 0 {
			logrus.Errorf("hander gets ip empty")
			return // 无法获取，退出，重新请求后会重写触发
		}
		for _, ip := range ips {
			aa.qcache <- ip // 提取代理IP到队列中
		}
	}
}

// =========================================================
// 获取一个可用的代理IP， 带有缓存形式和预获取模式

func (aa *IpManager) Get(timeout time.Duration) (IpProxy, error) {
	if !aa.running {
		aa.RunAsync()
		aa.H.WaitForFirst() // 如果新运行，需要等待，一遍可远程获取
	}
	for {
		select {
		case ip := <-aa.qcache:
			if time.Until(ip.ExpAt()) < time.Second*10 {
				// 如果ip[已经/快要]过期，释放, 10s，防止无法完成一次请求
				aa.H.Free(ip)
				continue
			}
			return ip, nil
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout: ip manager get ip")
		}
	}
}
