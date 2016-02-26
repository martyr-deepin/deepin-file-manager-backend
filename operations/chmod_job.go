/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package operations

import (
	"errors"
	"gir/gio-2.0"
	"os"
)

// ChmodJob change the mode of a file/directory.
type ChmodJob struct {
	*CommonJob
	file       *gio.File
	permission uint32
}

func (job *ChmodJob) finalize() {
	if job.file != nil {
		job.file.Unref()
	}
	if job.CommonJob != nil {
		job.CommonJob.finalize()
	}
}

// Execute the ChmodJob.
func (job *ChmodJob) Execute() {
	defer finishJob(job)

	if job.file != nil {
		job.setError(os.Chmod(job.file.GetPath(), os.FileMode(job.permission)))
		return
	}
	job.setError(errors.New("no such a file"))
}

func newChmodJob(file *gio.File, permission uint32) *ChmodJob {
	job := &ChmodJob{
		CommonJob:  newCommon(nil),
		file:       file,
		permission: permission,
	}
	return job
}

// NewChmodJob creates a new ChmodJob.
func NewChmodJob(uri string, permission uint32) *ChmodJob {
	file := gio.FileNewForCommandlineArg(uri)
	return newChmodJob(file, permission)
}
