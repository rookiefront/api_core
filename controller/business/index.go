package business

import (
	"fmt"
	"github.com/front-ck996/csy"
	"github.com/rookiefront/api-core/define"
	"github.com/rookiefront/api-core/global"
	"github.com/rookiefront/api-core/model"
	"github.com/rookiefront/api-core/model/manage_api"
	"net/http"
	"reflect"
	"strings"
)

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

type reqInfo struct {
	// 简写
	tableName string
	// 带前缀
	fullTableName string
	primaryId     string
	getType       string
	pathSplit     []string
}

func EasyCURD(c *define.BasicContext) {
	reqData := c.GetReqData()
	rId := getId(reqData["id"])
	req := reqInfo{
		primaryId: rId,
	}
	req.tableName = strings.Trim(c.Param("table"), "/")
	req.getType = strings.Trim(c.Param("path"), "/")
	if !manage_api.IsRequestModule(req.tableName) {
		c.SendJsonErrCode("模块不存在", 10001)
		return
	}
	if !manage_api.GetModule(req.tableName).Enable {
		c.SendJsonErrCode("模块未启用", 10002)
		return
	}
	currentModule := manage_api.GetModule(req.tableName)

	req.fullTableName = currentModule.TaName
	switch c.Request.Method {
	case http.MethodGet:
		if !currentModule.Where {
			c.SendJsonErr("接口未激活")
			return
		}

		isArray := false
		if csy.SliceInclude[string]([]string{"list", "page", "lazy"}, req.getType) {
			isArray = true
		}
		dataStruct, err := model.GetModel(currentModule.TaName, isArray)
		if err != nil {
			c.SendJsonErrCode("模块不存在", 10003)
			return
		}
		Get(c, req, dataStruct, currentModule)
		return

	case http.MethodPost:
		if !currentModule.Create {
			c.SendJsonErr("接口未激活")
			return
		}
		Insert(c, req, reqData, currentModule)
		break
	case http.MethodPut:
		if !currentModule.Update {
			c.SendJsonErr("接口未激活")
			return
		}
		Update(c, req, reqData, currentModule)
		return
	case http.MethodDelete:
		if c.VerifyRequestQualification(currentModule.TaName+"_delete") != nil {
			return
		}
		if !currentModule.Delete {
			c.SendJsonErr("接口未激活")
			return
		}
		global.DB.Table(currentModule.TaName).Where("id = ?", rId).Delete(&model.Model{})
		c.SendJsonOk("ok")
		break
	}
}
