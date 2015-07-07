package delegator

import (
	"pkg.deepin.io/lib/dbus"
)

func (*QueryFileInfoJob) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       JobDestination,
		ObjectPath: JobObjectPath + "/FileInfo",
		Interface:  JobDestination + ".FileInfo",
	}
}
