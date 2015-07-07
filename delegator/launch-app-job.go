package delegator

import (
	"encoding/json"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/gio-2.0"
	"pkg.deepin.io/lib/operations"
)

var (
	_GetDefaultLaunchAppJobCount     uint64
	_GetRecommendedLaunchAppJobCount uint64
	_GetAllLaunchAppJobCount         uint64
)

// GetDefaultLaunchAppJob exports to dbus.
type GetDefaultLaunchAppJob struct {
	dbusInfo dbus.DBusInfo
	op       *operations.GetDefaultLaunchAppJob

	Done                 func(string)
	DefaultLaunchAppInfo func(string)
}

// GetDBusInfo returns dbus information.
func (job *GetDefaultLaunchAppJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

type DefaultLaunchAppInfo struct {
	Name string
	Id   string
	// Icon string
}

// Execute GetDefaultLaunchAppJob.
func (job *GetDefaultLaunchAppJob) Execute() {
	job.op.ListenDone(func(err error) {
		defer dbus.UnInstallObject(job)
		if err != nil {
			dbus.Emit(job, "Done", err.Error())
			return
		}

		app := job.op.Result().(*gio.AppInfo)
		info := &DefaultLaunchAppInfo{
			Name: app.GetName(),
			Id:   app.GetId(),
			// TODO: get icon.
			// Icon: getIconFromApp(app),
		}
		app.Unref()
		bInfos, err := json.Marshal(info)
		if err != nil {
			dbus.Emit(job, "Done", err.Error())
			return
		}

		dbus.Emit(job, "DefaultAppInfo", string(bInfos))
		dbus.Emit(job, "Done", "")
	})
	job.op.Execute()
}

// NewGetDefaultLaunchAppJob creates a new GetLaunchAppJob for dbus.
func NewGetDefaultLaunchAppJob(uri string, mustSupportURI bool) *GetDefaultLaunchAppJob {
	job := &GetDefaultLaunchAppJob{
		dbusInfo: genDBusInfo("GetDefaultLaunchAppJob", &_GetDefaultLaunchAppJobCount),
		op:       operations.NewGetDefaultLaunchAppJob(uri, mustSupportURI),
	}
	return job
}

type LaunchAppInfo struct {
	Names []string
	Ids   []string
	// Icons []string
}

type GetRecommendedLaunchAppsJob struct {
	dbusInfo dbus.DBusInfo
	op       *operations.GetRecommendedLaunchAppsJob

	Done                     func(string)
	RecommendedLaunchAppInfo func(string)
}

func (job *GetRecommendedLaunchAppsJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

func NewGetRecommendedLaunchAppsJob(uri string) *GetRecommendedLaunchAppsJob {
	return &GetRecommendedLaunchAppsJob{
		dbusInfo: genDBusInfo("GetRecommendedLaunchAppsJob", &_GetRecommendedLaunchAppJobCount),
		op:       operations.NewGetRecommendedLaunchAppsJob(uri),
	}
}

func (job *GetRecommendedLaunchAppsJob) Execute() {
	job.op.ListenDone(func(e error) {
		defer dbus.UnInstallObject(job)
		if e != nil {
			dbus.Emit(job, "Done", e.Error())
			return
		}

		apps := job.op.Result().([]*gio.AppInfo)
		info := LaunchAppInfo{
			Names: make([]string, len(apps)),
			Ids:   make([]string, len(apps)),
		}
		for i, app := range apps {
			info.Names[i] = app.GetName()
			info.Ids[i] = app.GetId()
			// TODO: get icon
			// app.GetIcon()
			app.Unref()
		}

		bInfos, err := json.Marshal(info)
		if err != nil {
			dbus.Emit(job, "Done", err.Error())
			return
		}

		dbus.Emit(job, "RecommendedLaunchAppInfo", string(bInfos))
		dbus.Emit(job, "Done", "")
	})
	job.op.Execute()
}

type GetAllLaunchAppsJob struct {
	dbusInfo dbus.DBusInfo
	op       *operations.GetAllLaunchAppsJob

	Done          func(string)
	LaunchAppInfo func(string)
}

func (job *GetAllLaunchAppsJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

func NewGetAllLaunchAppsJob() *GetAllLaunchAppsJob {
	return &GetAllLaunchAppsJob{
		dbusInfo: genDBusInfo("GetAllLaunchAppsJob", &_GetAllLaunchAppJobCount),
		op:       operations.NewGetAllLaunchAppsJob(),
	}
}

func (job *GetAllLaunchAppsJob) Execute() {
	job.op.ListenDone(func(e error) {
		defer dbus.UnInstallObject(job)
		if e != nil {
			dbus.Emit(job, "Done", e.Error())
			return
		}

		apps := job.op.Result().([]*gio.AppInfo)
		info := LaunchAppInfo{
			Names: make([]string, len(apps)),
			Ids:   make([]string, len(apps)),
		}
		for i, app := range apps {
			info.Names[i] = app.GetName()
			info.Ids[i] = app.GetId()
			// TODO: get icon
			// app.GetIcon()
			app.Unref()
		}

		bInfos, err := json.Marshal(info)
		if err != nil {
			dbus.Emit(job, "Done", err.Error())
			return
		}

		dbus.Emit(job, "LaunchAppInfo", string(bInfos))
		dbus.Emit(job, "Done", "")
	})
	job.op.Execute()
}

// SetDefaultLaunchAppJob exports to dbus.
type SetDefaultLaunchAppJob struct {
	dbusInfo dbus.DBusInfo
	op       *operations.SetDefaultLaunchAppJob

	Done func(string)
}

// GetDBusInfo returns dbus information.
func (job *SetDefaultLaunchAppJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

// Execute SetDefaultLaunchAppJob.
func (job *SetDefaultLaunchAppJob) Execute() {
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
	_SetDefaultLaunchAppJobCount uint64
)

// NewSetDefaultLaunchAppJob creates a new SetLaunchAppJob for dbus.
func NewSetDefaultLaunchAppJob(id string, mimeType string) *SetDefaultLaunchAppJob {
	job := &SetDefaultLaunchAppJob{
		dbusInfo: genDBusInfo("SetLaunchAppJob", &_SetDefaultLaunchAppJobCount),
		op:       operations.NewSetDefaultLaunchAppJob(id, mimeType),
	}
	return job
}
