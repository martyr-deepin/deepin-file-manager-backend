package main_test

import (
	. "deepin-file-manager/dbusproxy"
	. "github.com/smartystreets/goconvey/convey"
	"os/exec"
	"path/filepath"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"testing"
)

var _Dest = "com.deepin.filemanager.Backend.Operations"
var _Path = "/com/deepin/filemanager/Backend/Operations"
var _Iface = _Dest

func TestNewDeleteJob(t *testing.T) {
	conn, err := dbus.SessionBus()
	if err != nil {
		t.Error("get session bus failed")
		t.Fail()
	}

	backendProxy, err := NewDBusProxy(conn, _Dest, _Path, _Iface, 0)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	x := func(t *testing.T, fn func() (string, error)) func() {
		return func() {
			target, err := fn()
			if err != nil {
				t.Error(err)
			}

			f := gio.FileNewForCommandlineArg(target)
			So(f.QueryExists(nil), ShouldBeTrue)
			f.Unref()

			var objPath dbus.ObjectPath
			var iface string
			err = backendProxy.Call("NewDeleteJob", []string{target}, "", "", "").Store(&objPath, &iface)
			if err != nil {
				t.Error(err)
			}

			deleteProxy, err := NewDBusProxy(conn, _Dest, string(objPath), iface, 0)
			err = deleteProxy.Call("Execute").Store()
			if err != nil {
				t.Error(err.Error())
			}

			f = gio.FileNewForCommandlineArg(target)
			// TODO: this test result is not stable, add a timer to check.
			So(f.QueryExists(nil), ShouldBeFalse)
			f.Unref()
		}
	}

	Convey("should remove a file correctly", t, x(t, func() (string, error) {
		target, err := filepath.Abs("./testdata/todelete.txt")
		if err != nil {
			return "", err
		}

		err = exec.Command("touch", target).Run()
		return target, err
	}))

	Convey("should delete a empty directory correctly", t, x(t, func() (string, error) {
		target, err := filepath.Abs("./testdata/todelete.dir")
		if err != nil {
			return target, err
		}
		err = exec.Command("mkdir", target).Run()
		return target, err
	}))

	Convey("should delete a non-empty directory correctly", t, x(t, func() (string, error) {
		target, err := filepath.Abs("./testdata/todelete.dir")
		if err != nil {
			return target, err
		}
		err = exec.Command("mkdir", target).Run()
		err = exec.Command("touch", target+"/a").Run()
		return target, err
	}))
}

func TestStatJob(t *testing.T) {
	conn, err := dbus.SessionBus()
	if err != nil {
		t.Error("get session bus failed")
		t.Fail()
	}

	backendProxy, err := NewDBusProxy(conn, _Dest, _Path, _Iface, 0)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	var objPath dbus.ObjectPath
	var iface string
	err = backendProxy.Call("NewStatJob", "/tmp").Store(&objPath, &iface)
	if err != nil {
		t.Fail()
	}

	statProxy, err := NewDBusProxy(conn, _Dest, string(objPath), iface, 0)
	if err != nil {
		t.Fail()
	}

	err = statProxy.Call("Execute").Store()
	if err != nil {
		t.Fail()
	}
}
