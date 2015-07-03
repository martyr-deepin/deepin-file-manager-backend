package delegator

import (
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/operations"
)

var (
	_RenameJobCount uint64 = 0
)

type RenameJob struct {
	dbusInfo dbus.DBusInfo
	op       *operations.RenameJob

	Done    func(string)
	OldName func(string)
}

func (job *RenameJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

func (job *RenameJob) Execute() {
	defer dbus.UnInstallObject(job)
	job.op.ListenDone(func(err error) {
		var errMsg = ""
		if err != nil {
			errMsg = err.Error()
		}
		dbus.Emit(job, "Done", errMsg)
	})
	job.op.ListenOldName(func(oldName string) {
		dbus.Emit(job, "OldName", oldName)
	})
	job.op.ListenNewFile(func(fileURL string) {
		dbus.Emit(job, "NewFile", fileURL)
	})
	job.op.Execute()
}

func NewRenameJob(fileURL string, newName string) *RenameJob {
	job := &RenameJob{
		dbusInfo: genDBusInfo("RenameJob", &_RenameJobCount),
		op:       operations.NewRenameJob(fileURL, newName),
	}
	return job
}
