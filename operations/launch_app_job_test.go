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
	. "pkg.deepin.io/service/file-manager-backend/operations"
	"testing"
)

func TestGetLaunchAppInfo(t *testing.T) {
	// FIXME: how to make a stable test???
	// Convey("get launch app info", t, func() {
	// 	uri, _ := pathToURL("./testdata/launchapp/test.c")
	// 	job := NewLaunchAppJob(uri)
	// 	appInfo := job.Execute()
	// 	So(job.HasError(), ShouldBeFalse)
	// 	t.Log(appInfo)
	// })
}

func TestSetLaunchAppInfo(t *testing.T) {
	// FIXME: how to make a stable test???
	SkipConvey("set launch app info", t, func() {
		mimeType := "text/html"
		job := NewSetDefaultLaunchAppJob("google-chrome.desktop", mimeType)
		job.Execute()
		So(job.HasError(), ShouldBeFalse)
	})
}
