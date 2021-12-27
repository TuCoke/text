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

type RoomPricePipeline struct {}
func NewRoomPricePipeline() *RoomPricePipeline {return &RoomPricePipeline{}}

func (this *RoomPricePipeline) Process(items *page_items.PageItems, t com_interfaces.Task) {
	item,exist:=items.GetItem("roomList")
	if !exist {
		return
	}
	roomList:=item.([]*model.Room)
	if len(roomList)>0 {
		err:=service.BatchSaveRoom(roomList)
		if err==nil {
			this.TaskProcess(items,fmt.Sprintf("一房一价写入成功 %d条",len(roomList)))
		}else {
			g.Log().Errorf("一房一价写入失败：%s",err.Error())
		}
	}else {
		mlog.StraceInst().Println("没有一房一价数据")
	}
}

func (this *RoomPricePipeline) TaskProcess(items *page_items.PageItems,message string) {
	item,exist:=items.GetItem("roomTaskId")
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
