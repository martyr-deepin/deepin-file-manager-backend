package operations_test

import (
	. "github.com/smartystreets/goconvey/convey"
	. "pkg.deepin.io/service/file-manager-backend/operations"
	"testing"
)

func TestGetTemplateJob(t *testing.T) {
	Convey("Get template from directory which consists of directories", t, func() {
		uri := "./testdata"
		op := NewGetTemplateJob(uri)
		So(len(op.Execute()), ShouldEqual, 0)
	})
}
