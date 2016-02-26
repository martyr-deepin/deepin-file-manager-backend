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

var (
	_RenameJobCount uint64 = 0
)

type RenameJob struct {
	dbusInfo dbus.DBusInfo
	op       *operations.RenameJob

	Done func(string)
	// TODO: it maybe be better that using 'Renamed func(newFile string, oldName string)' as signal.
	OldName func(string)
	NewFile func(string)
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
