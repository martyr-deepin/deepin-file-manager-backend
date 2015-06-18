package monitor

import (
	"pkg.linuxdeepin.com/lib/dbus"
)

func (monitor *TrashMonitor) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.filemanager.Backend.Monitor",
		ObjectPath: "/com/deepin/filemanager/Backend/Monitor",
		Interface:  "com.deepin.filemanager.Backend.TrashMonitor",
	}
}
