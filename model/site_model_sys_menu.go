package model

type SysMenu struct {
	Model
	MenuName  string         `json:"menuName" gorm:"column:menu_name;comment:菜单名称"`
	ParentId  PrimarykeyType `json:"parentId" gorm:"column:parent_id;comment:父级"`
	Show      int            `json:"show" gorm:"column:show;type:tinyint(1);default:2;comment:显示"`
	Module    string         `json:"module" gorm:"column:module;comment:模块"`
	Component string         `json:"component" gorm:"column:component;type:text;comment:组件"`
	Sort      int            `json:"sort" gorm:"column:sort;comment:排序"`
	Icon      string         `json:"icon" gorm:"column:icon;comment:图标"`
}

func (SysMenu) TableName() string {
	return "sys_menu"
}
