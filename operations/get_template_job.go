/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package operations

type GetTemplateJob struct {
	templateDirURI string
}

func NewGetTemplateJob(templateDirURI string) *GetTemplateJob {
	return &GetTemplateJob{
		templateDirURI: templateDirURI,
	}
}

func shouldShow(property ListProperty) bool {
	return property.CanRead && !property.IsHidden && !property.IsBackup && property.MIME != "inode/directory"
}

func (job *GetTemplateJob) Execute() []string {
	files := []string{}
	listJob := NewListDirJob(job.templateDirURI, ListJobFlagNone)
	listJob.ListenProperty(func(property ListProperty) {
		if shouldShow(property) {
			files = append(files, property.URI)
		}
	})
	listJob.Execute()

	return files
}
