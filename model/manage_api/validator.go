package manage_api

import (
	"github.com/rookiefront/api-core/config"
	"github.com/rookiefront/api-core/global"
	"github.com/rookiefront/api-core/utils/common"
	"sync"
)

type ManageApiValidator struct {
	Field string
	Rules []string
}
type ManageApiValidatorRule struct {
}

var tables map[string]ManageApiModule
var lock sync.Mutex

func ReloadTables() {
	tables = map[string]ManageApiModule{}
	tables = getModuleList()
}

func getModuleList() map[string]ManageApiModule {
	var result []ManageApiModule
	resultMap := map[string]ManageApiModule{}
	global.DB.Preload("Fields").Find(&result)
	for _, module := range result {
		resultMap[module.Url] = module
	}
	return resultMap
}

func IsRequestModule(table string) bool {
	if config.IsDev() || len(tables) == 0 {
		lock.Lock()
		tables = getModuleList()
		lock.Unlock()
	}
	if _, ok := tables[table]; ok {
		return true
	}
	return false
}

func GetModule(table string) ManageApiModule {
	var resutl ManageApiModule
	if _, ok := tables[table]; ok {
		resutl = tables[table]
	}
	return resutl
}
func GetModuleByTableName2(table string) ManageApiModule {
	var resutl ManageApiModule
	for _, module := range tables {
		if module.TaName2 == table {
			resutl = module
			break
		}
	}
	return resutl
}

func VerifyRules(tableName string, reqData map[string]interface{}, verifyType string) error {

	table := tables[tableName]
	for _, field := range table.Fields {
		continueCurrent := false
		currentValue := reqData[field.FrontField]
		if currentValue == nil {
			currentValue = ""
		}
		currentVerify := ""
		switch verifyType {
		case "where":
			if !field.Where || field.WhereVerify == "" {
				continueCurrent = true
			}
			currentVerify = field.WhereVerify
			break
		}
		if continueCurrent || currentVerify == "" {
			continue
		}
		return common.Validate.Var(currentValue, currentVerify, field.FrontField)

	}
	return nil
}

func FilterData(tableName string, reqData map[string]interface{}, verifyType string) map[string]interface{} {
	table := tables[tableName]
	filterReqData := map[string]interface{}{}
	for _, field := range table.Fields {
		if reqData[field.FrontField] == nil {
			continue
		}
		switch {
		case verifyType == "where" && field.Where:
			filterReqData[field.DbField] = reqData[field.FrontField]
			break
		}
	}

	return filterReqData
}
