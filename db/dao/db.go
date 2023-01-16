package dao

import (
	"fmt"
	"github.com/ChineseSubFinder/csf-supplier-base/db/models"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/settings"
	"github.com/WQGroup/logger"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	dbLogger "gorm.io/gorm/logger"
	"os"
	"path/filepath"
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

	logger.Infoln("DbType:", settings.Get().DBConnectConfig.DbType)

	if settings.Get().DBConnectConfig.DbType == "sqlite3" {
		// sqlite3
		nowDBFName := filepath.Join(".", dbFileName)
		dbDir := filepath.Dir(nowDBFName)
		if pkg.IsDir(dbDir) == false {
			err = os.MkdirAll(dbDir, os.ModePerm)
			if err != nil {
				return err
			}
		}
		db, err = gorm.Open(sqlite.Open(nowDBFName), &gorm.Config{})
		if err != nil {
			return errors.New(fmt.Sprintf("failed to connect database, %s", err.Error()))
		}
	} else if settings.Get().DBConnectConfig.DbType == "mysql" {
		// mysql
		db, err = gorm.Open(mysql.Open(settings.Get().DBConnectConfig.DataSource), &gorm.Config{})
		if err != nil {
			return errors.New(fmt.Sprintf("failed to connect database, %s", err.Error()))
		}
	} else {
		return errors.New("not support db type")
	}

	// 降低 gorm 的日志级别
	db.Logger = dbLogger.Default.LogMode(dbLogger.Silent)
	// 迁移 schema
	err = db.AutoMigrate(
		&models.Huo720Media{},
		&models.Movie{},
		&models.Tv{},

		&models.HotMovie{},
		&models.HotTV{},

		&models.DownloadedInfoZiMuKu{},
		//&models.ZiMuKuMovie{},
		//&models.ZiMuKuTV{},

		&models.HouseKeeping{},
		&models.HouseKeepingError{},
		&models.SubtitleMovie{},
		&models.SubtitleTV{},

		&models.ManualHandling{},

		&models.TvDetailInfo{},
		&models.SeasonDetailInfo{},
		&models.EpisodeDetailInfo{},
		&models.MovieDetailInfo{},

		&models.TopMovie{},
		&models.TopTv{},
	)
	if err != nil {
		return errors.New(fmt.Sprintf("db AutoMigrate error, %s", err.Error()))
	}

	return nil
}

const (
	dbFileName = "csf-supplier.db"
)

var (
	db   *gorm.DB
	once sync.Once
)
