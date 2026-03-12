package app

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"quickBillController/config"
)

var DB *gorm.DB

const TablePrefix = "sbs_"

func InitDatabase() {
	dbCf := config.GetCfg().Db
	cfg := fmt.Sprintf("host=%s port=%v user=%s dbname=%s sslmode=%s password=%s",
		dbCf.Host,
		dbCf.Port,
		dbCf.User,
		dbCf.Name,
		dbCf.Ssl,
		dbCf.Password)
	db, err := gorm.Open(postgres.Open(cfg), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   TablePrefix,
			SingularTable: false,
		},
		Logger:      dbLog(),
		PrepareStmt: true,
	})
	if err != nil {
		panic(fmt.Sprintf("Database connection failed %v", err.Error()))
	}
	db = db.Debug()
	if config.GetCfg().Db.Schema != "" {
		db.Exec(fmt.Sprintf("SET SEARCH_PATH TO '%s',\"$user\",'public'", config.GetCfg().Db.Schema))
		db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	}

	DB = db

}

func dbLog() logger.Interface {
	logLevel := logger.Info
	switch config.GetCfg().Db.LogLevel {
	case "debug":
		logLevel = logger.Info
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}
	newLogger := logger.New(
		Writer{},
		logger.Config{
			SlowThreshold:             time.Duration(2) * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logLevel,                            // Log level
			IgnoreRecordNotFoundError: true,                                // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                                // Disable color
		},
	)
	newLogger.LogMode(logLevel)
	return newLogger
}

func (w Writer) Printf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
