package global

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var Engine *gin.Engine
var DB *gorm.DB
