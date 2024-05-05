package config_test

import (
	"github.com/rookiefront/api_core/config"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func TestSetConfig(t *testing.T) {
	write_config := config.Config{
		System: config.System{
			Mode:                config.Development,
			Host:                "127.0.0.1",
			Port:                9099,
			RootDir:             "./",
			DbDir:               "/dbs",
			StaticDir:           "/static/deploy/",
			UploadDir:           "/static/deploy/uploads",
			ApiPrefix:           "/api/v1",
			BusinessTablePreFix: "bus_",
		},
		Db: config.DB{
			DSN:         "root:root@tcp(192.168.0.190:13306)/2023_12_19_site_api?charset=utf8mb4&parseTime=True&loc=Local",
			MaxOpenConn: 100,
			MaxIdleConn: 10,
		},
	}
	out, _ := yaml.Marshal(write_config)
	os.WriteFile("../config.yaml", out, 0644)
}
