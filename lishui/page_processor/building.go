package page_processor

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/guonaihong/gout"
	"github.com/tidwall/gjson"
	"go_spider/core/common/page"
	"go_spider/core/common/request"
	"go_spider/core/spider"
	"mian/model"
	"mian/model/task_model"
	"mian/pipline"
	"mian/service"
	"mian/task"
	"mian/util"
	"net/http"
)

var Pids string

type BuildingRecord struct {
	Presell *model.PresellViewModel
}

//region 抓取数据
type BuildingPageProcesser struct{
	City  *CityLiShui
	Token string
}
func NewBuildingPageProcesser(city *CityLiShui,token string) *BuildingPageProcesser { return &BuildingPageProcesser{City: city,Token: token} }
func (this *BuildingPageProcesser) Request(req *request.Request) {
	if(len(this.Token)==0){
		//region 获取token
		header := gout.H{
			"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
			"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36",
		}
		text, statusCode := "", 0
		gout.POST("http://183.246.198.92:8003/epoint-web-lishui/rest/oauth2/token").
			SetHeader(header).SetBody("client_id=28a990ac-49c4-425a-a2c5-a896a82fe9ad&client_secret=8d394188-a02f-4164-8e84-296a6060ece1&grant_type=client_credentials").
			BindBody(&text).Code(&statusCode).Do()
		this.Token=gjson.Get(text, "custom.access_token").String()
		fmt.Println("token:",gjson.Get(text, "custom.access_token").String())
	}
	presell := req.GetMeta().(*task_model.PresellViewModel)
	id := gjson.Get(presell.PresellInfo, "id").String()
	req.Postdata="params=%7B'id'%3A+'"+id+"'%7D&access_token="+this.Token+""
}

func (this *BuildingPageProcesser) Process(p *page.Page) {
	if !p.IsSucc() {
		return
	}
	fmt.Println("请求的状态code", p.GetStatusCode())
	if p.GetStatusCode() == 403 {
		this.Token=""
		p.SetStatus(true, "token过期")
		return
	}

	fmt.Println(p.GetBodyStr())
	fmt.Println("Building重新设置的header", p.GetRequest().Postdata)
	//请求失败
	if p.GetStatusCode() == 500 {
		p.SetSkip(true)
		fmt.Println("请求失败")
		return
	}
	fmt.Println("buildingcode", p.GetStatusCode())
	if p.GetStatusCode() != 200 {
		p.SetStatus(true, fmt.Sprintf("status code：%d", p.GetStatusCode()))
		return
	}
	fmt.Println("shuju", p.GetBodyStr())
	// jsp := gjson.Get(fmt.Sprintf("{data:%s}", p.GetBodyStr()), "")
	//fmt.Println("jsp", fmt.Sprintf(`{"data":%s}`, p.GetBodyStr()), "")
	//fmt.Println("data:",gjson.Get(fmt.Sprintf(`{"data":%s}`, p.GetBodyStr()),"data").Array())
	presell := p.GetRequest().GetMeta().(*task_model.PresellViewModel)
	taskId:=presell.TaskId
	//city := meta.City
	//presell := meta.Presell
	fmt.Println("PresellInfo", presell.PresellInfo)
	ids := gjson.Get(presell.PresellInfo, "id").String()
	var buildingList []*model.Building
	if gjson.Get(fmt.Sprintf(`{"data":%s}`, p.GetBodyStr()),"data").Array() ==nil{
		fmt.Println("GetBodyStrIsNil")
		return
	}
	//fmt.Println("postdata=====================:",p.GetRequest().GetPostdata())
	for _, building := range gjson.Get(fmt.Sprintf(`{"data":%s}`, p.GetBodyStr()), "data").Array() {
		buidlingName := building.Get("zh").String()
		//id:=building.Get("ysxmid").String()
		zrdid := building.Get("zid").String()
		fmt.Println("buidlingName", buidlingName)
		info, _ := json.Marshal(g.Map{"id": ids, "zrdid": zrdid})
		buildingHashId := util.Base64Encrypt(this.City.Host + "_" + zrdid)

		//fmt.Println("zid:",city.Host + "_" + zrdid)
		buildingList = append(buildingList, &model.Building{
			PresellName:    presell.PresellName,
			BuildingName:   buidlingName,
			Info:           string(info),
			HashId:         presell.HashId,
			PresellHashId:  presell.PresellHashId,
			BuildingHashId: buildingHashId,
		})
	}
	p.AddField("buildingTaskId",taskId)
	p.AddField("buildingList", buildingList)
}

