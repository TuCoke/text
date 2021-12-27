package service

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"mian/model/task_model"
	"mian/util"
)

var Task = new(taskService)
type taskService struct{}
//type=1楼盘、2预售证、3楼幢、4一房一价

//1
func (this *taskService) HomePageList(groupId int) ([]*task_model.HomePageViewModel1,error) {
	sql:=fmt.Sprintf("SELECT m.id task_id,m.info from s_task m where m.group_id=%d and m.type=%d and m.status=0 %s order by id asc",groupId,util.HouseType,"")
	var homePage []*task_model.HomePageViewModel1
	err:=g.DB().GetStructs(&homePage,sql)
	return homePage,err
}

//2
func (this *taskService) HouseList(groupId int,typeId int,where string) ([]*task_model.HouseViewModel1,error) {
	//util.PresellType
	sql:=fmt.Sprintf("SELECT n.id task_id,m.id,m.promotion_name,m.info,m.hash_id FROM house m,s_task n where m.id=n.pid and n.group_id=%d and n.type=%d and n.status=0 %s order by m.id asc",groupId,typeId,where)
	var house []*task_model.HouseViewModel1
	err:=g.DB().GetStructs(&house,sql)
	return house,err
}

//3
func (this *taskService) PresellList(groupId int,typeId int,where string,order string) ([]*task_model.PresellViewModel,error) {
	//util.BuildingType
	if len(order)==0 {
		order="n.presell_id ASC"
	}
	sql:=fmt.Sprintf("SELECT o.id task_id,n.presell_id,m.hash_id,n.presell_hash_id,m.promotion_name,n.presell_name,m.info house_info,n.info presell_info FROM house m,presell n,s_task o where m.hash_id=n.hash_id and n.presell_id=o.pid and o.group_id=%d and o.type=%d and o.status=0 %s ORDER BY %s",groupId,typeId,where,order)
	var presell []*task_model.PresellViewModel
	fmt.Println("+++++++++PresellList+++SQL+++++++++:",sql)
	err:=g.DB().GetStructs(&presell,sql)
	return presell,err
}

//5
//func (this *taskService) PresellStatsList(groupId int,where string) ([]*task_model.PresellViewModel,error) {
//	sql:=fmt.Sprintf("SELECT o.id task_id,n.presell_id,m.hash_id,n.presell_hash_id,n.presell_name,m.info house_info,n.info presell_info FROM house m,presell n,s_task o where m.hash_id=n.hash_id and n.presell_id=o.pid and o.group_id=%d and o.type=%d and o.status=0 %s ORDER BY n.presell_id ASC",groupId,util.PresellStatsType,where)
//	var presell []*task_model.PresellViewModel
//	err:=g.DB().GetStructs(&presell,sql)
//	return presell,err
//}

//4
func (this *taskService) BuildingList(groupId int,typeId int,where string) ([]*task_model.BuildingViewModel1,error) {
	//util.RoomType
	sql:=fmt.Sprintf("SELECT p.id task_id,o.building_id,m.hash_id,n.presell_hash_id,o.building_hash_id,m.promotion_name,o.building_name,m.info house_info,o.info building_info from house m,presell n,building o,s_task p where m.hash_id=n.hash_id and n.presell_hash_id=o.presell_hash_id and o.building_id=p.pid and p.group_id=%d and p.type=%d and p.status=0 %s order by o.building_id asc",groupId,typeId,where)
	var building []*task_model.BuildingViewModel1
	err:=g.DB().GetStructs(&building,sql)
	return building,err
}

func (this *taskService) UpdateTask(taskId int) error {
	_,err:=g.DB().Table("s_task").Update(g.Map{"status":1},"id=?",taskId)
	return err
}

func (this *taskService) SaveTask(task []*task_model.Task) error {
	_,err:=g.DB().Insert("s_task",task,500)
	//_,err:=g.DB().Table("s_task").Insert(task)
	fmt.Println("Task任务执行",err)
	return err
}

func (this *taskService) TaskCount(groupId int,typeId int) (int,error) {
	count,err:=g.DB().Table("s_task").Count("group_id=? and type=?",groupId,typeId)
	return count,err
}