package task_model

type TaskGroup struct {
	Id int `json:"id"`
	CityId int `json:"city_id"`
	AreaId int `json:"area_id"`
	Name string `json:"name"`
	AddDate string `json:"add_date"`
	AddTimestamp int64 `json:"add_timestamp"`
}
