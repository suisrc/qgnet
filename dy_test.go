package qgnet_test

import (
	"testing"

	"github.com/suisrc/qgnet"

	"github.com/guonaihong/gout"
	"github.com/test-go/testify/assert"
)

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
