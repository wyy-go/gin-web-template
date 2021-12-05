package dao

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSetup(t *testing.T) {
	c := Config{
		DriverName:     "mysql",
		DataSourceName: "root:123456@tcp(127.0.0.1:3306)/xxxx",
		MaxIdleConn:    3,
		MaxOpenConn:    10,
	}

	Convey("should get engine", t, func() {
		engine := Setup(c)
		So(engine, ShouldNotBeNil)
	})
}

func TestDataSource_Open(t *testing.T) {
	c := Config{
		DriverName:     "mysql",
		DataSourceName: "root:123456@tcp(127.0.0.1:3306)/xxxx",
		MaxIdleConn:    3,
		MaxOpenConn:    10,
	}

	dataSource := new(dataSource)
	Convey("service datasource", t, func() {
		Convey("should report error when close database if it's not opened", func() {
			err := dataSource.Close()
			So(err, ShouldEqual, DatabaseIsNotOpenedError)
		})
		Convey("should open database", func() {
			err := dataSource.Open(c)
			So(err, ShouldBeNil)
		})
		Convey("should close database", func() {
			err := dataSource.Close()
			So(err, ShouldBeNil)
		})
	})
}
