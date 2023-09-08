package qgnet

import (
	"fmt"
	"strconv"
	"time"

	"github.com/guonaihong/gout"
)

// 动态共享代理(世界)（按时）

type DyWorldIpByTimeClient struct {
	Endpoint string             // 服务地址
	Params   DyWorldIpProxiesCO // 公共参数
}

//====================================================================================

func NewDyWorldIpByTimeClient(username, password string) *DyWorldIpByTimeClient {
	return &DyWorldIpByTimeClient{
		Endpoint: "https://overseas.proxy.qg.net",
		Params: DyWorldIpProxiesCO{
			Key: username, // 是 公共参数，产品唯一标识。
			Pwd: password, // 否 密码鉴权，如果产品设置了密码，则必须传入密码。
		},
	}
}

//====================================================================================

var _ IpHandler = (*DyWorldIpByTimeClient)(nil)

func (cc *DyWorldIpByTimeClient) Gets(size int) ([]IpProxy, error) {
	pco := cc.Params
	pco.Num = size

	rst, err := cc.GetProxies(pco)
	if err != nil {
		return nil, err
	}

	ips := make([]IpProxy, len(rst.Data))
	for i := 0; i < len(rst.Data); i++ {
		ips[i] = &rst.Data[i]
	}
	return ips, nil
}

// 如果么有可用通道，q = true:复用可用通道
func (cc *DyWorldIpByTimeClient) Get1(q bool) (IpProxy, error) {
	pco := cc.Params
	pco.Num = 1

	rst, err := cc.GetProxies(pco)

	if err != nil {
		if err == ErrNoAvailableChannel && q {
			rst, err = cc.QueryProxies(pco.Key, pco.Pwd)
			if err != nil {
				// 无可用通道， 查询可用通道
				return nil, err
			} // else 有正在使用的可用通道
		} else {
			return nil, err // 获取代理失败
		}
	}
	if len(rst.Data) == 0 {
		return nil, ErrorOf(ErrCodeStatus, "获取代理失败: 无可用通道")
	}

	return &rst.Data[0], nil
}

func (cc *DyWorldIpByTimeClient) Free(info IpProxy) error {
	return nil // do nothing
}

func (cc *DyWorldIpByTimeClient) WaitForFirst() {
	time.Sleep(time.Second * 1) // 等待第一次获取
}

func (cc *DyWorldIpByTimeClient) SetToken(token string) {
	cc.Params.Key = token
}

// ====================================================================================
var _ IpProxy = (*DyWorldIpProxyRO)(nil)

func (cc *DyWorldIpProxyRO) String() string {
	return fmt.Sprintf("server: %s, area: %s, isp: %s, deadline: %s", cc.Server, cc.Area, cc.Isp, cc.Deadline)
}

func (cc *DyWorldIpProxyRO) Proxy(prof, user, pass string) string {
	if prof == "" {
		prof = "http"
	}
	if user == "" {
		return fmt.Sprintf("%s://%s", prof, cc.Server)
	} else if pass == "" {
		return fmt.Sprintf("%s://%s@%s", prof, user, cc.Server)
	}
	return fmt.Sprintf("%s://%s:%s@%s", prof, user, pass, cc.Server)
}

func (cc *DyWorldIpProxyRO) ExpAt() time.Time {
	return cc.Deadlin0
}

func (cc *DyWorldIpProxyRO) Serve() string {
	return cc.Server
}

//====================================================================================

//	{
//	    "code": "SUCCESS",
//	    "data": [{
//	        "proxy_ip": "123.54.55.24",
//	        "server": "123.54.55.24:59419",
//	        "area": "河南省商丘市",
//	        "isp": "电信",
//	        "deadline": "2023-02-25 15:38:36"
//	    }],
//	    "request_id": "83158ebe-be6c-40f7-a158-688741083edc"
//	}

// 获取代理
type DyWorldIpProxiesRO struct {
	ResultRO
	Data []DyWorldIpProxyRO `json:"data"`
}

// 代理详情
type DyWorldIpProxyRO struct {
	ProxyIp  string    `json:"proxy_ip"`
	Server   string    `json:"server"`
	Area     string    `json:"area"`
	Isp      string    `json:"isp"`
	Deadline string    `json:"deadline"`
	Deadlin0 time.Time `json:"-"` // Deadlin0 用于存储 Deadline 的时间格式
}

type DyWorldIpProxiesCO struct {
	Key    string `query:"key"`     // 是 公共参数，产品唯一标识。
	Pwd    string `query:"pwd"`     // 否 密码鉴权，如果产品设置了密码，则必须传入密码。
	Area   string `query:"area"`    // 否 按地区提取。支持多地区筛选，逗号隔开。比如：”350500,330700”。
	AreaEx string `query:"area_ex"` // 否 排除某些地区提取。支持多地区排除，用逗号隔开。比如：”440100,450000”。
	Isp    int    `query:"isp"`     // 否 按运营商提取。0：不筛选；1：电信；2：移动；3：联通。
	Num    int    `query:"num"`     // 否 提取个数，默认为1。
}

