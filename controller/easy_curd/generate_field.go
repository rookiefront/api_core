package easy_curd

import (
	"errors"
	"github.com/rookiefront/api-core/define"
	"github.com/rookiefront/api-core/global"
	"github.com/rookiefront/api-core/model/manage_api"
)

// 字段设置
type ApiManageApiField struct {
}

func (api ApiManageApiField) List(c *define.BasicContext) {
	req := c.GetReqData()
	var result []manage_api.ManageApiModuleField
	global.DB.Where(req).Find(&result)
	c.SendJsonOk(result)
}

func (api ApiManageApiField) Insert(c *define.BasicContext) {
	req := manage_api.ManageApiModuleField{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.SendJsonErr(err)
		return
	}

	tx := global.DB.Save(&req)
	if tx.Error != nil {
		c.SendJsonErr(tx.Error)
		return
	}
	c.SendJsonOk()
}

func (api ApiManageApiField) Update(c *define.BasicContext) {
	req := manage_api.ManageApiModuleField{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.SendJsonErr(err)
		return
	}
	if req.GetId() == 0 {
		c.SendJsonErr(errors.New("0"))
		return
	}
	tx := global.DB.Save(&req)
	if tx.Error != nil {
		c.SendJsonErr(tx.Error)
		return
	}
	c.SendJsonOk()
}

func (api ApiManageApiField) Delete(c *define.BasicContext) {
	id := c.Query("id")
	global.DB.Table(manage_api.ManageApiModuleField{}.TableName()).Where("id = ?", id).Delete(&manage_api.ManageApiModuleField{})
	c.SendJsonOk()
}
