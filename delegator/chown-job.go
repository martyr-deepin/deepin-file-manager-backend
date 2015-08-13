package delegator

import (
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/service/file-manager-backend/operations"
)

var (
	_ChownJobCount uint64
)

// ChownJob exports to dbus.
type ChownJob struct {
	dbusInfo dbus.DBusInfo
	op       *operations.ChownJob

	Done func(string)
}

// GetDBusInfo returns dbus information.
func (job *ChownJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

// Execute chown job.
func (job *ChownJob) Execute() {
	job.op.ListenDone(func(err error) {
		defer dbus.UnInstallObject(job)
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		dbus.Emit(job, "Done", errMsg)
	})
	job.op.Execute()
}

// NewChownJob creates a new chown job for dbus.
func NewChownJob(uri string, newOwner string, newGroup string) *ChownJob {
	job := &ChownJob{
		dbusInfo: genDBusInfo("ChownJob", &_ChownJobCount),
		op:       operations.NewChownJob(uri, newOwner, newGroup),
	}
	return job
}
