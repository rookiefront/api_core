package common

import (
	"github.com/rookiefront/api-core/cmd/model"
	"github.com/rookiefront/api-core/define"
	"github.com/rookiefront/api-core/global"
)

type ApiDict struct {
}

func (api *ApiDict) List(c *define.BasicContext) {
	value := c.Query("type")
	var list []model.SysDictItem
	if value == "" {
		c.SendJsonErr("类型输入错误")
		return
	}
	global.DB.Where(model.SysDictItem{DictType: value}).Find(&list)
	c.SendJsonOk(list)
}
