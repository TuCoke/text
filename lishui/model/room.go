package model

type Room struct {
	Id int `json:"id"`
	UpdateDate string `json:"update_date"`
	DealDate *string `json:"deal_date"`
	Status string `json:"status"`
	RoomName string `json:"room_name"`
	Area float64`json:"area"`
	PublicArea float64 `json:"public_area"`
	InsideArea float64 `json:"inside_area"`
	RecordPrice float64 `json:"record_price"`
	DecorationPrice float64 `json:"decoration_price"`
	TotalPrice float64 `json:"total_price"`
	RoomUsage string `json:"room_usage"`
	FloorName string `json:"floor_name"`
	Unitname string `json:"unit_name"`
	CellLocation string `json:"cell_location"`
	RoomId string `json:"room_id"`
	HashId string `json:"hash_id"`
	PresellHashId string `json:"presell_hash_id"`
	BuildingHashId string `json:"building_hash_id"`

	CityId int `json:"-"`
	AreaId int `json:"-"`
}
