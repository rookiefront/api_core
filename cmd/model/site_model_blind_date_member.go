package model

import (
	"time"
)

type BlindDateMember struct {
	Model
	Avatar         string                    `json:"avatar" gorm:"column:avatar;comment:头像"`
	UserId         int                       `json:"userId" gorm:"column:user_id;comment:UserId"`
	Area           int                       `json:"area" gorm:"column:area;comment:贯籍_区县"`
	House          string                    `json:"house" gorm:"column:house;comment:房"`
	City           int                       `json:"city" gorm:"column:city;comment:贯籍_市"`
	NickName       string                    `json:"nickName" gorm:"column:nick_name;comment:昵称"`
	Weight         int                       `json:"weight" gorm:"column:weight;comment:体重"`
	Height         int                       `json:"height" gorm:"column:height;comment:身高"`
	Birth          time.Time                 `json:"birth" gorm:"column:birth;comment:出生年月"`
	Married        string                    `json:"married" gorm:"column:married;comment:多久结婚"`
	Ta             any                       `json:"ta" gorm:"-;comment:心仪的他"`
	Nation         string                    `json:"nation" gorm:"column:nation;comment:民族"`
	Occupation     string                    `json:"occupation" gorm:"column:occupation;comment:职业"`
	Province       int                       `json:"province" gorm:"column:province;comment:籍贯_省"`
	Car            string                    `json:"car" gorm:"column:car;comment:车"`
	Smoking        string                    `json:"smoking" gorm:"column:smoking;comment:烟"`
	Drink          string                    `json:"drink" gorm:"column:drink;comment:酒"`
	Soliloquy      string                    `json:"soliloquy" gorm:"column:soliloquy;type:text;comment:内心独白"`
	MarkChildren   string                    `json:"markChildren" gorm:"column:mark_children;comment:多久要孩子"`
	Education      string                    `json:"education" gorm:"column:education;comment:学历"`
	ChildrenStatus string                    `json:"childrenStatus" gorm:"column:children_status;comment:子女"`
	FollowList     []BlindDateMember         `json:"followList" gorm:"many2many:m2m_blind_date_member_blind_date_member;comment:关注列表"`
	Interest       []BlindDateMemberInterest `json:"interest" gorm:"many2many:m2m_blind_date_member_blind_date_member_interest;comment:兴趣爱好"`
}

func (BlindDateMember) TableName() string {
	return "blind_date_member"
}
