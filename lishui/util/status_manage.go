package util

import "github.com/ahmetb/go-linq"

type SpiderStatus struct {
	Id int
	Page int
	Status bool
}

type StatusManage struct {
	Status []*SpiderStatus
}

func NewStatusManage() *StatusManage {
	status:=make([]*SpiderStatus,0)
	return &StatusManage{ Status: status}
}

func (statusManage *StatusManage) AddStatus(id int,page int)  {
	statusManage.Status = append(statusManage.Status, &SpiderStatus{id,page,false})
}

func (statusManage *StatusManage) SetStatus(id int,page int)  {
	linq.From(statusManage.Status).WhereT(func(s *SpiderStatus) bool{return s.Id==id && s.Page==page}).ForEachT(func(s *SpiderStatus) {
		s.Status=true
	})
}

func (statusManage *StatusManage) CloseStatus() {
	statusManage.Status=make([]*SpiderStatus,0)
}

func (statusManage *StatusManage) StatusSucc(id int) bool {
	query:=linq.From(statusManage.Status).WhereT(func(s *SpiderStatus) bool{return s.Id==id})
	if query.Count()>0 && !query.WhereT(func(s *SpiderStatus) bool{return s.Status==false}).Any() {
		return true
	}
	return false
}