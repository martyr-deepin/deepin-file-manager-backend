package delegator

import (
	"pkg.linuxdeepin.com/lib/dbus"
)

func (*QueryFileInfoJob) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       JobDestination,
		ObjectPath: JobObjectPath + "/FileInfo",
		Interface:  JobDestination + ".FileInfo",
	}
}
