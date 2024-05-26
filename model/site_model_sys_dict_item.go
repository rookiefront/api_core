package model

type SysDictItem struct {
	Model
	Label    string `json:"label" gorm:"column:label;comment:标签"`
	Value    string `json:"value" gorm:"column:value;comment:标签"`
	Sort     int    `json:"sort" gorm:"column:sort;comment:排序"`
	Remark   string `json:"remark" gorm:"column:remark;type:text;comment:备注"`
	DictId   int    `json:"dictId" gorm:"column:dict_id;comment:dict_id"`
	DictType string `json:"dictType" gorm:"column:dict_type;comment:类型"`
}

func (SysDictItem) TableName() string {
	return "sys_dict_item"
}
