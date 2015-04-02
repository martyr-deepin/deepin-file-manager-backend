package delegator

import (
	"deepin-file-manager/operations"
	"net/url"
	"pkg.linuxdeepin.com/lib/dbus"
	"sync/atomic"
)

var (
	_CopyJobCount uint64
	_MoveJobCount uint64
)

// CopyJob exports to dbus.
type CopyJob struct {
	dbusInfo dbus.DBusInfo

	uris []*url.URL
	op   *operations.CopyMoveJob

	Done            func()
	ProcessedAmount func(uint64, uint16)
	Copying         func(string)
	Aborted         func()
}

// GetDBusInfo returns dbus information.
func (job *CopyJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

// Abort the job.
func (job *CopyJob) Abort() {
	job.op.Abort()
	dbus.Emit(job, "Aborted")
}

func (job *CopyJob) listenSignals() {
	job.op.ListenProcessedAmount(func(amount int64, unit operations.AmountUnit) {
		dbus.Emit(job, "ProcessedAmount", amount, uint16(unit))
	})
	job.op.ListenCopying(func(srcURL string) {
		dbus.Emit(job, "Copying", srcURL)
	})
	job.op.ListenCreatingDir(func(string) {
		// TODO: why this signal???
		// dbus.Emit(job, "CreatingDir", dirURL)
	})
}

// Execute copy job.
func (job *CopyJob) Execute() {
	job.listenSignals()
	go func() {
		defer dbus.UnInstallObject(job)
		// TODO
		// operations.FileUndoManagerInstance().RecordJob(COPY, job.op)
		job.op.Execute()
		dbus.Emit(job, "Done")
	}()
}

// NewCopyJob creates a new copy job for dbus.
func NewCopyJob(srcUrls []*url.URL, destDirURL *url.URL, targetName string, uiDelegate IUIDelegate) *CopyJob {
	job := &CopyJob{
		dbusInfo: genDBusInfo("CopyJob", atomic.AddUint64(&_CopyJobCount, 1)),
		op:       operations.NewCopyJob(srcUrls, destDirURL, targetName, uiDelegate),
	}
	return job
}

// MoveJob exports to dbus.
type MoveJob struct {
	dbusInfo dbus.DBusInfo

	uris []*url.URL
	op   *operations.CopyMoveJob

	Done            func()
	ProcessedAmount func(uint64, uint16)
	Moving          func(string)
	Aborted         func()
}

// GetDBusInfo returns dbus information.
func (job *MoveJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

func (job *MoveJob) listenSignals() {
	job.op.ListenProcessedAmount(func(amount int64, unit operations.AmountUnit) {
		dbus.Emit(job, "ProcessedAmount", amount, uint16(unit))
	})
	job.op.ListenMoving(func(srcURL string) {
		dbus.Emit(job, "Moving", srcURL)
	})
	job.op.ListenCopyingMovingDone(func(srcURL string, destURL string) {
		// dbus.Emit(job, "MovingDone", srcURL, destURL)
	})
	job.op.ListenCreatingDir(func(dirURL string) {
		// dbus.Emit(job, "CreatingDir", dirURL)
	})
}

// Execute move job.
func (job *MoveJob) Execute() {
	job.listenSignals()
	go func() {
		job.op.Execute()
		dbus.Emit(job, "Done")
		dbus.UnInstallObject(job)
	}()
}

// Abort the job.
func (job *MoveJob) Abort() {
	job.op.Abort()
	dbus.Emit(job, "Aborted")
}

// NewMoveJob creates a new move job for dbus.
func NewMoveJob(srcUrls []*url.URL, destDirURL *url.URL, targetName string, uiDelegate IUIDelegate) *MoveJob {
	job := &MoveJob{
		dbusInfo: genDBusInfo("MoveJob", atomic.AddUint64(&_MoveJobCount, 1)),
		op:       operations.NewMoveJob(srcUrls, destDirURL, targetName, uiDelegate),
	}
	return job
}
