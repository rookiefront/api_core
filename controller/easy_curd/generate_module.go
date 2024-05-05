package easy_curd

import (
	"errors"
	"github.com/front-ck996/csy"
	"github.com/rookiefront/api-core/define"
	"github.com/rookiefront/api-core/global"
	"github.com/rookiefront/api-core/model/manage_api"
)

// 字段设置
type ApiManageApiModule struct {
}

func (api ApiManageApiModule) List(c *define.BasicContext) {
	req := c.GetReqData()
	var result []manage_api.ManageApiModule
	global.DB.Where(req).Find(&result)
	c.SendJsonOk(result)
}

func (api ApiManageApiModule) Insert(c *define.BasicContext) {
	req := manage_api.ManageApiModule{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.SendJsonErr(err)
		return
	}
	req.TaName2 = csy.StrCapitalize(csy.StrFirstToUpper(req.TaName))
	tx := global.DB.Save(&req)
	if tx.Error != nil {
		c.SendJsonErr(tx.Error)
		return
	}
	c.SendJsonOk()
}

func (api ApiManageApiModule) Update(c *define.BasicContext) {
	req := manage_api.ManageApiModule{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.SendJsonErr(err)
		return
	}
	if req.GetId() == 0 {
		c.SendJsonErr(errors.New("0"))
		return
	}
	req.TaName2 = csy.StrCapitalize(csy.StrFirstToUpper(req.TaName))
	tx := global.DB.Save(&req)
	if tx.Error != nil {
		c.SendJsonErr(tx.Error)
		return
	}
	c.SendJsonOk()
}

func (api ApiManageApiModule) Delete(c *define.BasicContext) {
	id := c.Query("id")
	global.DB.Table(manage_api.ManageApiModule{}.TableName()).Where("id = ?", id).Delete(&manage_api.ManageApiModule{})
	c.SendJsonOk()
}
