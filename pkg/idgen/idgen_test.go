package idgen

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/wyy-go/go-web-template/pkg/idgen/snowflake"
	"testing"
	"time"
)

func TestNext(t *testing.T) {
	Convey("generate id", t, func() {
		Convey("should generate id", func() {
			id := Next()
			So(id, ShouldNotEqual, 0)
		})

		Convey("should return zero when over the time limit", func() {
			// setup
			st := snowflake.Settings{
				StartTime: time.Date(1883, 1, 1, 0, 0, 0, 0, time.UTC),
				MachineID: getMachineId,
			}
			sf = snowflake.NewSnowflake(st)

			id := Next()
			So(id, ShouldEqual, 0)

			// teardown
			st = snowflake.Settings{
				MachineID: getMachineId,
			}
			sf = snowflake.NewSnowflake(st)
		})
	})
}

func TestGetOne(t *testing.T) {
	Convey("should generate one ID", t, func() {
		id := GetOne()
		So(id, ShouldNotEqual, 0)
	})
}

func TestGetMulti(t *testing.T) {
	Convey("should generate multiple IDs", t, func() {
		ids := GetMulti(3)
		So(len(ids), ShouldEqual, 3)
		for _, v := range ids[:] {
			So(v, ShouldBeGreaterThan, 0)
		}
	})
}
