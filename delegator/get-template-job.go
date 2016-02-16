/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
