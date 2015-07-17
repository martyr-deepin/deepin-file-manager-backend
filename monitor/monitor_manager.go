package monitor

import (
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/gio-2.0"
	"sync/atomic"
)

var _MonitorCounter uint32
var _WatcherCounter uint32

type MonitorManager struct {
	monitors map[MonitorID]*Monitor
	watchers map[WatcherID]*Watcher

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

	// watcher events
	FsNotifyCreated          uint32
	FsNotifyDeleted          uint32
	FsNotifyModified         uint32
	FsNotifyRename           uint32
	FsNotifyAttributeChanged uint32
}

func NewMonitorManager() *MonitorManager {
	manager := &MonitorManager{
		monitors: map[MonitorID]*Monitor{},
		watchers: map[WatcherID]*Watcher{},

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

		// watcher events
		FsNotifyCreated:          uint32(FsNotifyCreated),
		FsNotifyDeleted:          uint32(FsNotifyDeleted),
		FsNotifyModified:         uint32(FsNotifyModified),
		FsNotifyRename:           uint32(FsNotifyRename),
		FsNotifyAttributeChanged: uint32(FsNotifyAttributeChanged),
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

func (manager *MonitorManager) Monitor(fileURI string, flags uint32) (string, dbus.ObjectPath, string, error) {
	monitorID := atomic.AddUint32(&_MonitorCounter, 1)
	monitor, err := NewMonitor(monitorID, fileURI, gio.FileMonitorFlags(flags))
	if monitor == nil {
		return "", dbus.ObjectPath("/"), "", err
	}

	if err := dbus.InstallOnSession(monitor); err != nil {
		fmt.Println(err)
		monitor.finalize()
		return "", dbus.ObjectPath("/"), "", err
	}

	manager.monitors[MonitorID(monitorID)] = monitor
	dbusInfo := monitor.GetDBusInfo()
	return dbusInfo.Dest, dbus.ObjectPath(dbusInfo.ObjectPath), dbusInfo.Interface, nil
}

func (manager *MonitorManager) Unmonitor(id uint32) {
	monitorID := MonitorID(id)
	monitor, ok := manager.monitors[monitorID]
	if !ok {
		fmt.Printf("monitor %d not found", monitorID)
		return
	}

	fmt.Println("unmonitor", monitorID)
	dbus.UnInstallObject(monitor)
	monitor.finalize()
	delete(manager.monitors, monitorID)
}

func (manager *MonitorManager) Watch(fileURI string) (string, dbus.ObjectPath, string, error) {
	watcherID := atomic.AddUint32(&_WatcherCounter, 1)
	watcher, err := NewWatcher(watcherID, fileURI)
	if err != nil {
		return "", dbus.ObjectPath("/"), "", err
	}

	if err := dbus.InstallOnSession(watcher); err != nil {
		watcher.finalize()
		return "", dbus.ObjectPath("/"), "", err
	}

	manager.watchers[WatcherID(watcherID)] = watcher
	dbusInfo := watcher.GetDBusInfo()
	return dbusInfo.Dest, dbus.ObjectPath(dbusInfo.ObjectPath), dbusInfo.Interface, nil
}

func (manager *MonitorManager) Unwatcher(id uint32) {
	watcherID := WatcherID(id)
	watcher, ok := manager.watchers[watcherID]
	if !ok {
		fmt.Println("watcher %d not found", watcherID)
		return
	}

	fmt.Println("unwatcher", watcherID)
	dbus.UnInstallObject(watcher)
	watcher.finalize()
	delete(manager.watchers, watcherID)
}
