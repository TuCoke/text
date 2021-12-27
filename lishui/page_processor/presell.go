package page_processor

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/tidwall/gjson"
	"go_spider/core/common/com_interfaces"
	"go_spider/core/common/page"
	"go_spider/core/common/page_items"
	"go_spider/core/common/request"
	"go_spider/core/spider"
	"mian/model"
	"mian/service"
)

//region 抓取数据
type PressllDelProcesser struct{

}
func NewPressllDelProcesser() *PressllDelProcesser{return &PressllDelProcesser{}}
func(this *PressllDelProcesser) Request(req *request.Request){

}
func(this *PressllDelProcesser) Process(p *page.Page){
	if !p.IsSucc() {
		return
	}

	if p.GetStatusCode() == 500{
		fmt.Println("不存在预售证编号!")
		p.SetSkip(true)
		return
	}
	if p.GetStatusCode() != 200 {
		p.SetStatus(true, fmt.Sprintf("status code：%d", p.GetStatusCode()))
		return
	}
	var presellList []*model.Presell
	for _,presell:=range gjson.Get(p.GetBodyStr(),"data").Array() {
		presellName:=presell.Get("yszbhstr").String()
		fmt.Println("名称",presellName)
		ysxmid :=presell.Get("ysxmid").String()
		fmt.Println("id的值", presell.Get("ysxmid").String())
		presellList = append(presellList, &model.Presell{
			Info : ysxmid,
			PresellName:presellName,
		})


	}
	p.AddField("presellList",presellList)
}
//endregion
//region 写入
type PressllDelPipline struct {}
func NewPressllDelPipline() *PressllDelPipline{return &PressllDelPipline{}}
func(this *PressllDelPipline) Process(items *page_items.PageItems,t com_interfaces.Task){
	allItem:=items.GetAll()
	presellList :=allItem["presellList"].([]*model.Presell)

	for _,presell:= range presellList{
		info,_:=json.Marshal(g.Map{"id":presell.Info})
        err:=service.UpdatePresellName(string(info),presell.PresellName)
        if err!=nil{
			g.Log().Errorf("%s\n",err.Error())
		}
	}
}
//endregion

func (city *CityLiShui) PreallName()  {
	list,err:=service.PresellList(fmt.Sprintf("and m.city_id=%d and n.presell_name is not null AND n.presell_name =''" ,city.CityId))
	if err!=nil {
		g.Log().Panicf("获取宁波楼盘任务失败：%s\n",err.Error())
	}
	fmt.Println("list的长度",len(list))
	spider:=spider.NewSpider(NewPressllDelProcesser(), "").SetThreadnum(2).SetSleepTime("rand",300,600).AddPipeline(NewPressllDelPipline())
	for _,presell:=range list{
		id:=gjson.Get(presell.PresellInfo,"id").String()
		fmt.Println("截取的id为",id)
		url:=fmt.Sprintf("http://183.246.198.122:8099/api/Hourse/QueryXKZDetailsForXMID?id=%s",id)
		req:=request.NewRequest(url, "text","","GET","",city.Header(nil),nil,nil,nil)
		spider.AddRequest(req)
	}
	spider.Run()
	g.Log().Info(fmt.Sprintf("丽水市楼幢抓取完成 共%d",len(list)))
}
