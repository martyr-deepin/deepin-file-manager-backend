package monitor

import (
	"fmt"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"sync/atomic"
)

type MonitorID uint32

var _MonitorCounter uint32

type Monitor struct {
	file        *gio.File
	cancellable *gio.Cancellable
	flags       gio.FileMonitorFlags
	monitor     *gio.FileMonitor
	dbusInfo    dbus.DBusInfo

	ID uint32

	Changed func(string, string, uint32)
}

func NewMonitor(fileURI string, flags gio.FileMonitorFlags) *Monitor {
	file := gio.FileNewForCommandlineArg(fileURI)
	if file == nil {
		fmt.Printf("create file from %s failed\n", fileURI)
		return nil
	}

	monitor := &Monitor{
		file:        file,
		cancellable: gio.NewCancellable(),
		flags:       flags,
		ID:          atomic.AddUint32(&_MonitorCounter, 1),
	}

	monitor.dbusInfo = dbus.DBusInfo{
		Dest:       "com.deepin.filemanager.Backend.Monitor",
		ObjectPath: "/com/deepin/filemanager/Backend/Monitor",
		Interface:  fmt.Sprintf("com.deepin.filemanager.Backend.Monitor.m%d", monitor.ID),
	}

	var err error
	monitor.monitor, err = monitor.file.Monitor(flags, monitor.cancellable)
	if err != nil {
		fmt.Println("create file monitor failed:", err)
		monitor.finalize()
		return nil
	}
	monitor.monitor.Connect("changed", func(m *gio.FileMonitor, file *gio.File, newFile *gio.File, events gio.FileMonitorEvent) {
		newFileURI := ""
		if newFile != nil {
			newFileURI = newFile.GetUri()
		}
		err := dbus.Emit(monitor, "Changed", file.GetUri(), newFileURI, uint32(events))
		if err != nil {
			fmt.Println("emit signal failed:", err)
		}
	})

	return monitor
}

func (monitor *Monitor) GetDBusInfo() dbus.DBusInfo {
	return monitor.dbusInfo
}

func (monitor *Monitor) SetRateLimit(msecs int32) {
	monitor.monitor.SetRateLimit(msecs)
}

func (monitor *Monitor) Cancel() {
	monitor.monitor.Cancel()
}

func (monitor *Monitor) IsCancelled() bool {
	return monitor.monitor.IsCancelled()
}

func (monitor *Monitor) finalize() {
	if !monitor.cancellable.IsCancelled() {
		monitor.cancellable.Cancel()
	}
	monitor.cancellable.Unref()

	if monitor.monitor != nil {
		monitor.monitor.Unref()
	}

	if monitor.file != nil {
		monitor.file.Unref()
	}
}
