package delegator

import (
	"deepin-file-manager/operations"
	"net/url"
	"pkg.linuxdeepin.com/lib/dbus"
	"sync/atomic"
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
func NewListJob(path *url.URL, recusive bool, includeHidden bool) *ListJob {
	job := &ListJob{
		dbusInfo: genDBusInfo("ListJob", atomic.AddUint64(&_ListJobCount, 1)),
		op:       operations.NewListDirJob(path, recusive, includeHidden),
	}
	return job
}

// Execute list job.
func (job *ListJob) Execute() {
	job.op.ListenProcessedAmount(func(size int64, unit operations.AmountUnit) {
		dbus.Emit(job, "ProcessedAmount", size, uint16(unit))
	})
	job.op.ListenProperty(func(property operations.ListProperty) {
		dbus.Emit(job, "EntryInfo",
			property.DisplayName,
			property.BaseName,
			property.URI,
			property.MIME,
			property.Icon,
			property.Size,
			property.FileType,
			property.IsHidden,
			property.IsReadOnly,
			property.IsSymlink,
			property.CanDelete,
			property.CanExecute,
			property.CanRead,
			property.CanRename,
			property.CanTrash,
			property.CanWrite)
	})
	go func() {
		job.op.ListenDone(func(err error) {
			defer dbus.UnInstallObject(job)
			errMsg := ""
			if err != nil {
				errMsg = err.Error()
			}
			// time.Sleep(time.Microsecond * 200)
			dbus.Emit(job, "Done", errMsg)
		})
		job.op.Execute()
	}()
}

// Abort the job.
func (job *ListJob) Abort() {
	job.op.Abort()
	dbus.Emit(job, "Aborted")
	dbus.UnInstallObject(job)
}
