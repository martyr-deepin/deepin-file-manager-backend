package delegator

import (
	"deepin-file-manager/operations"
	"net/url"
	"pkg.linuxdeepin.com/lib/dbus"
	"sync/atomic"
)

var (
	_ChmodJobCount uint64
)

// ChmodJob exports to dbus.
type ChmodJob struct {
	dbusInfo dbus.DBusInfo
	op       *operations.ChmodJob

	Done func(string)
}

func (job *ChmodJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

// Execute chmod job.
func (job *ChmodJob) Execute() {
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

// NewChmodJob creates a new ChmodJob for dbus.
func NewChmodJob(uri *url.URL, permission uint32) *ChmodJob {
	job := &ChmodJob{
		dbusInfo: genDBusInfo("ChmodJob", atomic.AddUint64(&_ChmodJobCount, 1)),
		op:       operations.NewChmodJob(uri, permission),
	}
	return job
}
