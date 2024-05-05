package business

import (
	"encoding/json"
	"fmt"
	"github.com/front-ck996/csy"
	"github.com/google/uuid"
	"github.com/rookiefront/api-core/define"
	"github.com/rookiefront/api-core/global"
	"github.com/rookiefront/api-core/model/manage_api"
	"time"
)

func Update(c *define.BasicContext, info reqInfo, reqData map[string]any, currentModule manage_api.ManageApiModule) {
	insertData := map[string]any{}

	fieldInfo := currentModule.GetDbFieldInfo(info.getType)
	if fieldInfo.FrontField != "" {
		switch fieldInfo.Associations.Type {
		case "ManyToMany":
			m2mTableName := fmt.Sprintf("m2m_%s_%s", currentModule.TaName, csy.StrUpperToSplit(fieldInfo.Associations.Table, "_"))
			global.DB.Table(m2mTableName).Where(fmt.Sprintf("%s_id = ?", currentModule.TaName), reqData["id"]).Delete(map[string]any{})
			insertM2mData := []map[string]interface{}{}
			module := manage_api.GetModule(fieldInfo.Associations.Module)
			if module.TaName == "" {
				return
			}

			var ids []interface{}

			marshal, _ := json.Marshal(reqData["value"])
			err := json.Unmarshal(marshal, &ids)
			if err != nil {
				return
			}
			for i := 0; i < len(ids); i++ {
				_inertData := map[string]interface{}{}

				_inertData[fmt.Sprintf("%s_id", currentModule.TaName)] = reqData["id"]
				if currentModule.TaName != module.TaName {
					_inertData[fmt.Sprintf("%s_id", module.TaName)] = ids[i]
				} else {
					_inertData[fmt.Sprintf("%s_id", fieldInfo.DbField)] = ids[i]
				}
				insertM2mData = append(insertM2mData, _inertData)
			}

			if len(insertM2mData) >= 1 {
				tx := global.DB.Table(m2mTableName).Create(&insertM2mData)
				if tx.Error != nil {
					c.SendJsonErrCode("接口参数错误", 9000001)
					return
				}
			}

			c.SendJsonOk()
			return
		case "HasMany":
			//linkModel := manage_api.GetModuleByTableName2(fieldInfo.Associations.Table)
			//var ids []interface{}
			//
			//marshal, _ := json.Marshal(reqData["value"])
			//err := json.Unmarshal(marshal, &ids)
			//if err != nil {
			//	c.SendJsonErr(err)
			//	return
			//}
			//
			//dbFieldInfo := linkModel.GetDbFieldInfo(fieldInfo.Associations.Field)
			//if dbFieldInfo.FrontField == "" {
			//	c.SendJsonErrCode("接口参数错误", 9000002)
			//	return
			//}
			//tx := global.DB.Table(linkModel.TaName).Where("id in (?)", ids).Updates(map[string]any{
			//	dbFieldInfo.DbField: reqData["id"],
			//})
			//if tx.Error != nil {
			//	c.SendJsonErr(tx.Error)
			//	return
			//}
			//c.SendJsonOk()
			return
		}
		c.SendJsonErrCode("接口参数错误", 9000002)
		return
	}

	for _, field := range currentModule.Fields {
		if field.DbFieldType == "json" || field.DbFieldType == "images" {
			marshal, err := json.Marshal(reqData[field.FrontField])
			reqData[field.FrontField] = ""
			if err == nil {
				reqData[field.FrontField] = string(marshal)
			}
		}

		if field.DbFieldType == "uuid" {
			if _, ok := reqData["uuid"]; ok {
				if len(reqData["uuid"].(string)) == 0 {
					newUUID, _ := uuid.NewUUID()
					insertData[field.DbField] = newUUID.String()
				}
			}
			continue
		}

		if field.Associations.Type == "BelongsTo" {
			insertData[field.DbField] = reqData[field.FrontField+"Id"]
			continue
		}
		if field.Update {
			if _, ok := reqData[field.FrontField]; ok {
				// 如果是多对多
				if field.Associations.Type == "ManyToMany" || field.Associations.Type == "HasMany" {
					continue
				} else if field.Associations.Type == "Link" {
					dbField := currentModule.GetDbField(field.Associations.Field)
					if field.Associations.CurrentField != "" {
						insertData[dbField] = reqData[field.Associations.CurrentField]
					} else {
						insertData[dbField] = reqData["id"]
					}
				} else {
					insertData[field.DbField] = reqData[field.FrontField]
				}
			}
		}
	}
	if len(insertData) == 0 {
		c.SendJsonErr("未传递数据")
		return
	}

	insertData["updated_at"] = time.Now()
	tx := global.DB.Table(currentModule.TaName).Where("id = ?", reqData["id"]).Updates(&insertData)
	if tx.Error != nil {
		c.SendJsonErr(tx.Error)
		return
	}
	c.SendJsonOk("ok")
}
