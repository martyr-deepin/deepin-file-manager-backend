package delegator

import (
	"encoding/json"
	"net/url"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/operations"
)

var (
	_GetLaunchAppJobCount uint64
)

// GetLaunchAppJob exports to dbus.
type GetLaunchAppJob struct {
	dbusInfo dbus.DBusInfo
	op       *operations.LaunchAppJob

	Done           func(string)
	LaunchAppInfos func(string)
}

// GetDBusInfo returns dbus information.
func (job *GetLaunchAppJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

// Execute GetLaunchAppJob.
func (job *GetLaunchAppJob) Execute() {
	job.op.ListenDone(func(err error) {
		defer dbus.UnInstallObject(job)
		if err != nil {
			dbus.Emit(job, "Done", err.Error())
			return
		}

		info := job.op.Result().(*operations.LaunchAppInfo)
		bInfos, err := json.Marshal(info)
		if err != nil {
			dbus.Emit(job, "Done", err.Error())
			return
		}

		dbus.Emit(job, "LaunchAppInfos", string(bInfos))
		dbus.Emit(job, "Done", "")
	})
	job.op.Execute()
}

// NewGetLaunchAppJob creates a new GetLaunchAppJob for dbus.
func NewGetLaunchAppJob(uri *url.URL) *GetLaunchAppJob {
	job := &GetLaunchAppJob{
		dbusInfo: genDBusInfo("GetLaunchAppJob", &_GetLaunchAppJobCount),
		op:       operations.NewLaunchAppJob(uri),
	}
	return job
}

// SetLaunchAppJob exports to dbus.
type SetLaunchAppJob struct {
	dbusInfo dbus.DBusInfo
	op       *operations.SetLaunchAppJob

	Done func(string)
}

// GetDBusInfo returns dbus information.
func (job *SetLaunchAppJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

// Execute SetLaunchAppJob.
func (job *SetLaunchAppJob) Execute() {
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

var (
	_SetLaunchAppJobCount uint64
)

// NewSetLaunchAppJob creates a new SetLaunchAppJob for dbus.
func NewSetLaunchAppJob(id string, mimeType string) *SetLaunchAppJob {
	job := &SetLaunchAppJob{
		dbusInfo: genDBusInfo("SetLaunchAppJob", &_SetLaunchAppJobCount),
		op:       operations.NewSetLaunchAppJob(id, mimeType),
	}
	return job
}
