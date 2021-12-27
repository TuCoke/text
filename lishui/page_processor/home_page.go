package page_processor

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/guonaihong/gout"
	"github.com/tidwall/gjson"
	"go_spider/core/common/mlog"
	"go_spider/core/common/page"
	"go_spider/core/common/request"
	"mian/model/task_model"
	"mian/service"
	"mian/util"
	"strings"
)

//region 抓取数据
type HomePagePageProcesser struct {}
func NewHomePagePageProcesser() *HomePagePageProcesser {return &HomePagePageProcesser{}}
func (this *HomePagePageProcesser) Request(req *request.Request) {
	//req.ProxyHost=util.DefProxyIP()
}

func (this *HomePagePageProcesser) Process(p *page.Page) {
	if !p.IsSucc() {
		return
	}
	if p.GetStatusCode()!=200 {
		p.SetStatus(true,fmt.Sprintf("status code：%d",p.GetStatusCode()))
		return
	}
	query:=p.GetHtmlParser()
	href,_:=query.Find(`li[class="PagedList-skipToLast"]>a`).Attr("href")
	//pageIdnex pageSize pageCount
	pageCount:=gconv.Int(href[strings.LastIndex(href,"_")+1:len(href)])
	p.AddField("pageCount",gconv.Int(pageCount))
}
//endregion

func (city *CityLiShui) HomePage() {
	groupId:=city.TaskGroup.Id
	typeId:=util.HouseType
	count,err:=service.Task.TaskCount(groupId,typeId)
	if err!=nil {
		g.Log().Panicf("获取楼盘任务数量失败，%s\n",err.Error())
	}
	if count>0 {
		mlog.StraceInst().Println("当前已经存在楼盘任务")
		return
	}
	json1,code:="",0
	gout.GET("http://183.246.198.122:8099/api/Hourse/QueryTrading?xzqh=1&pageSize=25&pageNumber=1").BindBody(&json1).Code(&code).Do()
	ss1 := gconv.Int(gjson.Get(json1,"data.lpxxtotal"))
	pageCount := (ss1) / 25
	//fmt.Println("json1", gjson.Get(json1,"data.lpxxtotal").String())
	//doc, err := goquery.NewDocument("http://jsj.lishui.gov.cn/col/col1229219553/index.html?type=all")
	//fmt.Println("html",doc.Find("#pageTow"))
	//spider:=spider.NewSpider(NewHomePagePageProcesser(),"").SetThreadnum(1)
	//url:="http://183.246.198.122:8099/api/Hourse/QueryTrading?xzqh=1&pageSize=25&pageNumber=1"
	//req:=request.NewRequest(url, "json","","GET","",city.Header(nil),nil,nil,nil)
	//item:=spider.GetByRequest(req)
	//pageCount,exist:=item.GetItem("pageTow")
	fmt.Println("页码", pageCount)
	if pageCount<= 0 {
		g.Log().Panic("获取页码失败")
	}

	var task []*task_model.Task

	for i:=1;i<=pageCount + 1;i++{
		url:=fmt.Sprintf("http://183.246.198.122:8099/api/Hourse/QueryTrading?xzqh=1&pageSize=25&pageNumber=%d",i)
		info,_:=json.Marshal(g.Map{"url":url})
		task= append(task, &task_model.Task{GroupId:groupId,Name:url,Info:string(info),Type:typeId,Pid:0})
	}
	fmt.Println("task的值", task)
	err=service.Task.SaveTask(task)
	if err!=nil {
		g.Log().Panicf("写入楼盘任务失败：%s\n",err.Error())
	}else {
		mlog.StraceInst().Println("写入楼盘任务成功")
	}
}