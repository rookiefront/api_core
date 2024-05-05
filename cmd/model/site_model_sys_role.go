package model

import ()

type SysRole struct {
	Model
	RoleName       string    `json:"roleName" gorm:"column:role_name;comment:角色名称"`
	Enable         int       `json:"enable" gorm:"column:enable;type:tinyint(1);default:2;comment:启用"`
	PermissionList string    `json:"permissionList" gorm:"column:permission_list;comment:权限字符串列表"`
	MenuList       []SysMenu `json:"menuList" gorm:"many2many:m2m_sys_role_sys_menu;comment:可访问菜单"`
}

func (SysRole) TableName() string {
	return "sys_role"
}
