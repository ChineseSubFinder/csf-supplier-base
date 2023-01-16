package dao

import (
	"fmt"
	models2 "github.com/ChineseSubFinder/csf-supplier-base/pkg/imdb_info_center/models"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/settings"
	"github.com/WQGroup/logger"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	dbLogger "gorm.io/gorm/logger"
	"sync"
)

// Get 获取数据库实例
func Get() *gorm.DB {
	if db == nil {
		once.Do(func() {
			err := initDB()
			if err != nil {
				logger.Panicln(err)
			}
		})
	}
	return db
}

func initDB() error {

	var err error
	// mysql
	db, err = gorm.Open(mysql.Open(settings.Get().ImdbInfoCenterConfig.DBConnectConfig.DataSource), &gorm.Config{})
	if err != nil {
		return errors.New(fmt.Sprintf("failed to connect database, %s", err.Error()))
	}
	// 降低 gorm 的日志级别
	db.Logger = dbLogger.Default.LogMode(dbLogger.Silent)
	// 迁移 schema
	err = db.AutoMigrate(
		&models2.TitleBasic{},
		&models2.TitleEpisode{},
		&models2.TitleRatings{},
	)
	if err != nil {
		return errors.New(fmt.Sprintf("db AutoMigrate error, %s", err.Error()))
	}

	return nil
}

var (
	db   *gorm.DB
	once sync.Once
)
