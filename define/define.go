package define

import (
	"errors"
	"fmt"
	"github.com/front-ck996/csy"
	"github.com/gin-gonic/gin"
	"github.com/rookiefront/api-core/config"
	"github.com/rookiefront/api-core/global"
	"github.com/rookiefront/api-core/model"
	"github.com/rookiefront/api-core/service"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type BasicContext struct {
	*gin.Context
	Where any
}

func (c *BasicContext) SendJsonOk(data ...interface{}) {
	var message interface{}
	if len(data) >= 1 {
		message = data[0]
	}
	outerJson := gin.H{
		"msg":  "ok",
		"code": 200,
		"data": message,
	}
	if config.IsDev() {
		outerJson["where"] = c.Where
	}
	c.JSON(200, outerJson)
}
func (c *BasicContext) SendJsonOkPage(data interface{}, pageData gin.H) {
	outerJson := gin.H{
		"msg":  "ok",
		"code": 200,
		"data": data,
	}
	for k, v := range pageData {
		outerJson[k] = v
	}
	if config.IsDev() {
		outerJson["where"] = c.Where
	}
	c.JSON(200, outerJson)
}
func (c *BasicContext) SendJsonErr(err any) {
	if csy.IsError(err) && err != nil {
		err = err.(error).Error()
	}
	c.JSON(200, gin.H{
		"msg":  err,
		"code": 500,
		"data": nil,
	})
}
func (c *BasicContext) SendJsonErrCode(err any, code any) {
	if csy.IsError(err) && err != nil {
		err = err.(error).Error()
	}
	c.JSON(200, gin.H{
		"msg":  err,
		"code": code,
		"data": nil,
	})
}

func WrapHandler(handler func(c *BasicContext)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		myCtx := &BasicContext{Context: ctx}
		handler(myCtx)
	}
}

func (c *BasicContext) GetPostFormParams() (map[string]any, error) {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		if !errors.Is(err, http.ErrNotMultipart) {
			return nil, err
		}
	}
	var postMap = make(map[string]any, len(c.Request.PostForm))
	for k, v := range c.Request.PostForm {
		if len(v) > 1 {
			postMap[k] = v
		} else if len(v) == 1 {
			postMap[k] = v[0]
		}
	}

	return postMap, nil
}

func (c *BasicContext) GetQueryParams() map[string]any {
	query := c.Request.URL.Query()
	var queryMap = make(map[string]any, len(query))
	for k := range query {
		queryMap[k] = c.Query(k)
	}
	return queryMap
}

// 获得请求参数 GET POST FormData JSON 合并
func (c *BasicContext) GetReqData() (reqData map[string]any) {
	query := c.GetQueryParams()
	postQuery, err := c.GetPostFormParams()
	if err == nil {
		for m, v := range postQuery {
			query[m] = v
		}
		if len(postQuery) != 0 {
			return query
		}
	}
	var jsonData map[string]any
	c.ShouldBindJSON(&jsonData)
	for m, v := range jsonData {
		query[m] = v
	}
	return query
}
func (c *BasicContext) GetToken() string {
	return c.GetHeader("X-Token")
}
func (c *BasicContext) GetCurrentUserId() string {
	token := c.GetHeader("X-Token")
	userMap, err := service.User.ParseToken(token)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.0f", userMap["id"])
}

// GetCurrentUserIdToInt 获取用户ID 转换成 Int 类型
func (c *BasicContext) GetCurrentUserIdToInt() int {
	token := c.GetHeader("X-Token")
	userMap, err := service.User.ParseToken(token)
	if err != nil {
		return 0
	}
	int64Num, _ := strconv.ParseInt(fmt.Sprintf("%.0f", userMap["id"]), 10, 64)
	return int(int64Num)
}

// VerifyRequestQualification 验证用户权限是否存在
func (c *BasicContext) VerifyRequestQualification(verifyPermission string) error {
	userId := c.GetCurrentUserId()
	userInfo := model.SysUser{}
	global.DB.Where("id = ?", userId).Preload("RoleList").First(&userInfo)
	if userInfo.UserName == "user_root" {
		return nil
	}
	for _, role := range userInfo.RoleList {
		permissionList := strings.Split(role.PermissionList, ",")
		if csy.SliceInclude[string](permissionList, verifyPermission) {
			return nil
		}
	}
	c.SendJsonErrCode("无权限访问", 40002)
	return errors.New("无权限")
}

func (c *BasicContext) VerifyDbExecuteIsOk(tx *gorm.DB) bool {
	if tx.Error != nil {
		c.SendJsonErrCode(tx.Error, 40003)
		return false
	}
	return true
}
