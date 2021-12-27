package task_model

type PresellViewModel struct {
	TaskId int `json:"task_id"`
	PresellId int `json:"presell_id"`
	HashId string `json:"hash_id"`
	PresellHashId string `json:"presell_hash_id"`
	PromotionName string `json:"promotion_name"`
	PresellName string `json:"presell_name"`
	HouseInfo string `json:"house_info"`
	PresellInfo string `json:"presell_info"`
}