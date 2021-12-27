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
	"mian/model/task_model"
	"mian/pipline"
	"mian/service"
	"mian/util"
)

//region 抓取数据
type HousePageProcesser struct {
	City *CityLiShui
}
func NewHousePageProcesser(city *CityLiShui) *HousePageProcesser {return &HousePageProcesser{city}}
func (this *HousePageProcesser) Request(req *request.Request) {}

func (this *HousePageProcesser) Process(p *page.Page) {
	if !p.IsSucc() {
		return
	}
	if p.GetStatusCode() != 200 {
		p.SetStatus(true, fmt.Sprintf("status code：%d", p.GetStatusCode()))
		return
	}
	fmt.Println("json的格式开头", p.GetBodyStr())
	city:= this.City
	//homePage:=p.GetRequest().GetMeta().(*task_model.HomePageViewModel1)
	var houseList []*model.House
	var presellList []*model.Presell
    for _,house:=range gjson.Get(p.GetBodyStr(),"data.jrZrcjxx").Array(){
		// fmt.Println("house", house.Get("id").String())
		//id
		id:=house.Get("id").String()
		// name:=house.Get("ysxmmc").String()

		lpmc:=house.Get("lpmc").String()
		//可售套数
		ksts:=house.Get("ksts").String()
		//销售套数
		slts:=house.Get("slts").String()
		//开发商
		// developer:=house.Get("xmgs").String()
		//地址
		//address:=house.Get("xmzl").String()
		//预售证
		//presellName:=house.Get("xkz").String()
		//公示时间
		//qfrq:=house.Get("qfrq").String()
		kprq:=house.Get("kprq").String()
		//楼盘
		info,_:=json.Marshal(g.Map{"id":id})
		// 生成 url + id
		hashId:= util.Base64Encrypt(city.Host+"_"+id +"")
		houseList = append(houseList, &model.House{
			//城市
			CityId:city.CityId,
			//名称
			Name:lpmc,
			//推广名
			PromotionName:lpmc,
			//销售状态
			Status:"",
			//地址
			Address:"",
			//开发商
			Developer:"",
			Property:"",
			Info:string(info),
			HashId:hashId,
		})

		//预售证
		presellHashId:=hashId
		var presellDate *string
		if len(kprq)>0 {
			presellDate=&kprq
		}else {
			presellDate=nil
		}
		//30_10
		status :=ksts+"_"+slts
		presellList= append(presellList, &model.Presell{
			//PresellName:name,
			PresellDate:presellDate,
			Info:string(info),
			StatsTime :city.TaskGroup.AddTimestamp,
			Stats: status,
			HashId:hashId,
			PresellHashId:presellHashId,
		})
	}
	p.AddField("houseList",houseList)
	p.AddField("presellList",presellList)
    meta:=p.GetRequest().GetMeta()
	if meta!= nil {
		homePage:=meta.(*task_model.HomePageViewModel1)
		p.AddField("homeTaskId", homePage.TaskId)
	}
}
//endregion

type HousePipeline struct {}
func NewHousePipeline() *HousePipeline {return &HousePipeline{}}

func (this *HousePipeline) Process(items *page_items.PageItems, t com_interfaces.Task) {
	allItems:=items.GetAll()
	houseList:=allItems["houseList"].([]*model.House)
	presellList:=allItems["presellList"].([]*model.Presell)
	succ:=true
	for _,house:=range houseList{
		err:=service.SaveHouse(house)
		if err!=nil {
			g.Log().Errorf("%s\n",err.Error())
			succ=false
		}
	}
	for _,presell:=range presellList{
		err:=service.SavePresell(presell)
		if err!=nil {
			g.Log().Errorf("%s\n",err.Error())
			succ=false
		}
	}
	if succ {
		fmt.Printf("楼盘写入成功 %d条，预售证写入成功 %d条 \n",len(houseList),len(presellList))
	}else {
		fmt.Println("楼盘写入失败")
	}
}
//endregion
func (city *CityLiShui) House() {
	groupId :=city.TaskGroup.Id
	//获取页码
	city.HomePage()
	for i:=0;i<3;i++{
		list,err:=service.Task.HomePageList(groupId)
		if err!=nil {
			g.Log().Panicf("获取宁波楼盘任务失败：%s\n",err.Error())
		}
		spider:=spider.NewSpider(NewHousePageProcesser(city), "").SetThreadnum(2).SetSleepTime("rand",300,600).AddPipeline(pipline.NewTaskHousePipeline())
		for _,homePage:=range list{
			//if(i>3){
			//	continue
			//}
			//url:=fmt.Sprintf("http://183.246.198.122:8099/api/Hourse/QueryTrading?xzqh=1&pageSize=25&pageNumber=%d",)
			url:=gjson.Get(homePage.Info,"url").String()
			fmt.Println("infoUrl")
			req:=request.NewRequest(url, "text","","GET","",city.Header(nil),nil,nil,homePage)
			spider.AddRequest(req)
		}
		spider.Run()
		g.Log().Info(fmt.Sprintf("丽水市楼盘抓取完成 共%d",len(list)))
	}

}
