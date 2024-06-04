package easy_curd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/front-ck996/csy"
	"github.com/rookiefront/api-core/define"
	"github.com/rookiefront/api-core/global"
	"github.com/rookiefront/api-core/model/manage_api"
	"os"
	"regexp"
	"sort"
	"strings"
)

func GenerateBusTableApi(c *define.BasicContext) {
	data := c.GetReqData()
	if data["id"] != "" {
		if err := generateBusTable(data["id"]); err != nil {
			c.SendJsonErr(err)
			return
		}
		c.SendJsonOk("ok")
	} else {
		c.SendJsonErr(errors.New("id 必传"))
	}

}
func generateBusTable(id any) error {
	var findData = manage_api.ManageApiModule{}
	global.DB.Table("manage_api_module").Where("id = ?", id).Preload("Fields").First(&findData)
	sort.Slice(findData.Fields, func(i, j int) bool {
		return findData.Fields[i].ApiManageSort < findData.Fields[j].ApiManageSort
	})
	if findData.GetId() == 0 {
		return errors.New("id 未找到")
	}
	// 生成模型字符串
	fields := []string{}
	existTimeField := ""
	for _, v := range findData.Fields {
		codeType := "string"
		dbAppend := []string{
			"column:" + v.DbField,
		}
		if v.DbFieldType == "time" {
			existTimeField = "\"time\""
		}
		if _, ok := dbFields[v.DbFieldType]; ok {
			if dbFields[v.DbFieldType]["code"] != "" {
				codeType = dbFields[v.DbFieldType]["code"]
			}
			if dbFields[v.DbFieldType]["db"] != "" {
				dbAppend = append(dbAppend, dbFields[v.DbFieldType]["db"])
			}
		}

		if v.Associations.Type == "HasMany" {
			codeType = "[]" + v.Associations.Table
			dbAppend = []string{
				"foreignKey:" + csy.StrCapitalize(csy.StrFirstToUpper(v.Associations.Field)),
			}
		}

		if v.Associations.Type == "HasOne" {
			codeType = v.Associations.Table
			dbAppend = []string{
				"foreignKey:" + csy.StrCapitalize(csy.StrFirstToUpper(v.Associations.Field)),
			}
		}

		if v.Associations.Type == "ManyToMany" {
			codeType = "[]" + v.Associations.Table
			dbAppend = []string{
				"many2many:m2m_" + findData.TaName + "_" + csy.StrUpperToSplit(v.Associations.Table, "_"),
			}
		}

		if v.Associations.Type == "BelongsTo" {
			fields = append(fields,
				fmt.Sprintf("%sID %s `json:\"%sId\" gorm:\"%s;comment:%s\"`",
					csy.StrCapitalize(csy.StrFirstToUpper(v.DbField)),
					codeType,
					v.FrontField,
					strings.Join(dbAppend, ";"),
					v.Comment),
			)
			fields = append(fields,
				fmt.Sprintf("%s %s `json:\"%s\"`",
					csy.StrCapitalize(csy.StrFirstToUpper(v.DbField)),
					v.Associations.Table,
					v.FrontField,
				),
			)

		} else if v.Associations.Type == "Link" {
			codeType = "any"
			fields = append(fields,
				fmt.Sprintf("%s %s `json:\"%s\" gorm:\"-;comment:%s\"`",
					csy.StrCapitalize(csy.StrFirstToUpper(v.DbField)),
					codeType,
					v.FrontField,
					v.Comment),
			)
		} else {
			fields = append(fields,
				fmt.Sprintf("%s %s `json:\"%s\" gorm:\"%s;comment:%s\"`",
					csy.StrCapitalize(csy.StrFirstToUpper(v.DbField)),
					codeType,
					v.FrontField,
					strings.Join(dbAppend, ";"),
					v.Comment),
			)
		}

	}
	modelTemplateStr := fmt.Sprintf(`
package model

import (
	%s
	"github.com/rookiefront/api-core/model"
)

type %s struct {
		model.Model
	%s
}

func (%s) TableName() string {
	return "%s"
}


type %s struct {
	%s
}

func (%s) TableName() string {
	return "%s"
}
`,
		existTimeField,

		csy.StrCapitalize(csy.StrFirstToUpper(findData.TaName)),
		strings.Join(fields, "\r\n"),
		csy.StrCapitalize(csy.StrFirstToUpper(findData.TaName)),
		findData.TaName,
		// 模型名称
		"NoModel"+csy.StrCapitalize(csy.StrFirstToUpper(findData.TaName)),
		strings.Join(fields, "\r\n"),

		// 表名
		"NoModel"+csy.StrCapitalize(csy.StrFirstToUpper(findData.TaName)),
		findData.TaName,
	)

	// 格式化代码
	code, err2 := csy.FormatCode([]byte(modelTemplateStr))
	if err2 == nil {
		modelTemplateStr = string(code)
	}
	// 生成数据库模型
	err := os.WriteFile("cmd/model/site_model_"+findData.TaName+".go", []byte(modelTemplateStr), 0644)
	if err != nil {
		return err
	}

	var tables []manage_api.ManageApiModule
	global.DB.Find(&tables)
	hashAppend := []string{}
	for _, v := range tables {
		hashAppend = append(hashAppend, fmt.Sprintf("hash[\"%s\"] = &%s{}", v.TaName, csy.StrCapitalize(csy.StrFirstToUpper(v.TaName))))
		hashAppend = append(hashAppend, fmt.Sprintf("hash[\"%s_arr\"] = []%s{}", v.TaName, csy.StrCapitalize(csy.StrFirstToUpper(v.TaName))))
	}
	//写入hash
	modelHashTemplateStr := fmt.Sprintf(`
package model
import "errors"
func GetModel(tableName string,isArray bool) (interface{}, error) {
	if isArray {
		tableName += "_arr"
	}
	hash := map[string]any{}
	%s
	if _, ok := hash[tableName]; ok {
		return hash[tableName], nil
	}
	return "", errors.New("未找到")
}

func GetHash() map[string]any {
	hash := map[string]any{}
	%s
	return hash
}
`, strings.Join(hashAppend, "\r\n"), strings.Join(hashAppend, "\r\n"))
	// 格式化代码
	code, err2 = csy.FormatCode([]byte(modelHashTemplateStr))
	if err2 == nil {
		modelHashTemplateStr = string(code)
	}

	err = os.WriteFile("cmd/model/hash.go", []byte(modelHashTemplateStr), 0644)
	if err != nil {
		return err
	}
	cmd := csy.NewCMD()
	_, err = cmd.Run([]string{fmt.Sprintf("go test -v cmd/model/generate_test.go --args %s", findData.TaName)})
	if err != nil {
		return err
	}
	var cmdTestResult define.ResultJSON
	submatch := regexp.MustCompile(`result=>(.*?)<=result`).FindStringSubmatch(cmd.StdoutText)
	if len(submatch) != 2 {
		return errors.New("生成数据库模型,没有返回值")
	}
	json.Unmarshal([]byte(submatch[1]), &cmdTestResult)
	return cmdTestResult.Err
}
