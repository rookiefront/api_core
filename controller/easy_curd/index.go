package easy_curd

import (
	"fmt"
	"github.com/rookiefront/api-core/define"
	"reflect"
)

type reqInfo struct {
	tableName string
	primaryId string
	getType   string
	pathSplit []string
}

var dbFields = map[string]map[string]string{}
var dbFieldArray = []map[string]string{
	{
		"indexName": "tree_parent",
		"code":      "model.PrimarykeyType",
	},
	{
		"indexName": "icon",
	},
	{
		"indexName": "enable",
		"code":      "int",
		"db":        "type:tinyint(1);default:2",
	},
	{
		"indexName": "bool",
		"code":      "bool",
		"db":        "type:tinyint(1);",
	},
	{
		"indexName": "time",
		"code":      "time.Time",
	},
	{
		"indexName": "date-time",
		"code":      "time.Time",
	},
	{
		"indexName": "image_only",
	},
	{
		"indexName": "images",
		"code":      "model.DataJSONArray",
		"db":        "type:longtext",
	},
	{
		"indexName": "text",
		"db":        "type:text",
	},
	{
		"indexName": "text-no-row",
		"db":        "type:text",
	},
	{
		"indexName": "json",
		"code":      "model.DataJSON",
		"db":        "type:longtext",
	},
	{
		"indexName": "int",
		"code":      "int",
	},
	{
		"indexName": "city_province",
		"code":      "int",
	},
	{
		"indexName": "city_city",
		"code":      "int",
	},
	{
		"indexName": "city_area",
		"code":      "int",
	},
	{
		"indexName": "uuid",
	},
}

func init() {
	for _, m := range dbFieldArray {
		dbFields[m["indexName"]] = m
	}
}

var reqFieldVerify = []map[string]string{
	{
		"label": "必填项",
		"value": "required",
	},
	{
		"label": "输入",
		"value": "other",
	},
}

func getId(value any) string {
	if value == nil {
		return ""
	}
	of := reflect.TypeOf(value).Name()
	switch of {
	case "string":
		return value.(string)
	case "int":
		return fmt.Sprintf("%d", value)
	}
	return fmt.Sprintf("%v", value)
}

func GetReqVerify(c *define.BasicContext) {
	c.SendJsonOk(reqFieldVerify)
}

func GetDbFields(c *define.BasicContext) {
	var result []string
	for _, v := range dbFieldArray {
		result = append(result, v["indexName"])
	}
	c.SendJsonOk(result)
}
