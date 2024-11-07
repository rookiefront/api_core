package config

import (
	"context"
	"fmt"
	"github.com/rookiefront/api-core/global"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)
import "gorm.io/driver/mysql"

type dbLogger struct {
	logger *log.Logger
}

func (s dbLogger) LogMode(level logger.LogLevel) logger.Interface {
	return s
}

func (s dbLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	fmt.Println(msg)
	// 可以选择是否记录 Info 级别的日志
}

func (s dbLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	// 记录警告日志
	s.logger.Printf(msg, data...)
}

func (s dbLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	// 记录错误日志
	s.logger.Printf(msg, data...)
}

func (s dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if elapsed > time.Second { // 将慢查询的阈值设定为 1 秒
		s.logger.Printf("SLOW QUERY: %s [%v] (%d rows)", sql, elapsed, rows)
	}
}

func DbConnect() {
	conf := GetConfig()
	slowLogger := log.New(os.Stdout, "SLOW QUERY: ", log.LstdFlags)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: conf.Db.DSN,
		//"gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
		//SkipDefaultTransaction:                   true,
		Logger: dbLogger{logger: slowLogger},
	})
	if err != nil {
		log.Fatalln(err)
		return
	}

	db.Logger = logger.Default.LogMode(logger.Error)
	global.DB = db
	s, err := db.DB()
	if err != nil {
		log.Fatalln(err)
		return
	}
	// 最大空闲链接数
	s.SetMaxIdleConns(conf.Db.MaxIdleConn)
	// 最大链接数
	s.SetMaxOpenConns(conf.Db.MaxOpenConn)
	// 链接最大空闲时间
	s.SetConnMaxLifetime(time.Hour)
}
