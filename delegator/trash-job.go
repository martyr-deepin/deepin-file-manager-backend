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
	"sync"
)

var (
	_TrashJobCount     uint64
	_TrashJobCountLock sync.Mutex
)

// TrashJob exports to dbus.
type TrashJob struct {
	dbusInfo dbus.DBusInfo
	uris     []string
	op       *operations.DeleteJob

	Done             func()
	Trashing         func(string)
	Deleting         func(string)
	ProcessedAmount  func(int64, uint16)
	ProcessedPercent func(int64)
	Aborted          func()
}

// GetDBusInfo returns dbus information.
func (job *TrashJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

// NewTrashJob creates a new trash job for dbus.
func NewTrashJob(urls []string, shouldConfirmTrash bool, uiDelegate IUIDelegate) *TrashJob {
	return &TrashJob{
		dbusInfo: genDBusInfo("TrashJob", &_TrashJobCount),
		op:       operations.NewTrashJob(urls, shouldConfirmTrash, uiDelegate),
	}
}

// Execute trash job.
func (job *TrashJob) Execute() {
	job.op.ListenProcessedAmount(func(size int64, unit operations.AmountUnit) {
		dbus.Emit(job, "ProcessedAmount", size, uint16(unit))
	})
	job.op.ListenPercent(func(percent int64) {
		dbus.Emit(job, "ProcessedPercent", percent)
	})
	job.op.ListenAborted(func() {
		defer dbus.UnInstallObject(job)
		dbus.Emit(job, "Aborted")
	})
	job.op.ListenDone(func(err error) {
		dbus.Emit(job, "Done")
	})
	// TODO: fill signals.
	defer dbus.UnInstallObject(job)
	job.op.Execute()
}

// Abort the job.
func (job *TrashJob) Abort() {
	job.op.Abort()
}
