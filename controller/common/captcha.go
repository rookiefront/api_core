package common

import (
	"github.com/gin-gonic/gin"
	"github.com/rookiefront/api-core/define"
	"github.com/rookiefront/api-core/utils/common"
)

func GetCaptcha(c *define.BasicContext) {
	id, b64s, err := common.Captcha.Generate()
	if err != nil {
		c.SendJsonErr(err)
		return
	}
	c.SendJsonOk(gin.H{
		"captchaId": id,
		"img":       b64s,
	})
}
