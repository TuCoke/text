package service

import (
	"github.com/gogf/gf/frame/g"
	"mian/model"
	"github.com/ahmetb/go-linq"
	"mian/util"
	"time"
)


//批量保存一房一价
func BatchSaveRoom(list []*model.Room) error {
	var buildingHashId []string
	var roomId []string
	linq.From(list).DistinctByT(func(r *model.Room) string{return r.BuildingHashId}).SelectT(func(r *model.Room) string{return r.BuildingHashId}).ToSlice(&buildingHashId)
	linq.From(list).SelectT(func(r *model.Room) string{return r.RoomId}).ToSlice(&roomId)

	var oldList []*model.Room
	err:=g.DB().Table("room").Structs(&oldList,"building_hash_id in(?) and room_id in(?)",buildingHashId,roomId)
	if err!=nil {
		return err
	}
	for _,room:=range list{
		room.UpdateDate = util.FormatDate(time.Now(),1)
		//判断是否存在数据
		query:=linq.From(oldList).WhereT(func(r *model.Room) bool{return (r.BuildingHashId==room.BuildingHashId && r.RoomId==room.RoomId)})
		if query.Any() {
			oldRoom:=query.First().(*model.Room)
			//一房一价Id
			room.Id=oldRoom.Id
			//成交状态
			if room.Status=="已售" && oldRoom.Status!="已售" {
				//util.FormatDate(time.Now(),1)
				dealDate:= util.GetDealDate()
				room.DealDate=&dealDate
			}else {
				room.DealDate=oldRoom.DealDate
			}
			//旧版api数据判断
			// 分摊面积
			if room.PublicArea==0 {
				room.PublicArea=oldRoom.PublicArea
			}
			//套内面积
			if room.InsideArea==0 {
				room.InsideArea=oldRoom.InsideArea
			}
			//楼层(没有楼层)
			if len(room.FloorName)==0 {
				room.FloorName=oldRoom.FloorName
			}

			switch room.CityId {
			//温州价格判断
			case 31:
				room.Area       =oldRoom.Area        //面积
				room.PublicArea =oldRoom.PublicArea  //公摊面积
				room.InsideArea =oldRoom.InsideArea  //套内面积
				room.RoomUsage  =oldRoom.RoomUsage   //用途
				room.RecordPrice=oldRoom.RecordPrice //备案价
				room.TotalPrice =oldRoom.TotalPrice  //总价
			}
		}
	}
	////_,err=g.DB().Table("room").Data(list).Save() //.Batch(100)
	//批量写入数据
	_,err=g.DB().Save("room",list,500)
	return err
}
