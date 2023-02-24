package bootstrap

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/hildxd/backend-template/app/models"
	"github.com/hildxd/backend-template/global"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitializeDB() *gorm.DB {
	var db *gorm.DB
	// 根据驱动配置进行初始化
	switch global.App.Config.Database.Driver {
	case "postgres":
		db = initPostgresqlGorm()
	default:
		db = initPostgresqlGorm()
	}
	global.App.DB = db
	return db
}

// 数据库表初始化
func initPostgresqlTables(db *gorm.DB) {
	err := db.AutoMigrate(
		models.User{},
	)
	if err != nil {
		global.App.Log.Error("migrate table failed", zap.Any("err", err))
		os.Exit(0)
	}
}

func initPostgresqlGorm() *gorm.DB {
	dbConfig := global.App.Config.Database

	if dbConfig.Database == "" {
		return nil
	}

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", dbConfig.UserName, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)

	pgConfig := postgres.Config{
		DSN: dsn,
	}
	if db, err := gorm.Open(postgres.New(pgConfig), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
		Logger:                                   getGormLogger(),
	}); err != nil {
		global.App.Log.Error("初始化数据库连接失败", zap.Any("err", err))
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
		sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
		initPostgresqlTables(db)
		return db
	}
}

func getGormLogWriter() logger.Writer {
	var writer io.Writer

	// 是否启用日志文件
	if global.App.Config.Database.EnableFileLogWriter {
		// 自定义writer
		writer = &lumberjack.Logger{
			Filename:   global.App.Config.Log.RootDir + "/" + global.App.Config.Database.LogFilename,
			MaxSize:    global.App.Config.Log.MaxSize,
			MaxBackups: global.App.Config.Log.MaxBackups,
			MaxAge:     global.App.Config.Log.MaxAge,
			Compress:   global.App.Config.Log.Compress,
		}
	} else {
		writer = os.Stdout
	}
	return log.New(writer, "\r\n", log.LstdFlags)
}

func getGormLogger() logger.Interface {
	var logMode logger.LogLevel

	switch global.App.Config.Database.LogMode {
	case "silent":
		logMode = logger.Silent
	case "error":
		logMode = logger.Error
	case "warn":
		logMode = logger.Warn
	case "info":
		logMode = logger.Info
	default:
		logMode = logger.Info
	}

	return logger.New(getGormLogWriter(), logger.Config{
		SlowThreshold:             time.Millisecond * 200, // 慢 SQL 阈值
		LogLevel:                  logMode,
		IgnoreRecordNotFoundError: false,                                           // 忽略 ErrRecordNotFound 错误
		Colorful:                  !global.App.Config.Database.EnableFileLogWriter, // 禁用彩色打印
	})
}
