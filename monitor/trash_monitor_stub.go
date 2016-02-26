/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package monitor

import (
	"pkg.deepin.io/lib/dbus"
)

func (monitor *TrashMonitor) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.filemanager.Backend.Monitor",
		ObjectPath: "/com/deepin/filemanager/Backend/Monitor/TrashMonitor",
		Interface:  "com.deepin.filemanager.Backend.Monitor.TrashMonitor",
	}
}
