package model_test

import (
	"fmt"
	"github.com/rookiefront/api-core/cmd/model"
	"github.com/rookiefront/api-core/config"
	"github.com/rookiefront/api-core/define"
	"github.com/rookiefront/api-core/global"
	"os"
	"testing"
)

func TestGenerateModel(t *testing.T) {
	args := os.Args
	tableName := args[len(args)-1]
	//tableName = "sys_step"
	//go test -v cmd/model/generate_test.go --args daily_note
	fmt.Println("tableName", tableName, args[len(args)-1])
	getModel, err := model.GetModel(tableName, false)
	var result define.ResultJSON
	if err != nil {
		result.Err = err
		fmt.Println(result.ToTagJSON())
		return
	}
	config.LoadConfig()
	config.DbConnect()
	err = global.DB.AutoMigrate(&getModel)
	if err != nil {
		result.Err = err
		fmt.Println(result.ToTagJSON())
		return
	}
	fmt.Println(result.ToTagJSON())
}