//endregion

func (city *CityLiShui) Building() {
	groupId:=city.TaskGroup.Id
	typeId:=util.BuildingType
	//增加任务
	task.BuildingTask(groupId,typeId, util.CitySprintf(city.CityId,0,"and n.stats_time>=%d",city.TaskGroup.AddTimestamp))
	//执行任务
	list,err:=service.Task.PresellList(groupId,typeId,util.CitySprintf(city.CityId,0,""), "m.id ASC,n.presell_date ASC,n.presell_id ASC")
	if err != nil {
		g.Log().Panicf("获取预售证任务失败：%s",err.Error())
	}
	header1 := make(http.Header)
	header1.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	spider := spider.NewSpider(NewBuildingPageProcesser(city,""), "").SetThreadnum(3).SetSleepTime("rand", 200, 600).AddPipeline(pipline.NewBuildingPipeline())
	for index, presell := range list {
		//if(index>3){
		//	continue
		//}
		fmt.Println("index", index)
		urlPost := "http://183.246.198.92:8003/epoint-web-lishui/rest/propertiesrest/getpermitsale"
		//fmt.Println("presellid", presell.PresellInfo)
		//data := make(url.Values)
		id := gjson.Get(presell.PresellInfo, "id").String()
		Pids = id
		fmt.Println("Buildingrecord的ids的值", Pids)
		//tokens := "dc1697295ee9530716c9485a9f720d09"
		req := request.NewRequest(urlPost, "text", "", "POST", "", header1, nil, nil, presell)
		spider.AddRequest(req)
	}
	spider.Run()
	g.Log().Info(fmt.Sprintf("丽水市楼幢抓取完成 共%d", len(list)))

}

//region 自定义抓取
func (city *CityLiShui) WhereBuilding(strwhere string) {
	groupId:=city.TaskGroup.Id
	typeId:=util.BuildingType
	where:=util.CitySprintf(city.CityId,0,"%s",strwhere)
	//增加任务
	task.BuildingTask(groupId,typeId, where)
	//执行任务
	list,err:=service.Task.PresellList(groupId, typeId, where,"m.id ASC,n.presell_date ASC,n.presell_id ASC")
	if err!=nil {
		g.Log().Panicf("获取预售证任务失败：%s",err.Error())
	}
	header1 := make(http.Header)
	header1.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	spider := spider.NewSpider(NewBuildingPageProcesser(city,""), "").SetThreadnum(3).SetSleepTime("rand", 200, 600).AddPipeline(pipline.NewBuildingPipeline())
	for index, presell := range list {
		fmt.Println("index", index)
		urlPost := "http://183.246.198.92:8003/epoint-web-lishui/rest/propertiesrest/getpermitsale"
		id := gjson.Get(presell.PresellInfo, "id").String()
		Pids = id
		fmt.Println("Buildingrecord的ids的值", Pids)
		//tokens := "dc1697295ee9530716c9485a9f720d09"
		req := request.NewRequest(urlPost, "text", "", "POST", "", header1, nil, nil, presell)
		spider.AddRequest(req)
	}
	spider.Run()
	g.Log().Info(fmt.Sprintf("丽水市楼幢抓取完成 共%d", len(list)))

}
//endregion