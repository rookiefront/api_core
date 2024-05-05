package router

import (
	"github.com/rookiefront/api-core/controller/common"
	"github.com/rookiefront/api-core/controller/easy_curd"
	"github.com/rookiefront/api-core/define"
	"github.com/rookiefront/api-core/global"
	"github.com/rookiefront/api-core/middleware"
)

func registerSystemRouter() {
	commonApi := global.Engine.Group("/api/common")
	commonApiAuti := global.Engine.Group("/api/common").Use(middleware.InterceptNotLoggedIn)
	apiUser := common.ApiUser{}

	commonApi.GET("/captcha", define.WrapHandler(common.GetCaptcha))
	//// 用户相关
	// 注册
	commonApi.POST("/user/register", define.WrapHandler(apiUser.Register))
	//
	// 登录
	commonApi.POST("/user/login", define.WrapHandler(apiUser.Login))
	//
	// 当前登录用户的详情
	commonApiAuti.GET("/user/info", define.WrapHandler(apiUser.UserInfo))
	//
	// 退出登录
	commonApiAuti.GET("/user/logout", define.WrapHandler(apiUser.Logout))
	commonApiAuti.POST("/user/update", define.WrapHandler(apiUser.Update))
	//// 获取角色的菜单列表
	commonApiAuti.GET("/user/role_menu_list", define.WrapHandler(apiUser.RoleMenuList))
	commonApiAuti.GET("/user/current_user_menu", define.WrapHandler(apiUser.CurrentUserMenu))
	commonApiAuti.POST("/user/set_password", define.WrapHandler(apiUser.SetPassword))

	//// 可以访问的接口资源
	commonApiAuti.GET("/user/request_resource", define.WrapHandler(apiUser.RequestResource))
	//
	commonApiAuti.PUT("/user/request_resource", define.WrapHandler(apiUser.UpdateRequestResource))
	// 获取配置详情
	commonApiAuti.GET("/manage/config", define.WrapHandler(easy_curd.ManageConfig))

	///// 字典
	//apiDict := common.ApiDict{}
	//commonApi.GET("/dict/list", define.WrapHandler(apiDict.List))
	/////	城市选择
	//apiCity := common.ApiCity{}
	//commonApi.GET("/city/lazy", define.WrapHandler(apiCity.Lazy))
	//
	/// 步骤选择
	apiStep := common.ApiStep{}
	commonApi.GET("/step/refresh_full_path", define.WrapHandler(apiStep.RefreshFulPath))

	/// 文件上传
	apiUpload := common.ApiUpload{}
	commonApi.POST("/upload/image", define.WrapHandler(apiUpload.UploadImage))
	commonApi.POST("/upload/file", define.WrapHandler(apiUpload.UploadFile))
}
