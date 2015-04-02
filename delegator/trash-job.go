package delegator

import (
	"deepin-file-manager/operations"
	"net/url"
	"pkg.linuxdeepin.com/lib/dbus"
	"sync"
)

var (
	_TrashJobCount     uint64
	_TrashJobCountLock sync.Mutex
)

// TrashJob exports to dbus.
type TrashJob struct {
	dbusInfo dbus.DBusInfo
	uris     []*url.URL
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
func NewTrashJob(urls []*url.URL, shouldConfirmTrash bool, uiDelegate IUIDelegate) *TrashJob {
	return &TrashJob{
		dbusInfo: genDBusInfo("TrashJob", _TrashJobCount),
		op:       operations.NewTrashJob(urls, shouldConfirmTrash, uiDelegate),
	}
}

// Execute trash job.
func (job *TrashJob) Execute() {
	job.op.ListenProcessedAmount(func(size int64, unit operations.AmountUnit) {
		dbus.Emit(job, "ProcessedAmount", size, uint16(unit))
	})
	// TODO: fill signals.
	go func() {
		defer dbus.UnInstallObject(job)
		job.op.Execute()
		dbus.Emit(job, "Done")
	}()
}

// Abort the job.
func (job *TrashJob) Abort() {
	job.op.Abort()
}

var (
	_EmptyTrashJobCount     uint64
	_EmptyTrashJobCountLock sync.Mutex
)

// EmptyTrashJob for dbus.
type EmptyTrashJob struct {
	dbusInfo dbus.DBusInfo
	op       *operations.EmptyTrashJob

	Trashing        func(string)
	Deleting        func(string)
	Done            func()
	ProcessedAmount func(int64, uint16)
	Aborted         func()
}

// GetDBusInfo returns dbus information.
func (job *EmptyTrashJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

// NewEmptyTrashJob creates a new EmptyTrashJob for dbus.
func NewEmptyTrashJob(shouldConfirmTrash bool, uiDelegate IUIDelegate) *EmptyTrashJob {
	_EmptyTrashJobCountLock.Lock()
	defer _EmptyTrashJobCountLock.Unlock()
	job := &EmptyTrashJob{
		dbusInfo: genDBusInfo("EmptyTrashJob", _EmptyTrashJobCount),
		op:       operations.NewEmptyTrashJob(shouldConfirmTrash, uiDelegate),
	}
	_EmptyTrashJobCount++
	return job
}

// TODO:
func (job *EmptyTrashJob) listenSignals() {
}

func (job *EmptyTrashJob) executeJob() {
	go func() {
		defer dbus.UnInstallObject(job)
		job.op.Execute()
		dbus.Emit(job, "Done")
	}()
}

// Execute empty trash job.
func (job *EmptyTrashJob) Execute() {
	job.listenSignals()
	job.executeJob()
}
