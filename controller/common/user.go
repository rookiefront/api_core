package common

import (
	"encoding/json"
	"fmt"
	"github.com/front-ck996/csy"
	"github.com/front-ck996/csy/store"
	"github.com/gin-gonic/gin"
	"github.com/rookiefront/api-core/define"
	"github.com/rookiefront/api-core/global"
	"github.com/rookiefront/api-core/model"
	"github.com/rookiefront/api-core/model/manage_api"
	"github.com/rookiefront/api-core/service"
	"github.com/rookiefront/api-core/utils/common"
	"strings"
)

type ApiUser struct {
}

type userRegister struct {
	NickName  string `json:"nickName" validate:"required"`
	UserName  string `json:"userName" validate:"required"`
	PassWord  string `json:"password" validate:"required"`
	CaptchaId string `json:"captchaId"`
	Captcha   string `json:"captcha"`
}

type userLogin struct {
	UserName  string `json:"userName" validate:"required"`
	PassWord  string `json:"password" validate:"required"`
	CaptchaId string `json:"captchaId" validate:"required"`
	Captcha   string `json:"captcha"`
}
type userUpdate struct {
	model.Model
	UserName string `json:"userName" validate:"required"`
}
type userSetPassword struct {
	UserId   model.PrimarykeyType `json:"userId"`
	Password string               `json:"password"`
}

func (api *ApiUser) SetPassword(c *define.BasicContext) {
	if c.VerifyRequestQualification("btn_set_password") != nil {
		return
	}
	req := userSetPassword{}
	c.ShouldBindJSON(&req)
	if req.UserId == 0 || req.Password == "" {
		c.SendJsonErr("参数为空")
		return
	}

	currentUser := model.SysUser{}
	global.DB.Where("id = ?", req.UserId).First(&currentUser)
	if !currentUser.IdTure() {
		c.SendJsonErr("修改失败")
		return
	}
	if currentUser.Sign == "" {
		currentUser.Sign = csy.RandomString(8)
	}

	global.DB.Where("id = ?", req.UserId).Updates(&model.SysUser{
		Password: service.User.Encrypt(req.Password, currentUser.Sign),
		Sign:     currentUser.Sign,
	})
	c.SendJsonOk()
}

func (api *ApiUser) Register(c *define.BasicContext) {
	var req userRegister
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.SendJsonErr(err)
		return
	}
	// 验证码是否输如正确
	match := common.Captcha.Verify(req.CaptchaId, req.Captcha, true)
	if !match {
		c.SendJsonErr("验证码输入错误")
		return
	}

	err = common.Validate.Struct(req)
	if err != nil {
		c.SendJsonErr(err)
		return
	}
	userSign := csy.RandomString(8)
	role := model.SysRole{}
	role.ID = 2
	userData := model.SysUser{
		NickName: req.NickName,
		UserName: req.UserName,
		Sign:     userSign,
		Password: service.User.Encrypt(req.PassWord, userSign),
		Enable:   1,
		RoleList: []model.SysRole{role},
	}
	err = service.User.VerifyRegister(userData)

	if err != nil {
		c.SendJsonErr(err)
		return
	}
	tx := global.DB.Save(&userData)
	if tx.Error != nil {
		c.SendJsonErr(tx.Error)
		return
	}
	token := service.User.GenerateToken(userData)

	c.SendJsonOk(gin.H{
		"token": token,
	})
}

func (api *ApiUser) Login(c *define.BasicContext) {
	var req userLogin
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.SendJsonErr(err)
		return
	}
	// 验证码是否输如正确
	match := common.Captcha.Verify(req.CaptchaId, req.Captcha, true)
	if !match {
		c.SendJsonErr("验证码输入错误")
		return
	}

	err = common.Validate.Struct(req)
	if err != nil {
		c.SendJsonErr(err)
		return
	}
	where := model.SysUser{
		UserName: req.UserName,
	}
	currentUser := model.SysUser{}
	global.DB.Model(model.SysUser{}).Where(where).First(&currentUser)
	if currentUser.Model.ID == 0 {
		c.SendJsonErr("用户不存在")
		return
	}
	password := service.User.Encrypt(req.PassWord, currentUser.Sign)

	if currentUser.Password != password {
		c.SendJsonErr("密码错误")
		return
	}
	token := service.User.GenerateToken(currentUser)

	c.SendJsonOk(gin.H{
		"token": token,
	})
}

func (api *ApiUser) UserInfo(c *define.BasicContext) {
	id := c.GetCurrentUserId()
	var result model.SysUser
	global.DB.Model(model.SysUser{}).Omit("Password").Preload("RoleList").Where("id = ?", id).First(&result)
	c.SendJsonOk(result)
}

func (api *ApiUser) CurrentUserMenu(c *define.BasicContext) {
	id := c.GetCurrentUserId()
	var result model.SysUser
	global.DB.Model(model.SysUser{}).Preload("RoleList").Where("id = ?", id).First(&result)
	var menuList []model.SysMenu
	if result.UserName == "user_root" {
		global.DB.Find(&menuList)
		c.SendJsonOk(menuList)
		return
	}

	hash := map[model.PrimarykeyType]struct{}{}
	for i, _ := range result.RoleList {
		global.DB.Preload("MenuList").First(&result.RoleList[i])
		for _, menu := range result.RoleList[i].MenuList {
			if _, ok := hash[menu.Model.ID]; ok {
				continue
			}
			hash[menu.Model.ID] = struct{}{}
			menuList = append(menuList, menu)
		}
	}
	c.SendJsonOk(menuList)
}

