package dbusproxy_test

import (
	. "pkg.deepin.io/service/file-manager-backend/dbusproxy"
	. "github.com/smartystreets/goconvey/convey"
	"pkg.deepin.io/lib/dbus"
	"testing"
)

func TestProxy(t *testing.T) {
	Convey("DBusProxy", t, func() {
		Convey("should handle nil connection", func() {
			proxy, err := NewDBusProxy(nil,
				"org.freedesktop.DBus",
				"/",
				"org.freedesktop.DBus",
				0,
			)
			So(proxy, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, ErrorProxyNilConnection)
		})

		Convey("should handle empty destination", func() {
			conn, err := dbus.SessionBus()
			So(err, ShouldBeNil)
			proxy, err := NewDBusProxy(conn,
				"",
				"/",
				"org.freedesktop.DBus",
				0,
			)
			So(proxy, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, ErrorProxyEmptyDestination)
		})

		Convey("should handle empty object path", func() {
			conn, err := dbus.SessionBus()
			So(err, ShouldBeNil)

			proxy, err := NewDBusProxy(conn,
				"org.freedesktop.DBus",
				"",
				"org.freedesktop.DBus",
				0,
			)
			So(proxy, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, ErrorProxyEmptyObjectPath)
		})

		Convey("should handle empty interface path", func() {
			conn, err := dbus.SessionBus()
			So(err, ShouldBeNil)

			proxy, err := NewDBusProxy(conn,
				"org.freedesktop.DBus",
				"/",
				"",
				0,
			)
			So(proxy, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, ErrorProxyEmptyInterface)
		})

		Convey("should handle invalid object path", func() {
			conn, err := dbus.SessionBus()
			So(err, ShouldBeNil)

			proxy, err := NewDBusProxy(conn,
				"org.freedesktop.DBus",
				"1",
				"org.freedesktop.DBus",
				0,
			)
			So(proxy, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, ErrorProxyInvalidObjectPath)
		})

		Convey("call method", func() {
			conn, err := dbus.SessionBus()
			So(err, ShouldBeNil)
			proxy, err := NewDBusProxy(conn,
				"org.freedesktop.DBus",
				"/",
				"org.freedesktop.DBus",
				0,
			)

			So(err, ShouldBeNil)
			var x string
			err = proxy.Call("GetId").Store(&x)
			So(err, ShouldBeNil)
			So(x, ShouldNotEqual, "")
		})

		Convey("call a not existed method", func() {
			conn, err := dbus.SessionBus()
			So(err, ShouldBeNil)
			proxy, err := NewDBusProxy(conn,
				"org.freedesktop.DBus",
				"/",
				"org.freedesktop.DBus",
				0,
			)

			So(err, ShouldBeNil)
			err = proxy.Call("GetXXX").Store()
			So(err, ShouldNotBeNil)
		})
	})
}
