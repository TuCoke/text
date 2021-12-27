package util

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/gomodule/redigo/redis"
	"github.com/guonaihong/gout"
	"go_spider/core/common/mlog"
	"net/http"
	"regexp"
	"time"
)

//获取JsessionId
func HuJsessionId(name string,domain string,account string) (string,error) {
	cachekey:=fmt.Sprintf("%s_jsessionid",domain)
	exist,_:=redis.Bool(g.Redis().Do("EXISTS",cachekey))
	if exist {
		cache,_:=g.Redis().DoVar("GET",cachekey)
		return gconv.String(cache),nil
	}
	//自动登录
	var code = 0
	header:=&RspHeader{}
	err:=gout.New(&http.Client{}).POST(fmt.Sprintf("http://%s/mem/WebMemLoginAction_doMemLogin.jspx",domain)).SetTimeout(time.Second*15).SetHeader(gout.H{
		"Content-Type":"application/x-www-form-urlencoded; charset=UTF-8",
		"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36",
	}).SetBody(account).SetTimeout(time.Second*15).BindHeader(header).Code(&code).Do()
	if err!=nil {
		return "",err
	}
	if code!=http.StatusOK {
		return "",errors.New(fmt.Sprintf("status code:%d",code))
	}
	if header.Cookie==nil || len(header.Cookie)==0 {
		return "",errors.New("cookie notfound")
	}
	jsessionid:=regexp.MustCompile(`JSESSIONID=(.*?);`).FindStringSubmatch(header.Cookie[0])[0]
	g.Redis().Do("SETEX", cachekey, 60*20, jsessionid)
	mlog.StraceInst().Println(fmt.Sprintf("%s获取Cookie：%s",name,jsessionid))
	return jsessionid,nil
}
//endregion