package delegator

import (
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/operations"
	"sync"
)

var (
	_StatJobCount     uint64
	_StatJobCountLock sync.Mutex
)

// StatJob exports to dbus.
type StatJob struct {
	dbusInfo dbus.DBusInfo
	op       *operations.StatJob

	Stat func(operations.StatProperty)
	Done func(string)
}

// GetDBusInfo returns dbus information.
func (job *StatJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

// Execute stat job.
func (job *StatJob) Execute() {
	job.op.ListenDone(func(err error) {
		defer dbus.UnInstallObject(job)
		if err != nil {
			dbus.Emit(job, "Done", err.Error())
			return
		}

		dbus.Emit(job, "Stat", job.op.Result().(operations.StatProperty))
		dbus.Emit(job, "Done", "")
	})
	job.op.Execute()
}

// NewStatJob creates a new stat job for dbus.
func NewStatJob(uri string) *StatJob {
	_StatJobCountLock.Lock()
	defer _StatJobCountLock.Unlock()
	job := &StatJob{
		dbusInfo: genDBusInfo("StatJob", &_StatJobCount),
		op:       operations.NewStatJob(uri),
	}
	_StatJobCount++
	return job
}
