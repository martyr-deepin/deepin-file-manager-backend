package delegator

import (
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/service/file-manager-backend/operations"
)

var _GetTemplateJobCount uint64

type GetTemplateJob struct {
	dbusInfo dbus.DBusInfo
	op       *operations.GetTemplateJob
}

// GetDBusInfo returns dbus information.
func (job *GetTemplateJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

// NewListJob creates a new list job for dbus.
func NewGetTemplateJob(templateDirURI string) *GetTemplateJob {
	job := &GetTemplateJob{
		dbusInfo: genDBusInfo("GetTemplateJob", &_GetTemplateJobCount),
		op:       operations.NewGetTemplateJob(templateDirURI),
	}
	return job
}

func (job *GetTemplateJob) Execute() []string {
	defer dbus.UnInstallObject(job)
	return job.op.Execute()
}
