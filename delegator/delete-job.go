package delegator

import (
	"deepin-file-manager/operations"
	"net/url"
	"pkg.linuxdeepin.com/lib/dbus"
	"sync/atomic"
)

var (
	_DeleteJobCount uint64
)

// DeleteJob exports to dbus.
type DeleteJob struct {
	dbusInfo dbus.DBusInfo
	uris     []*url.URL
	op       *operations.DeleteJob

	Done            func()
	ProcessedAmount func(int64, uint16)
	Aborted         func()
	Deleting        func(string)
}

// GetDBusInfo returns dbus information.
func (job *DeleteJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

// NewDeleteJob creates a new delete job for dbus.
func NewDeleteJob(urls []*url.URL, shouldConfirmTrash bool, uiDelegate IUIDelegate) *DeleteJob {
	job := &DeleteJob{
		dbusInfo: genDBusInfo("DeleteJob", atomic.AddUint64(&_DeleteJobCount, 1)),
		op:       operations.NewDeleteJob(urls, shouldConfirmTrash, uiDelegate),
	}

	return job
}

func (job *DeleteJob) listenSignals() {
	job.op.ListenProcessedAmount(func(size int64, unit operations.AmountUnit) {
		dbus.Emit(job, "ProcessedAmount", size, uint16(unit))
	})

	job.op.ListenDeleting(func(deletingURL string) {
		dbus.Emit(job, "Deleting", deletingURL)
	})
}

func (job *DeleteJob) executeJob() {
	go func() {
		defer dbus.UnInstallObject(job)
		job.op.Execute()
		dbus.Emit(job, "Done")
	}()
}

// Execute delete job.
func (job *DeleteJob) Execute() {
	job.listenSignals()
	job.executeJob()
}

// Abort the job.
func (job *DeleteJob) Abort() {
	job.op.Abort()
	dbus.Emit(job, "Aborted")
}
