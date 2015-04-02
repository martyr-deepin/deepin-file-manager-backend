package delegator

import (
	"deepin-file-manager/operations"
	"net/url"
	"pkg.linuxdeepin.com/lib/dbus"
	"sync/atomic"
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

	Done func(string)
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

		dbus.Emit(job, "Done", errMsg)
		// TODO:
		// job.commandRecorder
		// operations.FileUndoManagerInstance().RecordJob(create, job.op)
	}()
}

// NewCreateFileJob creates a new create job to create a new file.
func NewCreateFileJob(destDirURL *url.URL, filename string, initContent []byte, uiDelegate IUIDelegate) *CreateJob {
	job := &CreateJob{
		dbusInfo: genDBusInfo("CreateFileJob", atomic.AddUint64(&_CreateFileJobCount, 1)),
		op:       operations.NewCreateFileJob(destDirURL, filename, initContent, uiDelegate),
	}
	return job
}

// NewCreateFileFromTemplateJob creates a new create job to create a new file from a template.
func NewCreateFileFromTemplateJob(uri *url.URL, tempateURL *url.URL, uiDelegate IUIDelegate) *CreateJob {
	job := &CreateJob{
		dbusInfo: genDBusInfo("CreateFileFromTemplateJob", atomic.AddUint64(&_CreateTemplateJobCount, 1)),
		op:       operations.NewCreateFileFromTemplateJob(uri, tempateURL, uiDelegate),
	}
	return job
}

// NewCreateDirectoryJob creates a new create job to create a new directory.
func NewCreateDirectoryJob(destDirURL *url.URL, dirname string, uiDelegate IUIDelegate) *CreateJob {
	job := &CreateJob{
		dbusInfo: genDBusInfo("CreateDirJob", atomic.AddUint64(&_CreateDirJobCount, 1)),
		op:       operations.NewCreateDirectoryJob(destDirURL, dirname, uiDelegate),
	}
	return job
}

// NewLinkJob creates a new job to creates a link.
func NewLinkJob(srcURL *url.URL, destDirURL *url.URL, uiDelegate IUIDelegate) *CreateJob {
	job := &CreateJob{
		dbusInfo: genDBusInfo("CreateLinkJob", atomic.AddUint64(&_CreateLinkJobCount, 1)),
		op:       operations.NewCreateLinkJob(srcURL, destDirURL, uiDelegate),
	}
	return job
}
