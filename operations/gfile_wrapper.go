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
	"gir/gio-2.0"
	"sync"
)

// GFileWrapper wraps *gio.File to make gio.File safer and easy to use.
type GFileWrapper struct {
	*gio.File // this field should not be operated directly in common.
	lock      sync.Mutex
}

// Unref the object.
func (o *GFileWrapper) Unref() {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.unref()
}

// to avoid dead-lock problem, unref is used to do the real unref operations.
func (o *GFileWrapper) unref() {
	o.File.Unref()
	o.File = nil
}

// Reset unref the old GFile that belongs to GFileWrapper and set a new GFile to it.
func (o *GFileWrapper) Reset(x *gio.File) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.unref()
	o.File = x
}

// IsNil checks whether the GFile is nil or not.
func (o *GFileWrapper) IsNil() bool {
	return o.File == nil
}

// NewGFileWrapper creates a new GFileWrapper.
func NewGFileWrapper(file *gio.File) *GFileWrapper {
	return &GFileWrapper{
		File: file,
	}
}
