/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package operations_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"os/exec"
	"gir/gio-2.0"
	. "pkg.deepin.io/service/file-manager-backend/operations"
	"testing"
)

func _TestChownJob(t *testing.T) {
	Convey("test chown with a file", t, func() {
		exec.Command("cp", "./testdata/chown/testfile", "./testdata/chown/afile").Run()
		defer exec.Command("rm", "./testdata/chown/afile").Run()
		u, err := pathToURL("./testdata/chown/afile")
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		// TODO: generate a not existed group.
		job := NewChownJob(u.String(), "xx", "xx")
		job.Execute()
		So(job.HasError(), ShouldBeTrue)

		// TODO: make sure the original group is not targetGroup.
		// and permission denied won't happen.
		targetGroup := "video"
		job2 := NewChownJob(u.String(), "", targetGroup)
		job2.Execute()
		So(job2.HasError(), ShouldBeFalse)

		f := gio.FileNewForPath("./testdata/chown/afile")
		info, _ := f.QueryInfo(gio.FileAttributeOwnerGroup, gio.FileQueryInfoFlagsNofollowSymlinks, nil)
		So(info.GetAttributeString(gio.FileAttributeOwnerGroup), ShouldEqual, targetGroup)
		info.Unref()

	})
}
