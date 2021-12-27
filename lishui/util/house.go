package util

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/util/gconv"
	"github.com/tidwall/gjson"
	"net/http"
)

func GetSalesStatus(status string) string {
	salesStatus:=map[string][]string{
		"未售":{"可售","未网签","在售"},
		"已售":{"已签约","已网签"},
		"预订":{"已经预定","已预定"},
		"限制":{"已限制"},
	}
	for k,m:=range salesStatus{
		for _,n:=range m{
			if status==n {
				return k
			}
		}
	}
	return status
}

func GetRoomUsage(usage string) string {
	roomUsage:=map[string][]string{
		"办公":{"办公用房","办公楼"},
	}
	for k,m:=range roomUsage{
		for _,n:=range m{
			if usage==n {
				return k
			}
		}
	}
	return usage
}

////region 一房一价
////物业类型
//func GetHousetype(types int) string {
//	housetype:=map[int]string{
//		1:"住宅",
//		2:"商贸",
//		3:"写字楼",
//		4:"其他",
//	}
//	if v,ok:=housetype[types];ok {
//		return v
//	}
//	return ""
//}
////销售状态
//func GetHousestate(status int) string {
//	housestate:=map[int]string{
//		1:"未售",           //可售
//		2:"已售",           //
//		3:"预订",           //已经预定
//		4:"限制房产",       //
//		5:"未纳入网上销售", //
//		6:"即将解限",       //
//		7:"待现售",         //
//	}
//	if v,ok:=housestate[status];ok {
//		return v
//	}
//	return "";
//}
////endregion

func GetNumber(s *goquery.Selection) string {
	number:=map[string]string{
		"numberzero" :"0",
		"numberone"  :"1",
		"numbertwo"  :"2",
		"numberthree":"3",
		"numberfour" :"4",
		"numberfive" :"5",
		"numbersix"  :"6",
		"numberseven":"7",
		"numbereight":"8",
		"numbernine" :"9",
		"numberdor"  :".",
	}
	var str bytes.Buffer
	s.Each(func(i int, s *goquery.Selection) {
		n,exit:=s.Attr("class")
		if exit {
			if v,ok:=number[n];ok {
				str.WriteString(v)
			}
		}
	})
	return str.String()
}

func CitySprintf(cityId int,areaId int,format string,args ...interface{}) string {
	return fmt.Sprintf("and m.city_id=%d and m.area_id=%d %s", cityId,areaId,fmt.Sprintf(format,args...))
}

func RoomRecord(record gdb.Record) (string,string,string,string,string,string,string) {
	house_info :=gjson.GetMany(gconv.String(record["house_info"]),"propertyid","siteid")
	building_info :=gjson.GetMany(gconv.String(record["building_info"]),"presellid","buildingid")
	siteid:=house_info[1].String()
	propertyid:=house_info[0].String()
	presellid:=building_info[0].String()
	buildingid:=building_info[1].String()
	hash_id:=gconv.String(record["hash_id"])
	presell_hash_id:=gconv.String(record["presell_hash_id"])
	building_hash_id:=gconv.String(record["building_hash_id"])
	return siteid,propertyid,presellid,buildingid,hash_id,presell_hash_id,building_hash_id
}

func RoomHeader(regain bool) http.Header {
	var cookie,ua string
	if !regain {
		cookie,ua,_=Tmsfcookie()
	}else {
		cookie,ua,_=RegainTmsfcookie()
	}
	header := make(http.Header)
	//header.Set("Accept","text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	//header.Set("Accept-Language","zh-CN,zh;q=0.9,en;q=0.8")
	header.Set("Host","www.tmsf.com")
	header.Set("Cookie",cookie)
	header.Set("User-Agent",ua)
	return header
}

//自定义用cus前缀