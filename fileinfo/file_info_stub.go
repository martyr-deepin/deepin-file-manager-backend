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
