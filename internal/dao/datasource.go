package dao

import (
	"errors"
	"github.com/wyy-go/go-web-template/pkg/env"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

type Config struct {
	DriverName     string
	DataSourceName string
	MaxIdleConn    int
	MaxOpenConn    int
}

type DataSource interface {
	Open(c Config) error
	IsOpened() bool
	Close() error
	Engine() *gorm.DB
}

type dataSource struct {
	engine *gorm.DB
}

var (
	DatabaseIsNotOpenedError = errors.New("database is not opened")
	ds                       *dataSource
	dsOnce                   sync.Once
)

func GetDataSource() DataSource {
	dsOnce.Do(func() {
		ds = new(dataSource)
	})
	return ds
}

func (d *dataSource) Open(c Config) error {
	gormConfig := gorm.Config{}

	if env.IsDeployDebug() {
		gormConfig.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // 慢 SQL 阈值
				LogLevel:      logger.Info, // Log level
				Colorful:      false,       // 禁用彩色打印
			},
		)
	}

	db, err := gorm.Open(mysql.Open(c.DataSourceName), &gormConfig)
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	if c.MaxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	}

	if c.MaxOpenConn > 0 {
		sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	}

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	d.engine = db
	return nil
}

func (d *dataSource) IsOpened() bool {
	return d.engine != nil
}

func (d *dataSource) Close() error {
	if d.engine != nil {
		sqlDB, err := d.engine.DB()
		if err != nil {
			return err
		}
		sqlDB.Close()
		d.engine = nil
		return err
	}
	return DatabaseIsNotOpenedError
}

func (d *dataSource) Engine() *gorm.DB {
	return d.engine
}

func Setup(c Config) *gorm.DB {
	dataSource := GetDataSource()
	if !dataSource.IsOpened() {
		if err := dataSource.Open(c); err != nil {
			panic(err)
		}
	}
	return dataSource.Engine()
}

func GetDB() *gorm.DB {
	return GetDataSource().Engine()
}
