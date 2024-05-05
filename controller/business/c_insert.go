package business

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rookiefront/api-core/define"
	"github.com/rookiefront/api-core/global"
	"github.com/rookiefront/api-core/model/manage_api"
	"time"
)

func Insert(c *define.BasicContext, info reqInfo, reqData map[string]any, currentModule manage_api.ManageApiModule) {
	insertData := map[string]any{}
	for _, field := range currentModule.Fields {
		if field.DbFieldType == "json" || field.DbFieldType == "images" {
			marshal, err := json.Marshal(reqData[field.FrontField])
			reqData[field.FrontField] = ""
			if err == nil {
				reqData[field.FrontField] = string(marshal)
			}
		}
		if field.DbFieldType == "uuid" {
			newUUID, _ := uuid.NewUUID()
			insertData[field.DbField] = newUUID.String()
			continue
		}
		if field.Associations.Type == "BelongsTo" {
			insertData[field.DbField] = reqData[field.FrontField+"Id"]
			continue
		}
		if field.Associations.Type == "Link" {
			//dbField := currentModule.GetDbField(field.Associations.Field)
			//if field.Associations.CurrentField != "" {
			//	insertData[dbField] = reqData[field.Associations.CurrentField]
			//} else {
			//	insertData[dbField] = reqData["id"]
			//}
			continue
		}
		if field.Insert {
			if _, ok := reqData[field.FrontField]; ok {
				insertData[field.DbField] = reqData[field.FrontField]
			}
		}
	}
	if len(insertData) == 0 {
		c.SendJsonErr("未传递数据")
		return
	}

	insertData["c_id"] = c.GetCurrentUserId()
	insertData["created_at"] = time.Now()

	tx := global.DB.Table(currentModule.TaName).Create(insertData)
	if tx.Error != nil {
		c.SendJsonErr(tx.Error)
		return
	}
	c.SendJsonOk("ok")
}
