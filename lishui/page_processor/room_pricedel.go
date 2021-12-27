package page_processor

import (
	"go_spider/core/common/page"
	"go_spider/core/common/request"
	"go_spider/core/spider"
)

type RoomDelPricePageProcesser struct {

}

func NewRoomDelPricePageProcesser() *RoomDelPricePageProcesser{return &RoomDelPricePageProcesser{}}

func (this *RoomDelPricePageProcesser) Request(req *request.Request)  {

}

func (this *RoomDelPricePageProcesser) Process(p *page.Page){

     spider:=spider.NewSpider(NewRoomDelPricePageProcesser(),"").SetThreadnum(2).SetSleepTime("rand",300,600)
     //.AddPipeline()
     url := "http://183.246.198.92:8003/epoint-web-lishui/rest/propertiesrest/getfwdetail"
     req :=request.NewRequest(url,"text","","POST","",nil,nil,nil,"")
     spider.AddRequest(req)
     spider.Run()
	 //g.Log().Info(fmt.Sprintf("丽水市楼幢楼层抓取完成 共%d", len(list)))
}
