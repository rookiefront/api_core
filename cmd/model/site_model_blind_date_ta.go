package model

import ()

type BlindDateTa struct {
	Model
	StartAge           int    `json:"startAge" gorm:"column:start_age;comment:开始年龄"`
	StopAge            int    `json:"stopAge" gorm:"column:stop_age;comment:结束年龄"`
	StartHeight        int    `json:"startHeight" gorm:"column:start_height;comment:开始身高"`
	StopHeight         int    `json:"stopHeight" gorm:"column:stop_height;comment:结束身高"`
	StartMonthlyIncome int    `json:"startMonthlyIncome" gorm:"column:start_monthly_income;comment:开始月收入"`
	StopMonthlyIncome  int    `json:"stopMonthlyIncome" gorm:"column:stop_monthly_income;comment:结束月收入"`
	Education          string `json:"education" gorm:"column:education;comment:学历"`
	Marriage           string `json:"marriage" gorm:"column:marriage;comment:婚姻状态"`
	ChildrenStatus     string `json:"childrenStatus" gorm:"column:children_status;comment:子女"`
	MarkChildren       string `json:"markChildren" gorm:"column:mark_children;comment:多久要孩子"`
	House              string `json:"house" gorm:"column:house;comment:房"`
	Car                string `json:"car" gorm:"column:car;comment:车"`
	Smoking            string `json:"smoking" gorm:"column:smoking;comment:烟"`
	Drink              string `json:"drink" gorm:"column:drink;comment:酒"`
	GiveTaDesc         string `json:"giveTaDesc" gorm:"column:give_ta_desc;type:text;comment:给Ta的一句话"`
	UserId             int    `json:"userId" gorm:"column:user_id;comment:UserID"`
}

func (BlindDateTa) TableName() string {
	return "blind_date_ta"
}
