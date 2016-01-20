package dbusproxy

import (
	"errors"
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"reflect"
	"runtime"
	"sync"
)

func SignalName(iface, name string) string {
	return iface + "." + name
}

type Signal struct {
	Chan <-chan *dbus.Signal
	Name string
}

// DBusProxy is a proxy for a dbus interface.
type DBusProxy struct {
	conn              *dbus.Conn
	obj               *dbus.Object
	dest              string
	objPath           string
	iface             string
	introspectorIface string
	flags             dbus.Flags
	sigChanMap        map[string][]Signal
	lock              sync.Mutex
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
	proxy := &DBusProxy{
		conn:              conn,
		obj:               obj,
		dest:              dest,
		objPath:           objPath,
		iface:             iface,
		introspectorIface: "org.freedesktop.DBus.Introspectable",
		flags:             flags,
		sigChanMap:        map[string][]Signal{},
	}
	runtime.SetFinalizer(proxy, func(proxy *DBusProxy) {
		proxy.finalize()
	})
	return proxy, nil
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

func (proxy *DBusProxy) createSignalChan(sigName string) Signal {
	proxy.lock.Lock()
	defer proxy.lock.Unlock()
	sigChan := proxy.conn.Signal()
	sig := Signal{
		Chan: sigChan,
		Name: sigName,
	}
	sigs := proxy.sigChanMap[sigName]
	proxy.sigChanMap[sigName] = append(sigs, sig)

	return sig
}

func (proxy *DBusProxy) removeSignalChan(removingSig Signal) {
	proxy.lock.Lock()
	defer proxy.lock.Unlock()

	for _, sig := range proxy.sigChanMap[removingSig.Name] {
		if sig.Chan == removingSig.Chan {
			proxy.conn.DetachSignal(sig.Chan)
			break
		}
	}
}

func (proxy *DBusProxy) deleteSignalChan(sigName string) {
	for _, sig := range proxy.sigChanMap[sigName] {
		proxy.conn.DetachSignal(sig.Chan)
	}
	rule := proxy.buildRule(sigName)
	proxy.conn.BusObject().Call("org.freedesktop.DBus.RemoveMatch", 0, rule)
	delete(proxy.sigChanMap, sigName)
}

func (proxy *DBusProxy) finalize() {
	proxy.lock.Lock()
	defer proxy.lock.Unlock()
	for sigName := range proxy.sigChanMap {
		proxy.deleteSignalChan(sigName)
	}
	runtime.SetFinalizer(proxy, nil)
}

func (proxy *DBusProxy) buildRule(sigName string) string {
	rule := fmt.Sprintf("type='signal',sender='%s',path='%s',interface='%s',member='%s'", proxy.dest, proxy.objPath, proxy.iface, sigName)
	return rule
}

func (proxy *DBusProxy) Subscribe(sigName string, f interface{}) func() {
	sig := proxy.createSignalChan(sigName)
	rule := proxy.buildRule(sigName)
	proxy.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, rule)
	go func() {
		fn := reflect.ValueOf(f)
		for v := range sig.Chan {
			if v.Name != SignalName(proxy.iface, sigName) || v.Path != dbus.ObjectPath(proxy.objPath) {
				continue
			}
			l := len(v.Body)
			args := make([]reflect.Value, l)
			for i, v := range v.Body {
				args[i] = reflect.ValueOf(v)
			}
			fn.Call(args)
		}
	}()
	return func() {
		proxy.removeSignalChan(sig)
	}
}

func (proxy *DBusProxy) Unsubscribe(sigName string) {
	proxy.lock.Lock()
	defer proxy.lock.Unlock()
	proxy.deleteSignalChan(sigName)
}

// TODO: handle property
