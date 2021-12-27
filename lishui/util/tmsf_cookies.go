package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/gomodule/redigo/redis"
	"github.com/guonaihong/gout"
	"github.com/tidwall/gjson"
	"go_spider/core/common/mlog"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

//var (
//	once sync.Once
//	seleniumService *selenium.Service
//)
type RspHeader struct {
	Cookie []string `header:"Set-Cookie"`
}
var TmsfCachekey = "tmsf_cookies_03"
//region 获取透明cookie
func Tmsfcookie() (string,string,error) {
	exist,_:=redis.Bool(g.Redis().Do("EXISTS",TmsfCachekey))
	//存在cookie
	if exist {
		cookieCache,_:=g.Redis().DoVar("GET",TmsfCachekey)
		cookie:=gjson.GetMany(gconv.String(cookieCache),"cookie","ua")
		return cookie[0].String(),cookie[1].String(),nil
	}
	jsessionId,err:=TmsfJsessionId()
	if err!=nil {
		return "","",err
	}
	//不存在cookie
	mlog.StraceInst().Println("开始获取cookie")
	rand.Seed(time.Now().UnixNano())
	ua:="Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/"+ gconv.String(rand.Intn(600)) +".36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/"+ gconv.String(rand.Intn(9999)) +".361"
	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.ProxyServer(ProxyIP(1)),
		chromedp.UserAgent(ua),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()
	// create context
	ctx, cancel := chromedp.NewContext(allocCtx,chromedp.WithLogf(log.Printf))
	defer cancel()
	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var bsfit_deviceid=""
	err = chromedp.Run(ctx,
		chromedp.Navigate(`http://www.tmsf.com/`),
		//设置
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := network.SetCookie("JSESSIONID", jsessionId).
				//WithExpires(&expr).
				WithDomain("www.tmsf.com").
				WithHTTPOnly(true).
				Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for i:=0;i<15;i++{
				cookes,err:=network.GetAllCookies().Do(ctx)
				if err!=nil {
					return err
				}
				exists:= linq.From(cookes).AnyWithT(func(c *network.Cookie) bool{return c.Name=="BSFIT_DEVICEID"})
				if(exists){
					for _, v := range cookes {
						bsfit_deviceid = bsfit_deviceid + v.Name + "=" + v.Value + ";"
					}
					break
				}
				time.Sleep(time.Second*1)
			}
			return nil
		}),
	)
	if err != nil {
		return  "","",err
	}
	json,err:=json.Marshal(g.Map{"cookie":bsfit_deviceid,"ua":ua})
	if err!=nil {
		return  "","",err
	}
	//写入redis
	g.Redis().Do("SETEX", TmsfCachekey, 60*2, string(json))
	mlog.StraceInst().Println(fmt.Sprintf("杭州获取Cookie：%s",bsfit_deviceid))
	return bsfit_deviceid,ua,nil
}
//获取JsessionId
func TmsfJsessionId() (string,error) {
	cachekey:="tmsf_jsessionid"
	exist,_:=redis.Bool(g.Redis().Do("EXISTS",cachekey))
	if exist {
		cache,_:=g.Redis().DoVar("GET",cachekey)
		return gconv.String(cache),nil
	}

	//获取JSESSIONID
	var code = 0
	header:=&RspHeader{}
	err:=gout.New(&http.Client{}).GET(`http://www.tmsf.com/mem/main.htm`).SetProxy(ProxyIP(1)).SetTimeout(time.Second*15).BindHeader(header).Code(&code).Do()
	if err!=nil {
		return "",err
	}
	if code!=http.StatusOK {
		return "",errors.New(fmt.Sprintf("status code:%d",code))
	}
	if header.Cookie==nil || len(header.Cookie)==0 {
		return "",errors.New("cookie notfound")
	}
	jsessionid:=regexp.MustCompile(`JSESSIONID=(.*?);`).FindStringSubmatch(header.Cookie[0])[1]
	//自动登录
	var res=""
	err=gout.New(&http.Client{}).POST(`http://www.tmsf.com/mem/WebMemLoginAction_doMemLogin.jspx`).SetProxy(ProxyIP(1)).SetTimeout(time.Second*15).SetHeader(gout.H{
		"Content-Type":"application/x-www-form-urlencoded; charset=UTF-8",
		"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36",
		"Cookie":fmt.Sprintf("JSESSIONID=%s",jsessionid),
	}).SetBody("webMem.username=18411631100&webMem.userpwd=hjkl6789").Code(&code).BindBody(&res).Do()
	if err!=nil {
		return "",err
	}
	if code!=http.StatusOK {
		return "",errors.New(fmt.Sprintf("status code:%d",code))
	}
	g.Redis().Do("SETEX", cachekey, 60*10, jsessionid)
	return jsessionid,nil
}
//endregion

//region 重新获取cookie
func RegainTmsfcookie() (string,string,error) {
	exist,_:=redis.Bool(g.Redis().Do("EXISTS",TmsfCachekey))
	if exist{
		g.Redis().Do("DEL",TmsfCachekey)
	}
	return Tmsfcookie()
}
//endregion