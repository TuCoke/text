package util

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/gomodule/redigo/redis"
	"github.com/guonaihong/gout"
	"go_spider/core/common/mlog"
	"net/http"
)

//var locker sync.Mutex
const ProxyIpCacheKey = "proxyip:%d"
//获取代理IP
func ProxyIP(idx int) string {
	cachekey:=fmt.Sprintf(ProxyIpCacheKey,idx)
	exist,_:=redis.Bool(g.Redis().Do("EXISTS",cachekey))
	if !exist {
		ip:=""
		err:=gout.New(&http.Client{}).GET("http://api.xiequ.cn/VAD/GetIp.aspx?act=get&num=1&time=30&plat=1&re=0&type=2&so=1&ow=1&spl=1&addr=&db=1").BindBody(&ip).Do()
		if err==nil && len(ip)>0 {
			g.Redis().Do("SETEX", cachekey,29, fmt.Sprintf("http://%s",ip))
			mlog.StraceInst().Println(fmt.Sprintf("proxy IP：%s",fmt.Sprintf("%d：http://%s",idx,ip)))
		}else {
			mlog.StraceInst().Println(fmt.Sprintf("ip获取异常，error：%s，ip：%s",err.Error(),ip))
		}
	}
	v, err := g.Redis().Do("GET", cachekey)
	if err==nil {
		return gconv.String(v)
	}else {
		return ""
	}
}

func RemoveProxyIP(idx int)  {
	cachekey:=fmt.Sprintf(ProxyIpCacheKey,idx)
	exist,_:=redis.Bool(g.Redis().Do("EXISTS",cachekey))
	if exist{
		g.Redis().Do("DEL",cachekey)
	}
}

//默认代理IP
func DefProxyIP() string {
	return ProxyIP(0)
}

//删除默认代理IP
func RemoveDefProxyIP() {
	RemoveProxyIP(0)
}