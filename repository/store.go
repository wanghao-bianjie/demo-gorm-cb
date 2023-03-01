package repository

import (
	"demo-gorm-cb/repository/callback"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	dbLogger "gorm.io/gorm/logger"
)

var (
	db *gorm.DB
)

func InitMysqlDB(dsn string) {
	//db logger
	// logger
	newLogger := dbLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		dbLogger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      dbLogger.Warn, // Log level
			Colorful:      false,         // Disable color
		},
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		logrus.Fatalf(dsn)
	}
}

func RegisterCallbacks() {
	var errs = []error{
		db.Callback().Create().Before("gorm:create").Register("my_plugin:before_create", callback.BeforeCreate),
		db.Callback().Create().After("gorm:create").Register("my_plugin:after_create", callback.AfterCreate),
		db.Callback().Update().Before("gorm:update").Register("my_plugin:before_update", callback.BeforeUpdate),
		db.Callback().Update().After("gorm:update").Register("my_plugin:after_update", callback.AfterUpdate),
		db.Callback().Query().After("gorm:query").Register("my_plugin:after_query", callback.AfterQuery),
	}
	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}

func GetDb() *gorm.DB {
	return db
}
