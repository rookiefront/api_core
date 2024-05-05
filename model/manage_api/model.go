package manage_api

import (
	model2 "github.com/rookiefront/api-core/cmd/model"
	"github.com/rookiefront/api-core/model"
)

type ManageApiModule struct {
	model.Model
	Name        string                   `json:"name" gorm:"column:name;comment:模块名"`
	Group       string                   `json:"group" gorm:"column:group;comment:表分组"`
	TaName      string                   `json:"table_name" gorm:"column:table_name;comment:数据库表名"`
	TaName2     string                   `json:"table_name_2" gorm:"column:table_name_2"`
	Url         string                   `json:"url" gorm:"comment:路由地址"`
	Enable      bool                     `json:"enable" gorm:"comment:是否启用"`
	Fields      []ManageApiModuleField   `json:"fields" gorm:"fields"`
	Tree        ManageApiModuleFieldTree `json:"tree" gorm:"column:tree;type:longtext"`
	Create      bool                     `json:"create"`
	Delete      bool                     `json:"delete"`
	Update      bool                     `json:"update"`
	Where       bool                     `json:"where"`
	DialogWidth int                      `json:"dialog_width"`

	TableOperateBtn model2.DataJSONArray `json:"table_operate_btn" gorm:"column:table_operate_btn;type:longtext"`
}

func (ManageApiModule) TableName() string {
	return "manage_api_module"
}

func (data ManageApiModule) GetDbField(s string) string {
	fieldName := ""
	for _, field := range data.Fields {
		if field.FrontField == s {
			fieldName = field.DbField
		}
	}
	return fieldName
}
func (data ManageApiModule) GetDbFieldInfo(s string) ManageApiModuleField {
	var d ManageApiModuleField
	for _, field := range data.Fields {
		if field.FrontField == s {
			d = field
		}
	}
	return d
}

type ManageApiModuleField struct {
	model.Model
	FrontField           string `json:"front_field" gorm:"comment:前端请求的字段名"`
	FrontFieldGroup      string `json:"front_field_group" gorm:"comment:前端字段分组"`
	FrontFieldGroupRatio int    `json:"front_field_group_ratio" gorm:"comment:前端字段分组比例"`
	FrontFieldGroupSort  int    `json:"front_field_group_sort" gorm:"comment:前端字段分组排序"`

	// 前端组件
	Component         ManageApiModuleFieldComponent `json:"component"`
	DbField           string                        `json:"db_field" gorm:"comment:数据库字段名"`
	DbFieldType       string                        `json:"db_field_type" gorm:"comment:数据库字段类型"`
	Dict              string                        `json:"dict"`
	Step              string                        `json:"step"`
	Comment           string                        `json:"comment" gorm:"size:0"`
	CommentDesc       string                        `json:"comment_desc" gorm:"size:0"`
	Where             bool                          `json:"where"`
	WhereVerify       string                        `json:"where_verify"`
	Insert            bool                          `json:"insert"`
	InsertVerify      string                        `json:"insert_verify"`
	Update            bool                          `json:"update"`
	UpdateVerify      string                        `json:"update_verify"`
	ManageApiModuleID uint                          `json:"manage_api_module_id"`
	SortBy            string                        `json:"sort_by" gorm:"comment:排序字段"`
	SortOrder         bool                          `json:"sort_order" gorm:"comment:排序顺序 asc desc"`
	HideForm          bool                          `json:"hide_form"`
	// 表格显示
	Table bool `json:"table" gorm:"comment:表格显示"`
	// 表格显示
	TableWhere bool `json:"table_where" gorm:"comment:表格下拉查询"`
	// 关联一对一, 多对多关系
	Associations ManageApiModuleFieldAssociations `json:"associations"`

	// 内部接口管理,排序
	ApiManageSort int `json:"api_manage_sort"`
}

func (ManageApiModuleField) TableName() string {
	return "manage_api_module_field"
}
