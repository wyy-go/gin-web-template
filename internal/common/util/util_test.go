package util

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/wyy-go/go-web-template/pkg/idgen"
	"testing"
)

func TestIdCodec(t *testing.T) {
	Convey("should equal IdEncode IdDecode", t, func() {
		id := idgen.Next()
		code := IdEncode(id)
		So(IdDecode(code), ShouldEqual, id)
	})
}
