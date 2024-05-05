package model_test

import (
	"github.com/rookiefront/api-core/cmd/model"
	"github.com/rookiefront/api-core/config"
	"github.com/rookiefront/api-core/global"
	"github.com/rookiefront/api-core/service"
	"testing"
)

func TestSetAdmin(t *testing.T) {

	config.LoadConfig()
	config.DbConnect()

	defaultUsers := []map[string]string{
		{
			"user_name": "admin",
			"password":  "admin@123",
		},
		{
			"user_name": "test_user_01",
			"password":  "test_user_01",
		},
	}
	for _, user := range defaultUsers {
		u_Admin := model.SysUser{}
		global.DB.Table("sys_user").Where("user_name = ?", user["user_name"]).First(&u_Admin)
		u_Admin.Password = service.User.Encrypt(user["password"], "")
		u_Admin.NickName = user["user_name"]
		u_Admin.UserName = user["user_name"]

		global.DB.Save(&u_Admin)
	}

}
