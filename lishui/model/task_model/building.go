package task_model

type BuildingViewModel1 struct {
	TaskId int `json:"task_id"`
	BuildingId int `json:"building_id"`
	HashId string `json:"hash_id"`
	PresellHashId string `json:"presell_hash_id"`
	BuildingHashId string `json:"building_hash_id"`
	PromotionName string `json:"promotion_name"`
	BuildingName string `json:"building_name"`
	HouseInfo string `json:"house_info"`
	BuildingInfo string `json:"building_info"`
}