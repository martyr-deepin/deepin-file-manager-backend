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
	// "time"
)

var (
	_ListJobCount uint64
)

// ListJob exports to dbus.
type ListJob struct {
	dbusInfo dbus.DBusInfo
	op       *operations.ListJob

	EntryInfo       func(operations.ListProperty)
	Done            func(string)
	ProcessedAmount func(int64, uint16)
	Aborted         func()
}

// GetDBusInfo returns dbus information.
func (job *ListJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

// NewListJob creates a new list job for dbus.
func NewListJob(path string, flags operations.ListJobFlag) *ListJob {
	job := &ListJob{
		dbusInfo: genDBusInfo("ListJob", &_ListJobCount),
		op:       operations.NewListDirJob(path, flags),
	}
	return job
}

// Execute list job.
func (job *ListJob) Execute() []operations.ListProperty {
	defer dbus.UnInstallObject(job)
	var files []operations.ListProperty
	job.op.ListenProcessedAmount(func(size int64, unit operations.AmountUnit) {
		dbus.Emit(job, "ProcessedAmount", size, uint16(unit))
	})
	job.op.ListenProperty(func(property operations.ListProperty) {
		dbus.Emit(job, "EntryInfo", property)
		files = append(files, property)
	})
	job.op.ListenDone(func(err error) {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		// time.Sleep(time.Microsecond * 200)
		dbus.Emit(job, "Done", errMsg)
	})
	job.op.ListenAborted(func() {
		defer dbus.UnInstallObject(job)
		dbus.Emit(job, "Aborted")
	})
	job.op.Execute()

	return files
}

// Abort the job.
func (job *ListJob) Abort() {
	job.op.Abort()
	dbus.Emit(job, "Aborted")
	dbus.UnInstallObject(job)
}
