package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

const Version = "1.0.0"

type Config struct {
	System System `yaml:"system"`
	Db     DB     `yaml:"db"`
}

const (
	Development = "development"
	Production  = "production"
)

type System struct {
	Mode                string `yaml:"mode"`
	Host                string `yaml:"host"`
	Port                int    `yaml:"port"`
	RootDir             string `yaml:"root_dir"`
	DbDir               string `yaml:"db_dir"`
	StaticDir           string `yaml:"static_dir"`
	UploadDir           string `yaml:"upload_dir"`
	SiteUploadDir       string `yaml:"-"`
	FullUploadDir       string `yaml:"-"`
	FullDbDir           string `yaml:"-"`
	ApiPrefix           string `yaml:"api_prefix"`             // api 地址前缀
	BusinessTablePreFix string `yaml:"business_table_pre_fix"` // 业务表前缀
}

type DB struct {
	DSN         string `yaml:"dsn"`
	MaxIdleConn int    `yaml:"maxIdleConn"`
	MaxOpenConn int    `yaml:"maxOpenConn"`
}

var _config = Config{}

func LoadConfig() {
	configFiles := []string{
		"config.yaml",
		"../../config.yaml",
	}
	file := []byte{}
	var err error
	for _, f := range configFiles {
		file, err = os.ReadFile(f)
		if err == nil {
			break
		}
	}
	err = yaml.Unmarshal(file, &_config)
	if err != nil {
		panic(err)
	}
	_config.System.DbDir = _config.System.RootDir + _config.System.DbDir
	_config.System.FullDbDir, _ = filepath.Abs(_config.System.DbDir)
	_config.System.SiteUploadDir = _config.System.UploadDir
	_config.System.UploadDir = _config.System.RootDir + "/public/" + _config.System.UploadDir
	_config.System.FullUploadDir, _ = filepath.Abs(_config.System.UploadDir)
	_config.System.StaticDir = _config.System.RootDir + _config.System.StaticDir
}

func GetConfig() Config {
	return _config
}

func IsDev() bool {
	return _config.System.Mode == Development
}
