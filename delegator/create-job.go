package delegator

import (
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/service/file-manager-backend/operations"
)

var (
	_CreateFileJobCount     uint64
	_CreateTemplateJobCount uint64
	_CreateDirJobCount      uint64
	_CreateLinkJobCount     uint64
)

// CreateJob exports to dbus.
type CreateJob struct {
	dbusInfo        dbus.DBusInfo
	op              *operations.CreateJob
	commandRecorder *operations.CommandRecorder

	Done func(string, string)
}

// GetDBusInfo returns dbus information.
func (job *CreateJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

// Execute create job.
func (job *CreateJob) Execute() {
	go func() {
		defer dbus.UnInstallObject(job)
		err := job.op.Execute()
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}

		dbus.Emit(job, "Done", job.op.Result().(string), errMsg)
		// TODO:
		// job.commandRecorder
		// operations.FileUndoManagerInstance().RecordJob(create, job.op)
	}()
}

// NewCreateFileJob creates a new create job to create a new file.
func NewCreateFileJob(destDirURL string, filename string, initContent []byte, uiDelegate IUIDelegate) *CreateJob {
	job := &CreateJob{
		dbusInfo: genDBusInfo("CreateFileJob", &_CreateFileJobCount),
		op:       operations.NewCreateFileJob(destDirURL, filename, initContent, uiDelegate),
	}
	return job
}

// NewCreateFileFromTemplateJob creates a new create job to create a new file from a template.
func NewCreateFileFromTemplateJob(uri string, tempateURL string, uiDelegate IUIDelegate) *CreateJob {
	job := &CreateJob{
		dbusInfo: genDBusInfo("CreateFileFromTemplateJob", &_CreateTemplateJobCount),
		op:       operations.NewCreateFileFromTemplateJob(uri, tempateURL, uiDelegate),
	}
	return job
}

// NewCreateDirectoryJob creates a new create job to create a new directory.
func NewCreateDirectoryJob(destDirURL string, dirname string, uiDelegate IUIDelegate) *CreateJob {
	job := &CreateJob{
		dbusInfo: genDBusInfo("CreateDirJob", &_CreateDirJobCount),
		op:       operations.NewCreateDirectoryJob(destDirURL, dirname, uiDelegate),
	}
	return job
}

// NewLinkJob creates a new job to creates a link.
func NewLinkJob(srcURL string, destDirURL string, uiDelegate IUIDelegate) *CreateJob {
	job := &CreateJob{
		dbusInfo: genDBusInfo("CreateLinkJob", &_CreateLinkJobCount),
		op:       operations.NewCreateLinkJob(srcURL, destDirURL, uiDelegate),
	}
	return job
}
