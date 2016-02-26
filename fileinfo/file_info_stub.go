/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package fileinfo

import (
	"pkg.deepin.io/lib/dbus"
)

func (*QueryFileInfoJob) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.filemanager.Backend.FileInfo",
		ObjectPath: "/com/deepin/filemanager/Backend/FileInfo",
		Interface:  "com.deepin.filemanager.Backend.FileInfo",
	}
}
