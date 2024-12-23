package router

import (
	"github.com/front-ck996/csy"
	"github.com/front-ck996/csy/gin_middleware"
	"github.com/gin-gonic/gin"
	config2 "github.com/rookiefront/api-core/config"
	"github.com/rookiefront/api-core/controller/business"
	"github.com/rookiefront/api-core/controller/common"
	"github.com/rookiefront/api-core/controller/easy_curd"
	"github.com/rookiefront/api-core/define"
	"github.com/rookiefront/api-core/global"
	"github.com/rookiefront/api-core/middleware"
	"github.com/rookiefront/api-core/model"
	"github.com/rookiefront/api-core/service"
	"net/http"
)

func Register() {
	config := config2.GetConfig()
	if config.ManageApi.Enable {
		// 接口管理模块
		apiManageApi := global.Engine.Group("/manage_api")
		apiManageApi.Use(gin_middleware.Cors())
		apiManageApi.Use(gin.BasicAuth(gin.Accounts{
			config.ManageApi.Account: config.ManageApi.Password,
		}))
		global.ApiManageRouter = apiManageApi
		apiManageApi.POST("/generate_bus_table", define.WrapHandler(easy_curd.GenerateBusTableApi))
		apiManageApi.POST("/db_field_conv_front", define.WrapHandler(easy_curd.FbFieldConvFront))
		apiManageApi.GET("/db_field_list", define.WrapHandler(easy_curd.GetDbFields))
		apiManageApi.GET("/req_verify", define.WrapHandler(easy_curd.GetReqVerify))
		apiManageApi.GET("/tables", define.WrapHandler(easy_curd.Tables))
		apiManageApi.GET("/step", define.WrapHandler(easy_curd.StepList))

		sApiManageApiField := easy_curd.ApiManageApiField{}
		apiManageApiField := apiManageApi.Group("/field")
		apiManageApiField.GET("/list", define.WrapHandler(sApiManageApiField.List))
		apiManageApiField.PUT("", define.WrapHandler(sApiManageApiField.Update))
		apiManageApiField.POST("", define.WrapHandler(sApiManageApiField.Insert))
		apiManageApiField.DELETE("", define.WrapHandler(sApiManageApiField.Delete))

		sApiManageApiModule := easy_curd.ApiManageApiModule{}
		apiManageApiModule := apiManageApi.Group("/module")
		apiManageApiModule.GET("/list", define.WrapHandler(sApiManageApiModule.List))
		apiManageApiModule.PUT("", define.WrapHandler(sApiManageApiModule.Update))
		apiManageApiModule.POST("", define.WrapHandler(sApiManageApiModule.Insert))
		apiManageApiModule.DELETE("", define.WrapHandler(sApiManageApiModule.Delete))

	}

	// 通用服务模块 && 系统内置
	registerSystemRouter()

	// 业务模块api
	global.ApiPrefix = global.Engine.Group(config.System.ApiPrefix)
	businessApi := global.Engine.Group(config.System.ApiPrefix)
	global.ApiPrefixAuth = businessApi.Use(middleware.InterceptNotLoggedIn)

	if config2.IsDev() {
		businessApi.Use(gin_middleware.Cors())
		global.ApiPrefix.Use(gin_middleware.Cors())
	}
	businessApi.Any("/:table/*path", define.WrapHandler(business.EasyCURD))
	businessApi.Any("/:table", define.WrapHandler(business.EasyCURD))

	if config2.IsDev() {
		rootUser := model.SysUser{
			UserName: "user_root",
		}
		user := model.SysUser{}
		global.DB.Unscoped().Where(rootUser).First(&user)
		if !user.IdTure() {
			userSign := csy.RandomString(8)
			global.DB.Save(&model.SysUser{
				UserName: "user_root",
				NickName: "超级管理员",
				Sign:     userSign,
				Password: service.User.Encrypt("user_root123456user_root123456", userSign),
				Enable:   1,
			})
		}
		initMenu := model.SysMenu{}
		initMenu.CreateUserID = 1
		global.DB.Unscoped().Where(model.SysMenu{
			Module: "menu",
		}).First(&initMenu)
		if initMenu.MenuName == "" {
			initMenu.MenuName = "菜单管理"
		}
		initMenu.Component = "/"
		global.DB.Save(&initMenu)
	}

}

func RegisterApi(label string, method, url string, f func(c *define.BasicContext)) {
	common.InputApiList = append(common.InputApiList, common.RequestResourceDefine{
		Label:      label,
		Permission: url,
		ParentID:   "api",
	})
	switch method {
	case http.MethodPost:
		global.ApiPrefixAuth.POST(url, define.WrapHandler(f))
	case http.MethodDelete:
		global.ApiPrefixAuth.DELETE(url, define.WrapHandler(f))
	case http.MethodPut:
		global.ApiPrefixAuth.PUT(url, define.WrapHandler(f))
	case http.MethodGet:
		global.ApiPrefixAuth.GET(url, define.WrapHandler(f))
	}
}

func RegisterPublicApi(method, url string, f func(c *define.BasicContext)) {
	switch method {
	case http.MethodPost:
		global.ApiPrefix.POST(url, define.WrapHandler(f))
	case http.MethodDelete:
		global.ApiPrefix.DELETE(url, define.WrapHandler(f))
	case http.MethodPut:
		global.ApiPrefix.PUT(url, define.WrapHandler(f))
	case http.MethodGet:
		global.ApiPrefix.GET(url, define.WrapHandler(f))
	}
}
