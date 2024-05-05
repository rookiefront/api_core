package model

type SysStep struct {
	Model
	Name     string         `json:"name" gorm:"column:name;comment:名称"`
	Show     int            `json:"show" gorm:"column:show;type:tinyint(1);default:2;comment:显示"`
	FullPath string         `json:"fullPath" gorm:"column:full_path;type:text;comment:完整路径"`
	Leaf     bool           `json:"leaf" gorm:"column:leaf;type:tinyint(1);;comment:子集"`
	Value    string         `json:"value" gorm:"column:value;comment:值"`
	Remark   string         `json:"remark" gorm:"column:remark;type:text;comment:备注"`
	ParentId PrimarykeyType `json:"parentId" gorm:"column:parent_id;comment:父级"`
}

func (SysStep) TableName() string {
	return "sys_step"
}
