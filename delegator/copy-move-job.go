package delegator

import (
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/gio-2.0"
	. "pkg.deepin.io/service/file-manager-backend/log"
	"pkg.deepin.io/service/file-manager-backend/operations"
)

var (
	_CopyJobCount uint64
	_MoveJobCount uint64
)

// CopyJob exports to dbus.
type CopyJob struct {
	dbusInfo dbus.DBusInfo

	uris []string
	op   *operations.CopyMoveJob

	Done            func(string)
	ProcessedAmount func(int64, uint16)
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
		Log.Debug("copying", srcURL)
		dbus.Emit(job, "Copying", srcURL)
	})
	job.op.ListenCreatingDir(func(dirURL string) {
		// TODO
		// dbus.Emit(job, "CreatingDir", dirURL)
	})
	job.op.ListenCopyingMovingDone(func(srcURL string, destURL string) {
	})
	job.op.ListenDone(func(err error) {
		defer dbus.UnInstallObject(job)
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		dbus.Emit(job, "Done", errMsg)
	})
}

// Execute copy job.
func (job *CopyJob) Execute() {
	job.listenSignals()
	// TODO
	// operations.FileUndoManagerInstance().RecordJob(COPY, job.op)
	job.op.Execute()
}

// NewCopyJob creates a new copy job for dbus.
func NewCopyJob(srcUrls []string, destDirURL string, targetName string, flags uint32, uiDelegate IUIDelegate) *CopyJob {
	job := &CopyJob{
		dbusInfo: genDBusInfo("CopyJob", &_CopyJobCount),
		op:       operations.NewCopyJob(srcUrls, destDirURL, targetName, gio.FileCopyFlags(flags), uiDelegate),
	}
	return job
}

// MoveJob exports to dbus.
type MoveJob struct {
	dbusInfo dbus.DBusInfo

	uris []string
	op   *operations.CopyMoveJob

	Done            func(string)
	ProcessedAmount func(int64, uint16)
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
	job.op.ListenDone(func(err error) {
		defer dbus.UnInstallObject(job)
		errMsg := ""
		if errMsg != "" {
			errMsg = err.Error()
		}
		dbus.Emit(job, "Done", errMsg)
	})
}

// Execute move job.
func (job *MoveJob) Execute() {
	job.listenSignals()
	job.op.Execute()
}

// Abort the job.
func (job *MoveJob) Abort() {
	job.op.Abort()
	dbus.Emit(job, "Aborted")
}

// NewMoveJob creates a new move job for dbus.
func NewMoveJob(srcUrls []string, destDirURL string, targetName string, flags uint32, uiDelegate IUIDelegate) *MoveJob {
	job := &MoveJob{
		dbusInfo: genDBusInfo("MoveJob", &_MoveJobCount),
		op:       operations.NewMoveJob(srcUrls, destDirURL, targetName, gio.FileCopyFlags(flags), uiDelegate),
	}
	return job
}
