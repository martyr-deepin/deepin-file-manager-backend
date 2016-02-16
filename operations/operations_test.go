/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package operations_test

import (
	. "pkg.deepin.io/service/file-manager-backend/operations"
	"net/url"
	"path/filepath"
)

var testdataDir, _ = filepath.Abs("./testdata")

type UIMock struct {
	skip bool
}

func (*UIMock) AskRetry(primaryText string, secondaryText string, detailText string) Response {
	return NewResponse(ResponseSkip, false)
}
func (*UIMock) AskDeleteConfirmation(primaryText string, secondaryText string, detailText string) bool {
	return true
}

func (*UIMock) AskDelete(primaryText string, secondaryText string, detailText string, flags UIFlags) Response {
	return NewResponse(ResponseSkip, true)
}

func (mock *UIMock) ConflictDialog() Response {
	code := ResponseAutoRename
	if mock.skip {
		code = ResponseSkip
	}
	return NewResponse(code, true)
}

func (*UIMock) AskSkip(primaryText string, secondaryText string, detailText string, flags UIFlags) Response {
	// TODO:
	return NewResponse(ResponseSkip, true)
}

func NewUIMock(skip bool) *UIMock {
	return &UIMock{
		skip: skip,
	}
}

var skipMock = NewUIMock(true)
var renameMock = NewUIMock(false)

func pathToURL(path string) (*url.URL, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	return url.Parse(absPath)
}
