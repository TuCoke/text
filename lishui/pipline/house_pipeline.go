package pipline

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"go_spider/core/common/com_interfaces"
	"go_spider/core/common/mlog"
	"go_spider/core/common/page_items"
	"mian/model"
	"mian/service"
)

type HousePipeline struct {}

func NewTaskHousePipeline() *HousePipeline {return &HousePipeline{}}

func (this *HousePipeline) Process(items *page_items.PageItems, t com_interfaces.Task) {
	allItems:=items.GetAll()
	houseList:=allItems["houseList"].([]*model.House)
	presellList:=allItems["presellList"].([]*model.Presell)
	succ:=true
	for _,house:=range houseList{
		err:=service.SaveHouse(house)
		if err!=nil {
			g.Log().Errorf("%s",err.Error())
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
		this.TaskProcess(items,fmt.Sprintf("楼盘写入成功 %d条",len(houseList)))
	}else {
		mlog.StraceInst().Println("楼盘写入失败")
	}
}

func (this *HousePipeline) TaskProcess(items *page_items.PageItems,message string) {
	item,exist:=items.GetItem("homeTaskId")
	if !exist {
		mlog.StraceInst().Println(message)
		return
	}
	err:=service.Task.UpdateTask(item.(int))
	if err!=nil {
		g.Log().Errorf("更新任务状态失败，%s",err.Error())
	}else {
		mlog.StraceInst().Println(fmt.Sprintf("%s，更新任务状态成功",message))
	}
}