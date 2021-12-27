package page_processor

import (
	"github.com/gogf/gf/frame/g"
	"mian/model/task_model"
	"mian/service"
	"net/http"
)

type CityLiShui struct {
	CityId    int
	Host      string
	TaskGroup *task_model.TaskGroup
}

//丽水市
func NewLiShui(addDate string) *CityLiShui {
	taskGroup, err := service.TaskGroup.GroupId(33, 0, addDate)
	if err!=nil {
		g.Log().Panic("获取taskGroup错误")
	}
	return &CityLiShui{CityId: 33, Host: "jsjlishui", TaskGroup: taskGroup}
}

func (city *CityLiShui) Header(h map[string]string) http.Header {
	header := make(http.Header)
	header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	header.Set("Host", "jsj.lishui.gov.cn")
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")
	if h != nil {
		for k, v := range h {
			header.Add(k, v)
		}
	}
	return header
}

//func (city *CityLiShui) GetStatus(result gjson.Result) string {
//	if
//	(result.Get("djmx").Value()!=nil && result.Get("djmx.djyxbz").Int()==1)||
//		(result.Get("fwgg.cfzt").Exists() && result.Get("fwgg.cfzt").Int()!=1300)||
//		(result.Get("fwgg.djbz").Exists() &&
//			((regexp.MustCompile("^[0-9]+$").MatchString(result.Get("fwgg.djbz").String()) && result.Get("fwgg.djbz").Int()!=0)||
//				(!regexp.MustCompile("^[0-9]+$").MatchString(result.Get("fwgg.djbz").String()) && result.Get("fwgg.djbz").String()!="0"))) {
//		return "暂停销售"
//	}
//	if result.Get("htqd").Value()!=nil {
//		sfbatg:=result.Get("htqd.sfbatg").String()
//		if regexp.MustCompile("^[0-9]+$").MatchString(sfbatg) {
//			if gconv.Int(sfbatg)==1 {
//				return "未售" //签约中
//			}else if gconv.Int(sfbatg)==2 {
//				return "已售" //资金审核中
//			}else if gconv.Int(sfbatg)==0 {
//				return "已售"
//			}else if gconv.Int(sfbatg)==5 {
//				return "资金审核未通过"
//			}else {
//				return "不可售"
//			}
//		}
//	}else if result.Get("fwgg").Value()!=nil {
//		if result.Get("fwgg.pzkszt").Int()==1010 && result.Get("fwgg.sjyxszt").Int()==1100 {
//			return "未售" //可售
//		}else if result.Get("fwgg.pzkszt").Int()==1020 {
//			return "不可售"
//		}else {
//			return "不可售"
//		}
//	}else {
//		return "不可售"
//	}
//	return ""
//}
//
//func (city *CityLiShui) GetRoomUsage(result gjson.Result) string{
//	ysfwyt:=result.Get("ysfw.ysfwyt").String()
//	if len(ysfwyt)>=2 && ysfwyt[0:2]=="10" {
//		return "住宅"
//	}
//	if len(ysfwyt)>=3 && ysfwyt[0:3]=="126" {
//		return "写字楼"
//	}
//	if len(ysfwyt)>=3 && ysfwyt[0:3]=="123" {
//		return "商贸"
//	}
//	if gconv.Int(ysfwyt) ==13103 {
//		return "车库"
//	}
//	return ""
//}
