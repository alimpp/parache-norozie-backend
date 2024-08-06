package api

import (
	"ecom/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitSqlDb(conf config.ConfStruct) *gorm.DB {
	if conf.DB.DriverName == "sqlite" {
		db, err := gorm.Open(sqlite.Open(conf.DB.DataSourceName), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		return db
	} else if conf.DB.DriverName == "postgres" {
		db, err := gorm.Open(postgres.Open(conf.DB.DataSourceName), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		return db
	} else {
		panic("unsupported database driver")
	}
}
