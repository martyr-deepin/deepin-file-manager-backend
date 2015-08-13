package delegator

import (
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/service/file-manager-backend/operations"
	// "time"
)

var (
	_ListJobCount uint64
)

// ListJob exports to dbus.
type ListJob struct {
	dbusInfo dbus.DBusInfo
	op       *operations.ListJob

	EntryInfo       func(string, string, string, string, string, int64, uint16, bool, bool, bool, bool, bool, bool, bool, bool, bool)
	Done            func(string)
	ProcessedAmount func(int64, uint16)
	Aborted         func()
}

// GetDBusInfo returns dbus information.
func (job *ListJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

// NewListJob creates a new list job for dbus.
func NewListJob(path string, flags operations.ListJobFlag) *ListJob {
	job := &ListJob{
		dbusInfo: genDBusInfo("ListJob", &_ListJobCount),
		op:       operations.NewListDirJob(path, flags),
	}
	return job
}

// Execute list job.
func (job *ListJob) Execute() []operations.ListProperty {
	defer dbus.UnInstallObject(job)
	var files []operations.ListProperty
	job.op.ListenProcessedAmount(func(size int64, unit operations.AmountUnit) {
		dbus.Emit(job, "ProcessedAmount", size, uint16(unit))
	})
	job.op.ListenProperty(func(property operations.ListProperty) {
		// TODO: read setting for icon size.
		icon := operations.GetThemeIcon(property.URI, 48)
		dbus.Emit(job, "EntryInfo",
			property.DisplayName,
			property.BaseName,
			property.URI,
			property.MIME,
			icon,
			property.Size,
			property.FileType,
			property.IsBackup,
			property.IsHidden,
			property.IsReadOnly,
			property.IsSymlink,
			property.CanDelete,
			property.CanExecute,
			property.CanRead,
			property.CanRename,
			property.CanTrash,
			property.CanWrite)
		files = append(files, property)
	})
	job.op.ListenDone(func(err error) {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		// time.Sleep(time.Microsecond * 200)
		dbus.Emit(job, "Done", errMsg)
	})
	job.op.Execute()

	return files
}

// Abort the job.
func (job *ListJob) Abort() {
	job.op.Abort()
	dbus.Emit(job, "Aborted")
	dbus.UnInstallObject(job)
}
