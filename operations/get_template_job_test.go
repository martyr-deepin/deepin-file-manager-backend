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

func TestGetTemplateJob(t *testing.T) {
	Convey("Get template from directory which consists of directories", t, func() {
		uri := "./testdata"
		op := NewGetTemplateJob(uri)
		So(len(op.Execute()), ShouldEqual, 0)
	})
}
