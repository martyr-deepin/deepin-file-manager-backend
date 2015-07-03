package delegator

import (
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/operations"
	"sync"
)

var (
	_TrashJobCount     uint64
	_TrashJobCountLock sync.Mutex
)

// TrashJob exports to dbus.
type TrashJob struct {
	dbusInfo dbus.DBusInfo
	uris     []string
	op       *operations.DeleteJob

	Done            func()
	Trashing        func(string)
	Deleting        func(string)
	ProcessedAmount func(int64, uint16)
	Aborted         func()
}

// GetDBusInfo returns dbus information.
func (job *TrashJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

// NewTrashJob creates a new trash job for dbus.
func NewTrashJob(urls []string, shouldConfirmTrash bool, uiDelegate IUIDelegate) *TrashJob {
	return &TrashJob{
		dbusInfo: genDBusInfo("TrashJob", &_TrashJobCount),
		op:       operations.NewTrashJob(urls, shouldConfirmTrash, uiDelegate),
	}
}

// Execute trash job.
func (job *TrashJob) Execute() {
	job.op.ListenProcessedAmount(func(size int64, unit operations.AmountUnit) {
		dbus.Emit(job, "ProcessedAmount", size, uint16(unit))
	})
	// TODO: fill signals.
	defer dbus.UnInstallObject(job)
	job.op.Execute()
	dbus.Emit(job, "Done")
}

// Abort the job.
func (job *TrashJob) Abort() {
	job.op.Abort()
}