// 查询资源
// 接口请求域名： overseas.proxy.qg.net。
// 本接口 (/get) 用于动态共享代理（全球HTTP）产品提取IP的接口。
// 默认接口请求频率限制：60/分钟。
func (cc *DyWorldIpByTimeClient) GetProxies(co DyWorldIpProxiesCO) (*DyWorldIpProxiesRO, error) {
	if co.Key == "" {
		co.Key, co.Pwd = cc.Params.Key, cc.Params.Pwd
	}

	body := DyWorldIpProxiesRO{}
	code := 0
	err := gout.GET(cc.Endpoint + "/get").SetQuery(co).BindJSON(&body).Code(&code).Do()
	if err != nil {
		return nil, err
	}
	if code != 200 && body.Code == "" {
		return nil, ErrorOf(ErrCodeStatus, "获取代理失败: code = "+strconv.Itoa(code))
	} else if body.Code != ErrSuccess.Code {
		return &body, GetError(body.Code, body.Message) // 返回错误码
	}

	// 解析时间, 本地时区
	for i := 0; i < len(body.Data); i++ {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", body.Data[i].Deadline, time.Local)
		if err != nil {
			return nil, ErrorOf(ErrCodeStatus, "获取代理失败: "+err.Error())
		}
		body.Data[i].Deadlin0 = t
	}

	return &body, nil
}

//====================================================================================

// 查询资源
// 接口请求域名： overseas.proxy.qg.net。
// 本接口 (/query) 用于动态共享代理（全球HTTP）产品查询IP的接口。
// 默认接口请求频率限制：60/分钟。
func (cc *DyWorldIpByTimeClient) QueryProxies(key, pwd string) (*DyWorldIpProxiesRO, error) {
	if key == "" {
		key, pwd = cc.Params.Key, cc.Params.Pwd
	}
	co := gout.H{"key": key, "pwd": pwd}

	body := DyWorldIpProxiesRO{}
	code := 0
	err := gout.GET(cc.Endpoint + "/query").SetQuery(co).BindJSON(&body).Code(&code).Do()
	if err != nil {
		return nil, err
	}
	if code != 200 && body.Code == "" {
		return nil, ErrorOf(ErrCodeStatus, "获取代理失败: code = "+strconv.Itoa(code))
	} else if body.Code != ErrSuccess.Code {
		return &body, GetError(body.Code, body.Message) // 返回错误码
	}

	// 解析时间, 本地时区
	for i := 0; i < len(body.Data); i++ {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", body.Data[i].Deadline, time.Local)
		if err != nil {
			return nil, ErrorOf(ErrCodeStatus, "获取代理失败: "+err.Error())
		}
		body.Data[i].Deadlin0 = t
	}

	return &body, nil
}

//====================================================================================

//	{
//	    "code": "SUCCESS",
//	    "data": [
//	        {
//	            "area": "新加坡",
//	            "area_code": 990100,
//	            "isp": "Oracle",
//	            "isp_code": 1,
//	            "available": true
//	        }
//	    ],
//	    "request_id": "51024a8b-a8a5-4e78-9301-cb500a8c083e"
//	}
type DyWorldIpResourcesRO struct {
	ResultRO
	Data []DyWorldIpResourceRO `json:"data"`
}

type DyWorldIpResourceRO struct {
	Area      string `json:"area"`      // 地区
	AreaCode  int    `json:"area_code"` // 地区编码
	Isp       string `json:"isp"`       // 运营商
	IspCode   int    `json:"isp_code"`  // 运营商编码
	Available bool   `json:"available"` // 是否可用
}

// 查询地区和运营商
// 接口请求域名： overseas.proxy.qg.net。
// 本接口 (/resources) 用于动态共享代理（全球HTTP）产品查询资源地区的接口。
// 默认接口请求频率限制：60/分钟。
func (cc *DyWorldIpByTimeClient) GetResources(key, pwd string) (*DyWorldIpResourcesRO, error) {
	if key == "" {
		key, pwd = cc.Params.Key, cc.Params.Pwd
	}
	co := gout.H{"key": key, "pwd": pwd}

	body := DyWorldIpResourcesRO{}
	code := 0
	err := gout.GET(cc.Endpoint + "/resources").SetQuery(co).BindJSON(&body).Code(&code).Do()
	if err != nil {
		return nil, err
	}
	if code != 200 && body.Code == "" {
		return nil, ErrorOf(ErrCodeStatus, "获取代理失败: code = "+strconv.Itoa(code))
	} else if body.Code != ErrSuccess.Code {
		return &body, GetError(body.Code, body.Message) // 返回错误码
	}

	return &body, nil
}
