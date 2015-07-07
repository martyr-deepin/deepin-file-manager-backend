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
