package common

import (
	"github.com/gin-gonic/gin"
	"github.com/rookiefront/api-core/config"
	"github.com/rookiefront/api-core/define"
	"github.com/xiaoqidun/qqwry"
)

var qqwryLoaded = false

func GeoIp(c *define.BasicContext) {
	currentConfig := config.GetConfig()

	// 从文件加载IP数据库
	if qqwryLoaded == false {
		if err := qqwry.LoadFile(currentConfig.System.RootDir + "/qqwry.dat"); err != nil {
			qqwryLoaded = true
			c.SendJsonErr(err)
			return
		}
	}
	clientIP := c.ClientIP()
	if clientIP == "127.0.0.1" {
		clientIP = "14.111.255.255"
	}
	// 从内存或缓存查询IP
	city, isp, err := qqwry.QueryIP(clientIP)
	if err != nil {
		c.SendJsonErr(err)
	}
	c.SendJsonOk(gin.H{
		"city": city,
		"isp":  isp,
		"ip":   clientIP,
	})
}
