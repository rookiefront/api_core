package model

type SysCity struct {
	Model
	Code       int    `json:"code"`
	ParentCode int    `json:"parent_code"`
	Name       string `json:"name"`
	//ProvinceCode int    `json:"province_code"`
	//CityCode     int    `json:"city_code"`
}

func (SysCity) TableName() string {
	return "sys_city"
}
