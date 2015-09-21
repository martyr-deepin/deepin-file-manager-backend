package monitor

import (
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/gio-2.0"
	. "pkg.deepin.io/service/file-manager-backend/log"
)

type MonitorID uint32

type Monitor struct {
	file        *gio.File
	cancellable *gio.Cancellable
	flags       gio.FileMonitorFlags
	monitor     *gio.FileMonitor
	dbusInfo    dbus.DBusInfo

	ID uint32

	Changed func(string, string, uint32)
}

func NewMonitor(id uint32, fileURI string, flags gio.FileMonitorFlags) (*Monitor, error) {
	file := gio.FileNewForCommandlineArg(fileURI)
	if file == nil {
		return nil, fmt.Errorf("create file from %s failed\n", fileURI)
	}

	monitor := &Monitor{
		file:        file,
		cancellable: gio.NewCancellable(),
		flags:       flags,
		ID:          id,
	}

	monitor.dbusInfo = dbus.DBusInfo{
		Dest:       "com.deepin.filemanager.Backend.Monitor",
		ObjectPath: fmt.Sprintf("/com/deepin/filemanager/Backend/Monitor/%d", monitor.ID),
		Interface:  "com.deepin.filemanager.Backend.Monitor",
	}

	var err error
	monitor.monitor, err = monitor.file.Monitor(flags, monitor.cancellable)
	if err != nil {
		monitor.finalize()
		return nil, fmt.Errorf("create file monitor failed: %s", err)
	}
	monitor.monitor.Connect("changed", func(m *gio.FileMonitor, file *gio.File, newFile *gio.File, events gio.FileMonitorEvent) {
		newFileURI := ""
		if newFile != nil {
			newFileURI = newFile.GetUri()
		}
		err := dbus.Emit(monitor, "Changed", file.GetUri(), newFileURI, uint32(events))
		if err != nil {
			Log.Warning("emit signal failed:", err)
		}
	})

	return monitor, nil
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
