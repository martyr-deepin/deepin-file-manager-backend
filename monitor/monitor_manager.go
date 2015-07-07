package monitor

import (
	"fmt"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/gio-2.0"
)

type MonitorManager struct {
	monitors map[MonitorID]*Monitor

	// monitor flags
	FileMonitorFlagsSendMoved      uint32
	FileMonitorFlagsWatchHardLinks uint32
	FileMonitorFlagsNone           uint32
	FileMonitorFlagsWatchMounts    uint32

	// monitor events
	FileMonitorEventMoved            uint32
	FileMonitorEventChanged          uint32
	FileMonitorEventCreated          uint32
	FileMonitorEventDeleted          uint32
	FileMonitorEventUnmounted        uint32
	FileMonitorEventPreUnmount       uint32
	FileMonitorEventAttributeChanged uint32
	FileMonitorEventChangesDoneHint  uint32
}

func NewMonitorManager() *MonitorManager {
	manager := &MonitorManager{
		monitors: map[MonitorID]*Monitor{},

		// monitor flags
		FileMonitorFlagsSendMoved:      uint32(gio.FileMonitorFlagsSendMoved),
		FileMonitorFlagsWatchHardLinks: uint32(gio.FileMonitorFlagsWatchHardLinks),
		FileMonitorFlagsNone:           uint32(gio.FileMonitorFlagsNone),
		FileMonitorFlagsWatchMounts:    uint32(gio.FileMonitorFlagsWatchMounts),

		// monitor events
		FileMonitorEventMoved:            uint32(gio.FileMonitorEventMoved),
		FileMonitorEventChanged:          uint32(gio.FileMonitorEventChanged),
		FileMonitorEventCreated:          uint32(gio.FileMonitorEventCreated),
		FileMonitorEventDeleted:          uint32(gio.FileMonitorEventDeleted),
		FileMonitorEventUnmounted:        uint32(gio.FileMonitorEventUnmounted),
		FileMonitorEventPreUnmount:       uint32(gio.FileMonitorEventPreUnmount),
		FileMonitorEventAttributeChanged: uint32(gio.FileMonitorEventAttributeChanged),
		FileMonitorEventChangesDoneHint:  uint32(gio.FileMonitorEventChangesDoneHint),
	}

	return manager
}

func (manager *MonitorManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.filemanager.Backend.Monitor",
		ObjectPath: "/com/deepin/filemanager/Backend/MonitorManager",
		Interface:  "com.deepin.filemanager.Backend.MonitorManager",
	}
}

func (manager *MonitorManager) Monitor(fileURI string, flags uint32) (string, dbus.ObjectPath, string) {
	monitor := NewMonitor(fileURI, gio.FileMonitorFlags(flags))
	if monitor == nil {
		return "", dbus.ObjectPath("/"), ""
	}
	err := dbus.InstallOnSession(monitor)
	if err != nil {
		fmt.Println(err)
		monitor.finalize()
		return "", dbus.ObjectPath("/"), ""
	}

	manager.monitors[MonitorID(monitor.ID)] = monitor
	dbusInfo := monitor.GetDBusInfo()
	return dbusInfo.Dest, dbus.ObjectPath(dbusInfo.ObjectPath), dbusInfo.Interface
}

func (manager *MonitorManager) Unmonitor(id uint32) {
	monitorID := MonitorID(id)
	monitor, ok := manager.monitors[MonitorID(monitorID)]
	if !ok {
		fmt.Printf("monitor %d not found", monitorID)
		return
	}

	fmt.Println("unmonitor", monitorID)
	monitor.finalize()
	delete(manager.monitors, MonitorID(monitorID))
}
