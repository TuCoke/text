package task_model

type Task struct {
	Id int `json:"id"`
	GroupId int `json:"group_id"`
	Name string `json:"name"`
	Info string `json:"info"`
	Type int `json:"type"`
	Pid int `json:"pid"`
	Status int `json:"status"`
}
