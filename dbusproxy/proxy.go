package dbusproxy

import (
	"errors"
	"pkg.deepin.io/lib/dbus"
)

// DBusProxy is a proxy for a dbus interface.
type DBusProxy struct {
	conn              *dbus.Conn
	obj               *dbus.Object
	dest              string
	objPath           string
	iface             string
	introspectorIface string
	flags             dbus.Flags
}

// ProxyError
var (
	ErrorProxyNilConnection     = errors.New("nil dbus connection")
	ErrorProxyEmptyDestination  = errors.New("empty destination")
	ErrorProxyEmptyObjectPath   = errors.New("empty object path")
	ErrorProxyInvalidObjectPath = errors.New("invalid object path")
	ErrorProxyEmptyInterface    = errors.New("empty interface")
)

// NewDBusProxy creates a new dbus proxy.
func NewDBusProxy(conn *dbus.Conn, dest string, objPath string, iface string, flags dbus.Flags) (*DBusProxy, error) {
	if conn == nil {
		return nil, ErrorProxyNilConnection
	}
	if dest == "" {
		return nil, ErrorProxyEmptyDestination
	}
	if objPath == "" {
		return nil, ErrorProxyEmptyObjectPath
	}
	if !dbus.ObjectPath(objPath).IsValid() {
		return nil, ErrorProxyInvalidObjectPath
	}
	if iface == "" {
		return nil, ErrorProxyEmptyInterface
	}
	obj := conn.Object(dest, dbus.ObjectPath(objPath))
	return &DBusProxy{
		conn:              conn,
		obj:               obj,
		dest:              dest,
		objPath:           objPath,
		iface:             iface,
		introspectorIface: "org.freedesktop.DBus.Introspectable",
		flags:             flags,
	}, nil
}

func (proxy *DBusProxy) fullName(name string) string {
	return proxy.iface + "." + name
}

// Call calls a method of a interface.
func (proxy *DBusProxy) Call(name string, args ...interface{}) *dbus.Call {
	return proxy.obj.Call(proxy.fullName(name), proxy.flags, args...)
}

// Introspect is a function calls the Introspect method of a object path.
func (proxy *DBusProxy) Introspect() (string, error) {
	var x string
	err := proxy.obj.Call(proxy.introspectorIface+".Introspect", proxy.flags).Store(&x)
	return x, err
}

// TODO: add property and signal supports.
