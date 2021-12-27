package service

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"go_spider/core/common/mlog"
	"mian/model/task_model"
	"mian/util"
	"time"
)

var TaskGroup = new(taskGroupService)
type taskGroupService struct{}

func (this *taskGroupService) GroupId(city_id int,area_id int,addDate string) (*task_model.TaskGroup,error) {
	if len(addDate)==0{
		if gtime.Now().Hour()<18 {
			addDate=util.FormatDate(time.Now().Add(-time.Hour*24),0)
		}else {
			addDate=util.FormatDate(time.Now(),0)
		}
	}
	//时间戳
	timestamp:=gtime.NewFromStr(addDate).Timestamp()
	//addDate
	mlog.StraceInst().Println(fmt.Sprintf(`=================** addDate:%s **=================`,addDate))
	oldTaskGroup:=&task_model.TaskGroup{}
	err:=g.DB().Table("s_task_group").Struct(oldTaskGroup,"city_id=? and area_id=? and add_date=? order by id desc",city_id, area_id, addDate)
	if err==sql.ErrNoRows {
		taskGroup:=&task_model.TaskGroup{}
		taskGroup.CityId=city_id
		taskGroup.AreaId=area_id
		taskGroup.Name=""
		taskGroup.AddDate=addDate
		taskGroup.AddTimestamp=timestamp

		result,err:=g.DB().Exec("INSERT INTO s_task_group(city_id,area_id,name,add_date,add_timestamp) VALUES(?,?,?,?,?)",
			taskGroup.CityId,taskGroup.AreaId,taskGroup.Name,taskGroup.AddDate,taskGroup.AddTimestamp)
		if err!=nil {
			return nil,err
		}
		insertId,err:=result.LastInsertId()
		if err!=nil {
			return nil,err
		}
		taskGroup.Id=gconv.Int(insertId)
		return taskGroup,nil
	}else if err==nil {
		return oldTaskGroup,nil
	}
	return nil,err
}