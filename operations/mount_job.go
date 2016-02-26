/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package operations

type IMountUI interface {
	Username() string
	Domain() string
	Password() string
	IsAnonymous() bool
	RememberPasswordFlags() int

	// TODO: ???
	AskPassword()
	AskQuestion()
	ShowProcess()
	ShowUnmountProcess()
}

type MountJob struct {
	*CommonJob
}

func newMountJob() *MountJob {
	job := &MountJob{
		CommonJob: newCommon(nil), // TODO
	}
	return job
}

func NewMountJob() *MountJob {
	return newMountJob()
}
