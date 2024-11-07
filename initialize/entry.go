package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rookiefront/api-core/config"
	"github.com/rookiefront/api-core/global"
	"github.com/rookiefront/api-core/model"
	"github.com/rookiefront/api-core/model/manage_api"
	"github.com/rookiefront/api-core/router"
)

func Init(callBack func()) {
	// 加载配置文件
	config.LoadConfig()
	// 链接数据库
	config.DbConnect()
	r := gin.Default()
	global.Engine = r
	currentConfig := config.GetConfig()
	r.Static(currentConfig.System.StaticPreFix, "./public/"+currentConfig.System.StaticDir+"/../")

	// 注册路由
	router.Register()

	global.DB.AutoMigrate(
		manage_api.ManageApiModule{},
		manage_api.ManageApiModuleField{},
	)

	tableList := []string{
		"sys_user",
		"sys_city",
		"sys_step",
	}
	for _, i := range tableList {
		getModel, err := model.GetModel(i, false)
		if err != nil {
			global.DB.AutoMigrate(getModel)
		}
	}
	callBack()
	// 开放端口
	r.Run(fmt.Sprintf("%s:%d", currentConfig.System.Host, currentConfig.System.Port))
}
