package model

type SysDict struct {
	Model
	Name     string        `json:"name" gorm:"column:name;comment:字典名称"`
	IsSystem bool          `json:"isSystem" gorm:"column:is_system;type:tinyint(1);;comment:系统字典"`
	Type     string        `json:"type" gorm:"column:type;comment:标记"`
	Remark   string        `json:"remark" gorm:"column:remark;comment:备注"`
	Items    []SysDictItem `json:"items" gorm:"foreignKey:DictId;comment:字典项"`
}

func (SysDict) TableName() string {
	return "sys_dict"
}
