package model

type BlindDateMemberInterest struct {
	Model
	UserId      int    `json:"userId" gorm:"column:user_id;comment:用户ID"`
	Label       string `json:"label" gorm:"column:label;comment:问题"`
	Value       string `json:"value" gorm:"column:value;comment:答案"`
	LabelCustom string `json:"labelCustom" gorm:"column:label_custom;comment:自定义提问"`
}

func (BlindDateMemberInterest) TableName() string {
	return "blind_date_member_interest"
}
