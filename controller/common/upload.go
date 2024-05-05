package common

import (
	"fmt"
	"github.com/front-ck996/csy"
	"github.com/gin-gonic/gin"
	config2 "github.com/rookiefront/api-core/config"
	"github.com/rookiefront/api-core/define"
	"time"
)

type ApiUpload struct {
}

func (api *ApiUpload) UploadImage(c *define.BasicContext) {
	config := config2.GetConfig()
	file, err := c.FormFile("file")
	if err != nil {
		c.SendJsonErr(err)
		return
	}
	fileName := fmt.Sprintf("/%s.png", time.Now().Format("2006_01_02_15_04_05")+csy.RandomString(5))
	saveFilePath := config.System.FullUploadDir + fileName
	filePath := config.System.SiteUploadDir + fileName
	err = c.SaveUploadedFile(file, saveFilePath)
	if err != nil {
		c.SendJsonErr(err)
		return
	}
	siteUrl := "http://"
	if c.Request.TLS != nil {
		siteUrl = "https://"
	}
	siteUrl += c.Request.Host
	c.SendJsonOk(gin.H{
		"file": filePath,
		"url":  siteUrl + filePath,
	})
}

func (api *ApiUpload) UploadFile(c *define.BasicContext) {

}
