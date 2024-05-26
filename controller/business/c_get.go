package business

import (
	"fmt"
	"github.com/front-ck996/csy"
	"github.com/gin-gonic/gin"
	"github.com/rookiefront/api-core/define"
	"github.com/rookiefront/api-core/global"
	"github.com/rookiefront/api-core/model"
	"github.com/rookiefront/api-core/model/manage_api"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func getVerify(c *define.BasicContext, info reqInfo, dbHandle *gorm.DB, reqData map[string]any) error {
	// 数据字段验证
	err := manage_api.VerifyRules(info.tableName, reqData, "where")
	if err != nil {
		c.SendJsonErrCode(err, 20001)
		return err
	}
	// 赛选数据
	filterReqData := manage_api.FilterData(info.tableName, reqData, "where")
	c.Where = filterReqData
	for field, value := range filterReqData {
		dbHandle.Where(fmt.Sprintf("%s = ?", field), value)
	}
	return nil
}

func Get(c *define.BasicContext, info reqInfo, result any, currentModule manage_api.ManageApiModule) {

	reqData := c.GetReqData()
	// 如果内置的字段存在
	if currentModule.TaName2 == "SysDictItem" {
		if _, ok := reqData["dictType"]; ok && model.GetSystemStatus(reqData["dictType"].(string)) != nil {
			c.SendJsonOk(model.GetSystemStatus(reqData["dictType"].(string)))
			return
		}
	}
	dbHandle := global.DB.Table(info.fullTableName)
	switch info.getType {
	// tree 数据懒加载
	case "lazy":
		if !currentModule.PublicWhere {
			if c.VerifyRequestQualification(currentModule.TaName+"_query") != nil {
				return
			}
		}
		if reqData[currentModule.Tree.PID] == nil {
			fieldName := currentModule.GetDbField(currentModule.Tree.PID)
			if fieldName == "" {
				c.SendJsonErr("接口不支持")
				return
			}
			dbHandle.
				Where(fmt.Sprintf("%s = 0", fieldName)).
				Or(fmt.Sprintf("%s is null", fieldName))
		}
		if err := getVerify(c, info, dbHandle, reqData); err != nil {
			c.SendJsonErr(err)
			return
		}

		dbHandle.Find(&result)
		c.SendJsonOk(result)
		break
	//	list 数据无分页
	case "list":
		if c.VerifyRequestQualification(currentModule.TaName+"_query") != nil {
			return
		}
		if err := getVerify(c, info, dbHandle, reqData); err != nil {
			return
		}
		if c.GetCurrentUserId() != "1" {
			reqData["c_id"] = c.GetCurrentUserId()
		}
		dbHandle.Find(&result)
		c.SendJsonOk(result)
		break
	//	分页数据
	case "page":
		if !currentModule.PublicWhere {
			if c.VerifyRequestQualification(currentModule.TaName+"_query") != nil {
				return
			}
		}
		if err := getVerify(c, info, dbHandle, reqData); err != nil {
			return
		}
		pageSize := 10
		pageNum := 1
		if _, ok := reqData["pageNum"]; ok {
			i, err := strconv.ParseInt(reqData["pageNum"].(string), 10, 64)
			if err == nil {
				pageNum = int(i)
			}
		}

		if _, ok := reqData["pageSize"]; ok {
			i, err := strconv.ParseInt(reqData["pageSize"].(string), 10, 64)
			if err == nil {
				pageSize = int(i)
			}
		}
		total := int64(0)
		dbHandle.Where("deleted_at is null").Count(&total)
		if pageNum > 1 {
			dbHandle.Offset((pageNum - 1) * pageSize)
		}
		if pageSize == 0 {
			pageSize = 10
		}
		if pageSize > 150 {
			pageSize = 150
		}
		dbHandle.Limit(pageSize)

		dbHandle.Find(&result)
		c.SendJsonOkPage(result, gin.H{
			"total":    total,
			"pageSize": pageSize,
			"pageNum":  pageNum,
		})
		break
	// 单条数据,主键查询
	case "info":
		if !currentModule.PublicWhere {
			if c.VerifyRequestQualification(currentModule.TaName+"_query") != nil {
				return
			}
		}
		tx := dbHandle.Where("id = ? ", reqData["id"]).First(&result)
		if tx.Error != nil {
			c.SendJsonOk(nil)
		} else {
			c.SendJsonOk(result)
		}
		break
	//	单挑数据,条件查询排除了主键
	case "find":
		if !currentModule.PublicWhere {
			if c.VerifyRequestQualification(currentModule.TaName+"_query") != nil {
				return
			}
		}

		if err := getVerify(c, info, dbHandle, reqData); err != nil {
			c.SendJsonErr(err)
			return
		}
		tx := dbHandle.First(&result)
		if tx.Error != nil {
			c.SendJsonOk(nil)
		} else {
			c.SendJsonOk(result)
		}
		break
	default:
		fieldInfo := currentModule.GetDbFieldInfo(info.getType)
		if fieldInfo.FrontField != "" {
			switch fieldInfo.Associations.Type {
			case "ManyToMany":
				id := c.Query("id")
				m2mTableName := fmt.Sprintf("m2m_%s_%s", currentModule.TaName, csy.StrUpperToSplit(fieldInfo.Associations.Table, "_"))
				var r []map[string]any
				if !strings.Contains(id, ",") {
					global.DB.Table(m2mTableName).Where(fmt.Sprintf("%s_id = ?", currentModule.TaName), reqData["id"]).Find(&r)
				} else {
					global.DB.Table(m2mTableName).Where(fmt.Sprintf("%s_id in (?)", currentModule.TaName), strings.Split(id, ",")).Find(&r)
				}
				linkModel := manage_api.GetModuleByTableName2(fieldInfo.Associations.Table)
				ids := []any{}
				for _, m := range r {
					if _, ok := reqData["map"]; ok {
						_data := map[string]any{}
						_data["id"] = m[fmt.Sprintf("%s_id", currentModule.TaName)]
						if currentModule.TaName != linkModel.TaName {
							_data["dataId"] = m[fmt.Sprintf("%s_id", linkModel.TaName)]
						} else {
							_data["dataId"] = m[fmt.Sprintf("%s_id", fieldInfo.DbField)]
						}
						ids = append(ids, _data)
					} else {
						ids = append(ids, m[fmt.Sprintf("%s_id", linkModel.TaName)])
					}
				}
				c.SendJsonOk(ids)
				return
			case "HasMany":
				//linkModel := manage_api.GetModuleByTableName2(fieldInfo.Associations.Table)
				//r, err := model.GetModel(linkModel.TaName, true)
				//if err != nil {
				//	c.SendJsonErr(err)
				//	return
				//}
				//where := map[string]any{}
				//for _, v := range linkModel.Fields {
				//	if v.FrontField == fieldInfo.Associations.Field {
				//		where = map[string]any{
				//			v.DbField: reqData["id"],
				//		}
				//		break
				//	}
				//}
				//if len(where) == 0 {
				//	c.SendJsonErrCode("接口参数错误", 9000003)
				//	return
				//}
				//global.DB.Model(r).Where(where).Find(&r)
				//c.SendJsonOk(r)
				return
			case "Link":
				linkModule, err := model.GetModel(fieldInfo.Associations.Module, false)
				if err != nil {
					c.SendJsonErr(err)
					return
				}
				global.DB.Model(linkModule).Where(fmt.Sprintf("%s = ?", currentModule.GetDbField(fieldInfo.Associations.Field)), reqData["id"]).First(&linkModule)
				c.SendJsonOk(linkModule)
				return
			}
			c.SendJsonErrCode("接口参数错误", 9000002)
			return
		}
		c.SendJsonErrCode("接口不存在", 20099)
	}
}
