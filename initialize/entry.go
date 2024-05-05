package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/rookiefront/api-core/config"
	"github.com/rookiefront/api-core/global"
	"github.com/rookiefront/api-core/router"
)

func Init() {
	// 加载配置文件
	config.LoadConfig()
	// 链接数据库
	config.DbConnect()
	r := gin.Default()
	global.Engine = r
	currentConfig := config.GetConfig()
	r.Static(currentConfig.System.StaticPreFix, "./public/"+currentConfig.System.StaticDir)

	// 注册路由
	router.Register()
	//
	//global.DB.AutoMigrate(
	//	manage_api.ManageApiModule{},
	//	manage_api.ManageApiModuleField{},
	//	model.SysUser{},
	//	model.SysCity{},
	//)
	//// 开放端口
	//r.Run(":8081")
}
