package model

type SysUser struct {
	Model
	UserName              string    `json:"userName" gorm:"column:user_name;comment:用户名"`
	NickName              string    `json:"nickName" gorm:"column:nick_name;comment:昵称"`
	Password              string    `json:"password" gorm:"column:password;comment:密码"`
	Enable                int       `json:"enable" gorm:"column:enable;type:tinyint(1);default:2;comment:激活"`
	BlindDateMember       any       `json:"blindDateMember" gorm:"-;comment:相亲信息"`
	BlindDateMemberEnable int       `json:"blindDateMemberEnable" gorm:"column:blind_date_member_enable;type:tinyint(1);default:2;comment:相亲信息激活"`
	RoleList              []SysRole `json:"roleList" gorm:"many2many:m2m_sys_user_sys_role;comment:角色"`
	Sign                  string    `json:"sign" gorm:"column:sign;comment:签名"`
}

func (SysUser) TableName() string {
	return "sys_user"
}
