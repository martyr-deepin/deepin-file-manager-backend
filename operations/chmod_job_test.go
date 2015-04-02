package operations_test

import (
	. "deepin-file-manager/operations"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"os/exec"
	"testing"
)

func TestChmodJob(t *testing.T) {
	Convey("Chmod a file to 0000", t, func() {
		source := "./testdata/chmod/test.file"
		target := "./testdata/chmod/afile"
		exec.Command("cp", source, target).Run()
		defer exec.Command("rm", target).Run()

		uri, err := pathToURL(target)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		NewChmodJob(uri, 0000).Execute()
		fi, err := os.Stat(target)
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		So(fi.Mode(), ShouldEqual, os.FileMode(0))
	})

	Convey("Chmod a dir to 0000", t, func() {
		source := "./testdata/chmod/test.dir"
		target := "./testdata/chmod/adir"
		exec.Command("cp", "-r", source, target).Run()
		defer exec.Command("rm", "-r", target).Run()

		uri, err := pathToURL(target)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		NewChmodJob(uri, 0000).Execute()
		fi, err := os.Stat(target)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		So(fi.Mode(), ShouldEqual, os.FileMode(os.ModeDir))
	})
}
