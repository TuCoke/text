package page_processor

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/guonaihong/gout"
	"github.com/tidwall/gjson"
	"go_spider/core/common/page"
	"go_spider/core/common/request"
	"go_spider/core/spider"
	"mian/model"
	"mian/pipline"
	"mian/service"
	"net/http"
)

var Ids string

//type RoomRecord struct {
//	City *CityLiShui
//	Building *model.BuildingViewModel
//}

//region 抓取数据
type RoomPricePageProcesser struct {
	Token string
	City  *CityLiShui
}

func NewRoomPricePageProcesser(token string, city *CityLiShui) *RoomPricePageProcesser {
	return &RoomPricePageProcesser{
		token, city,
	}
}
func (this *RoomPricePageProcesser) Request(req *request.Request) {
	//region 获取token
	if len(this.Token) == 0 {
		header := gout.H{
			"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
			"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36",
		}
		text, statusCode := "", 0
		gout.POST("http://183.246.198.92:8003/epoint-web-lishui/rest/oauth2/token").
			SetHeader(header).SetBody("client_id=28a990ac-49c4-425a-a2c5-a896a82fe9ad&client_secret=8d394188-a02f-4164-8e84-296a6060ece1&grant_type=client_credentials").
			BindBody(&text).Code(&statusCode).Do()
		this.Token = gjson.Get(text, "custom.access_token").String()
		//fmt.Println("token:",gjson.Get(text, "custom.access_token").String())
	}
	building := req.GetMeta().(*model.BuildingViewModel)
	id := gjson.Get(building.BuildingInfo, "zrdid").String()
	fmt.Println("data_id", id)
	req.Postdata = "params=%7B'DQZRZID'%3A+'" + id + "'%7D&access_token=" + this.Token + ""
	//endregion
}

func (this *RoomPricePageProcesser) Process(p *page.Page) {
	fmt.Println("请求的data数据", p.GetRequest().GetPostdata())
	if !p.IsSucc() {
		return
	}
	if p.GetStatusCode() != 200 {
		p.SetStatus(true, fmt.Sprintf("status code：%d", p.GetStatusCode()))
		return
	}
	fmt.Println("code值", gjson.Get(p.GetBodyStr(), "code"))
	//if gjson.Get(p.GetBodyStr(),"code").Value() != 1 {
	//	p.SetStatus(true,"数据获取失败")
	//	return
	//}
	if p.GetStatusCode() == 403 {
		return
	}
	building := p.GetRequest().Meta.(*model.BuildingViewModel)
	var roomList []*model.Room
	//fmt.Println("room_json", p.GetBodyStr())
	for _, room := range gjson.Get(p.GetBodyStr(), "fwlist").Array() {
		//房屋id
		roomId := room.Get("fwid").String()
		unitName := room.Get("dyh").String()
		roomName := room.Get("fh").String()
		area := room.Get("jzmj").Float()
		publicArea := room.Get("ftmj").Float()
		insideArea := room.Get("tnmj").Float()
		recordPrice := room.Get("sqbadj").Float()
		floorName := room.Get("lc").String()
		status := room.Get("zt").String()
		roomUsage := room.Get("fwyt").String()

		//fmt.Printf("%s，%s，%f，%f，%f，%f，%s \n",roomId,roomName,area,publicArea,insideArea,recordPrice,floorName)
		roomList = append(roomList, &model.Room{
			//销售状态
			Status: status,
			//房号
			RoomName: roomName,
			//建筑面积
			Area: area,
			//分摊面积
			PublicArea: publicArea,
			//套内建筑面积
			InsideArea: insideArea,
			//参考单价
			RecordPrice: recordPrice,
			//装修价
			DecorationPrice: 0,
			//总价
			TotalPrice: area * recordPrice,
			//用途
			RoomUsage: roomUsage,
			//楼层
			FloorName: floorName,
			//单元
			Unitname: unitName,
			//单元格
			CellLocation: "0,0",
			//房间id
			RoomId:         roomId,
			HashId:         building.HashId,
			PresellHashId:  building.PresellHashId,
			BuildingHashId: building.BuildingHashId,
		})
	}
	p.AddField("roomList", roomList)
}

//endregion

func (city *CityLiShui) RoomPrice() {
	list, err := service.BuildingList(fmt.Sprintf("and m.city_id=%d ", city.CityId))
	if err != nil {
		g.Log().Panic(err)
	}
	header1 := make(http.Header)
	header1.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	header1.Set("Referer", "http://jsj.lishui.gov.cn/")
	spider := spider.NewSpider(NewRoomPricePageProcesser("", city), "").SetThreadnum(8).SetSleepTime("rand", 200, 600).AddPipeline(pipline.NewRoomPricePipeline())
	for index, building := range list {
		//if(index > 3){
		//	continue
		//}
		fmt.Println("index", index)
		zrdid := gjson.Get(building.BuildingInfo, "zrdid").String()
		Ids = zrdid
		fmt.Println("zrdid的值", zrdid)
		url := fmt.Sprintf("http://183.246.198.92:8003/epoint-web-lishui/rest/propertiesrest/getbuildingfwlist")
		req := request.NewRequest(url, "text", "", "POST", "", header1, nil, nil, building)
		spider.AddRequest(req)
	}
	spider.Run()
	g.Log().Info(fmt.Sprintf("丽水市一房一价抓取完成 共%d", len(list)))
}
