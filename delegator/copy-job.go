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
	"gir/gio-2.0"
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/service/file-manager-backend/log"
	"pkg.deepin.io/service/file-manager-backend/operations"
)

var (
	_CopyJobCount uint64
)

// CopyJob exports to dbus.
type CopyJob struct {
	dbusInfo dbus.DBusInfo

	uris []string
	op   *operations.CopyMoveJob

	Done             func(string)
	TotalAmount      func(int64, uint16)
	ProcessedAmount  func(int64, uint16)
	ProcessedPercent func(int64)
	Copying          func(string)
	Aborted          func()
}

// GetDBusInfo returns dbus information.
func (job *CopyJob) GetDBusInfo() dbus.DBusInfo {
	return job.dbusInfo
}

// Abort the job.
func (job *CopyJob) Abort() {
	job.op.Abort()
	dbus.Emit(job, "Aborted")
}

func (job *CopyJob) listenSignals() {
	job.op.ListenTotalAmount(func(amount int64, unit operations.AmountUnit) {
		dbus.Emit(job, "TotalAmount", amount, uint16(unit))
	})
	job.op.ListenProcessedAmount(func(amount int64, unit operations.AmountUnit) {
		dbus.Emit(job, "ProcessedAmount", amount, uint16(unit))
	})
	job.op.ListenPercent(func(percent int64) {
		dbus.Emit(job, "ProcessedPercent", percent)
	})
	job.op.ListenCopying(func(srcURL string) {
		Log.Debug("copying", srcURL)
		dbus.Emit(job, "Copying", srcURL)
	})
	job.op.ListenCreatingDir(func(dirURL string) {
		// TODO
		// dbus.Emit(job, "CreatingDir", dirURL)
	})
	job.op.ListenCopyingMovingDone(func(srcURL string, destURL string) {
	})
	job.op.ListenDone(func(err error) {
		defer dbus.UnInstallObject(job)
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		dbus.Emit(job, "Done", errMsg)
	})
	job.op.ListenAborted(func() {
		defer dbus.UnInstallObject(job)
		dbus.Emit(job, "Aborted")
	})
}

// Execute copy job.
func (job *CopyJob) Execute() {
	job.listenSignals()
	// TODO
	// operations.FileUndoManagerInstance().RecordJob(COPY, job.op)
	job.op.Execute()
}

// NewCopyJob creates a new copy job for dbus.
func NewCopyJob(srcUrls []string, destDirURL string, targetName string, flags uint32, uiDelegate IUIDelegate) *CopyJob {
	job := &CopyJob{
		dbusInfo: genDBusInfo("CopyJob", &_CopyJobCount),
		op:       operations.NewCopyJob(srcUrls, destDirURL, targetName, gio.FileCopyFlags(flags), uiDelegate),
	}
	return job
}
