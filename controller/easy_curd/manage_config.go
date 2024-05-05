package easy_curd

import (
	"github.com/rookiefront/api-core/cmd/model"
	"github.com/rookiefront/api-core/define"
	"github.com/rookiefront/api-core/global"
	"github.com/rookiefront/api-core/model/manage_api"
	"sort"
	"strings"
)

func ManageConfig(c *define.BasicContext) {
	data := c.GetReqData()
	result := manage_api.ManageApiModule{}
	if _, ok := data["module"]; ok {
		global.DB.Model(result).Where("url = ? ", data["module"]).Preload("Fields").Find(&result)
	}
	sort.Slice(result.Fields, func(i, j int) bool {
		return result.Fields[i].ApiManageSort < result.Fields[j].ApiManageSort
	})
	c.SendJsonOk(result)
}

func Tables(c *define.BasicContext) {
	var list []manage_api.ManageApiModule
	global.DB.Model(manage_api.ManageApiModule{}).Preload("Fields").Find(&list)
	var newList []manage_api.ManageApiModule
	for _, v := range list {
		if !strings.Contains(v.TaName, "manage_api") {
			newList = append(newList, v)
		}
	}
	c.SendJsonOk(newList)
}

func StepList(c *define.BasicContext) {
	list := []model.SysStep{}
	global.DB.Where("parent_id = ?", 0).Or("parent_id is null").Find(&list)
	c.SendJsonOk(list)
}
