package task

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"go_spider/core/common/mlog"
	"mian/model/task_model"
	"mian/service"
)

//楼幢任务
func BuildingTask(groupId int,typeId int,where string)  {
	list,err:=service.PresellList(fmt.Sprintf("%s and not EXISTS(SELECT 1 FROM s_task a where a.group_id=%d and a.type=%d and a.pid=n.presell_id)",where,groupId,typeId))
	if len(list)>0 {
		var task []*task_model.Task
		for  _,presell:=range list{
			info,_:=json.Marshal(g.Map{})
			task= append(task, &task_model.Task{GroupId:groupId,Name:presell.PresellName,Info:string(info),Type:typeId,Pid:presell.PresellId})
		}
		err=service.Task.SaveTask(task)
		if err!=nil {
			g.Log().Panicf("写入楼幢任务失败：%s",err.Error())
		}else {
			mlog.StraceInst().Println("写入楼幢任务成功")
		}
	}else {
		mlog.StraceInst().Println(fmt.Sprintf("%d条楼幢任务可写入",len(list)))
	}
}