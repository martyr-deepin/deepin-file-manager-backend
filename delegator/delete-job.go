package delegator

import (
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/service/file-manager-backend/operations"
)

var (
	_DeleteJobCount uint64
)

// DeleteJob exports to dbus.
type DeleteJob struct {
	dbusInfo dbus.DBusInfo
	uris     []string
	op       *operations.DeleteJob

	Done             func()
	TotalAmount      func(int64, uint16)
	ProcessedAmount  func(int64, uint16)
	ProcessedPercent func(int64)
	Aborted          func()
	Deleting         func(string)
}

// GetDBusInfo returns dbus information.
func (job *DeleteJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

// NewDeleteJob creates a new delete job for dbus.
func NewDeleteJob(urls []string, shouldConfirm bool, uiDelegate IUIDelegate) *DeleteJob {
	job := &DeleteJob{
		dbusInfo: genDBusInfo("DeleteJob", &_DeleteJobCount),
		op:       operations.NewDeleteJob(urls, shouldConfirm, uiDelegate),
	}

	return job
}

func (job *DeleteJob) listenSignals() {
	job.op.ListenTotalAmount(func(amount int64, unit operations.AmountUnit) {
		dbus.Emit(job, "TotalAmount", amount, uint16(unit))
	})
	job.op.ListenProcessedAmount(func(size int64, unit operations.AmountUnit) {
		dbus.Emit(job, "ProcessedAmount", size, uint16(unit))
	})
	job.op.ListenPercent(func(percent int64) {
		dbus.Emit(job, "ProcessedPercent", percent)
	})

	job.op.ListenDeleting(func(deletingURL string) {
		dbus.Emit(job, "Deleting", deletingURL)
	})
}

func (job *DeleteJob) executeJob() {
	defer dbus.UnInstallObject(job)
	job.op.Execute()
	dbus.Emit(job, "Done")
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
