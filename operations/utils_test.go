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
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "pkg.deepin.io/service/file-manager-backend/operations"
)

func TestShortenUtf8String(t *testing.T) {
	Convey("shorten english", t, func() {
		a := "abc"
		b := ShortenUtf8String(a, 0)
		So(b, ShouldEqual, "abc")

		c := ShortenUtf8String(a, 1)
		So(c, ShouldEqual, "ab")
	})

	Convey("shorten chiness", t, func() {
		a := "ab我卡"

		b := ShortenUtf8String(a, 0)
		So(b, ShouldEqual, "ab我卡")

		c := ShortenUtf8String(a, 1)
		So(c, ShouldEqual, "ab我")

		d := ShortenUtf8String(a, 2)
		So(d, ShouldEqual, "ab")

		e := ShortenUtf8String(a, 3)
		So(e, ShouldEqual, "a")
	})

	Convey("the reduce num is out of range", t, func() {
		a := "aaa"
		b := ShortenUtf8String(a, len(a)+1)
		So(b, ShouldEqual, "")
	})

	Convey("the reduce num is a negative", t, func() {
		a := "xxx"
		b := ShortenUtf8String(a, -1)
		So(b, ShouldEqual, a)
	})
}

func TestParsePreviousDuplicateName(t *testing.T) {
	oldLang := os.Getenv("LANGUAGE")
	os.Setenv("LANGUAGE", "en_US")
	defer os.Setenv("LANGUAGE", oldLang)

	Convey("parse previous duplicate name without copy", t, func() {
		namebase, suffix, count := ParsePreviousDuplicateName("test.png")
		So(namebase, ShouldEqual, "test")
		So(suffix, ShouldEqual, ".png")
		So(count, ShouldEqual, 0)
	})

	Convey("parse previous duplicate name with copy2", t, func() {
		namebase, suffix, count := ParsePreviousDuplicateName("test (Copy).png")
		So(namebase, ShouldEqual, "test")
		So(suffix, ShouldEqual, ".png")
		So(count, ShouldEqual, 1)
	})

	Convey("parse previous duplicate name with copy 2", t, func() {
		namebase, suffix, count := ParsePreviousDuplicateName("test (Copy 2).png")
		So(namebase, ShouldEqual, "test")
		So(suffix, ShouldEqual, ".png")
		So(count, ShouldEqual, 2)
	})
}

func TestMakeNextDuplicateName(t *testing.T) {
	oldLang := os.Getenv("LANGUAGE")
	os.Setenv("LANGUAGE", "en_US")
	defer os.Setenv("LANGUAGE", oldLang)

	Convey("make next duplicate name with 0 count", t, func() {
		So(MakeNextDuplicateName("test", ".png", 0, 0), ShouldEqual, "test (Copy).png")
	})

	Convey("make next duplicate name with 1 count", t, func() {
		So(MakeNextDuplicateName("test", ".png", 1, 0), ShouldEqual, "test (Copy).png")
	})

	Convey("make next duplicate name with 2 count", t, func() {
		So(MakeNextDuplicateName("test", ".png", 2, 0), ShouldEqual, "test (Copy 2).png")
	})

	Convey("make next duplicate name with over length", t, func() {
		So(MakeNextDuplicateName("looooooooooooooooooongtest", ".png", 1, 20), ShouldEqual, "looooooooooooooooooo (Copy).png")
	})
}