func (api *ApiUser) Logout(c *define.BasicContext) {
	id := c.GetCurrentUserId()
	bucket := service.User.Store.SetBucket(id)
	err := store.Remove(bucket, c.GetToken())
	if err != nil {
		c.SendJsonErr(err)
		return
	}
	c.SendJsonOk()
}

func (api *ApiUser) Update(c *define.BasicContext) {
	var req userUpdate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.SendJsonErr(err)
		return
	}
}

func (api *ApiUser) RoleMenuList(c *define.BasicContext) {
	query, _ := c.GetQuery("id")
	if query == "" {
		c.SendJsonErr("")
		return
	}
	var result []map[string]interface{}
	global.DB.Table("m2m_sys_role_sys_menu").Where("sys_role_id = ?", query).Find(&result)
	var list []interface{}
	for _, m := range result {
		list = append(list, m["sys_menu_id"])
	}
	c.SendJsonOk(list)
}

type RequestResourceDefine struct {
	Label      string `json:"label"`
	Permission string `json:"permission"`
	ParentID   string `json:"parent_id"`
}

type RequestResourceDefineUpdate struct {
	RequestResourceDefine
	ID    model.PrimarykeyType `json:"id"`
	Value []string             `json:"value"`
}

func (api *ApiUser) UpdateRequestResource(c *define.BasicContext) {

	if err := c.VerifyRequestQualification("btn_allow_request_resource"); err != nil {
		c.SendJsonErr(err)
		return
	}

	var req RequestResourceDefineUpdate
	c.ShouldBindJSON(&req)
	saveData := model.SysRole{
		PermissionList: strings.Join(req.Value, ","),
	}
	saveData.ID = req.ID
	if saveData.ID == 0 {
		return
	}
	if saveData.PermissionList == "" {
		saveData.PermissionList = " "
	}
	tx := global.DB.Updates(&saveData)
	if tx.Error != nil {
		c.SendJsonErr(tx.Error)
		return
	}
	c.SendJsonOk("")
}

func (api *ApiUser) RequestResource(c *define.BasicContext) {
	if c.VerifyRequestQualification("btn_allow_request_resource") != nil {
		return
	}
	var moduleList []manage_api.ManageApiModule
	global.DB.Preload("Fields").Where("table_name = ?", "sys_user").Find(&moduleList)
	global.DB.Preload("Fields").Find(&moduleList)
	var result []RequestResourceDefine
	for _, module := range moduleList {
		result = append(result, RequestResourceDefine{
			Label:      "模块_" + module.Name,
			Permission: module.TaName,
			ParentID:   "",
		})
		result = append(result, RequestResourceDefine{
			Label:      "查询",
			Permission: fmt.Sprintf("%s_query", module.TaName),
			ParentID:   module.TaName,
		})

		result = append(result, RequestResourceDefine{
			Label:      "新增",
			Permission: fmt.Sprintf("%s_add", module.TaName),
			ParentID:   module.TaName,
		})

		result = append(result, RequestResourceDefine{
			Label:      "修改",
			Permission: fmt.Sprintf("%s_edit", module.TaName),
			ParentID:   module.TaName,
		})

		result = append(result, RequestResourceDefine{
			Label:      "删除",
			Permission: fmt.Sprintf("%s_delete", module.TaName),
			ParentID:   module.TaName,
		})

		for _, btn := range module.TableOperateBtn {
			var btn2 manage_api.ManageApiModuleFieldTableOperateBtn
			marshal, _ := json.Marshal(btn)
			json.Unmarshal(marshal, &btn2)
			result = append(result, RequestResourceDefine{
				Label:      btn2.Name + "_" + "[按钮]",
				Permission: fmt.Sprintf("btn_%s", btn2.Permission),
				ParentID:   module.TaName,
			})
		}

		for _, a := range module.Fields {
			if a.Associations.Type == "ManyToMany" {
				result = append(result, RequestResourceDefine{
					Label:      fmt.Sprintf("查询[%s][ManyToMany]", a.FrontField),
					Permission: fmt.Sprintf("%s_query_%s", module.TaName, a.FrontField),
					ParentID:   module.TaName,
				})
				result = append(result, RequestResourceDefine{
					Label:      fmt.Sprintf("修改[%s][ManyToMany]", a.FrontField),
					Permission: fmt.Sprintf("%s_edit_%s", module.TaName, a.FrontField),
					ParentID:   module.TaName,
				})
			}
		}
	}

	for _, module := range moduleList {
		for _, a := range module.Fields {
			if a.Associations.Type == "HasMany" {
				currentModule := manage_api.GetModule(a.Associations.Module)
				result = append(result, RequestResourceDefine{
					Label:      "子模块_" + currentModule.Name + "[HasMany]",
					Permission: currentModule.TaName,
					ParentID:   module.TaName,
				})
			}
		}
	}
	c.SendJsonOk(result)
}
