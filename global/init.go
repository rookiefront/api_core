package global

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var Engine *gin.Engine
var ApiManageRouter *gin.RouterGroup
var ApiPrefix *gin.RouterGroup
var ApiPrefixAuth gin.IRoutes
var DB *gorm.DB
