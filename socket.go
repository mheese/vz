package vz

/*
#cgo darwin CFLAGS: -x objective-c -fno-objc-arc
#cgo darwin LDFLAGS: -lobjc -framework Foundation -framework Virtualization
# include "virtualization.h"
*/
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/rs/xid"
)

// SocketDeviceConfiguration for a socket device configuration.
type SocketDeviceConfiguration interface {
	NSObject

	socketDeviceConfiguration()
}

type baseSocketDeviceConfiguration struct{}

func (*baseSocketDeviceConfiguration) socketDeviceConfiguration() {}

var _ SocketDeviceConfiguration = (*VirtioSocketDeviceConfiguration)(nil)

// VirtioSocketDeviceConfiguration is a configuration of the Virtio socket device.
//
// This configuration creates a Virtio socket device for the guest which communicates with the host through the Virtio interface.
// Only one Virtio socket device can be used per virtual machine.
// see: https://developer.apple.com/documentation/virtualization/vzvirtiosocketdeviceconfiguration?language=objc
type VirtioSocketDeviceConfiguration struct {
	pointer

	*baseSocketDeviceConfiguration
}

// NewVirtioSocketDeviceConfiguration creates a new VirtioSocketDeviceConfiguration.
func NewVirtioSocketDeviceConfiguration() *VirtioSocketDeviceConfiguration {
	config := &VirtioSocketDeviceConfiguration{
		pointer: pointer{
			ptr: C.newVZVirtioSocketDeviceConfiguration(),
		},
	}
	runtime.SetFinalizer(config, func(self *VirtioSocketDeviceConfiguration) {
		self.Release()
	})
	return config
}

type VirtioSocketDevice struct {
	pointer
	dispatchQueue unsafe.Pointer
}

func newVirtioSocketDevice(ptr, dispatchQueue unsafe.Pointer) *VirtioSocketDevice {
	d := &VirtioSocketDevice{
		pointer: pointer{
			ptr: ptr,
		},
		dispatchQueue: dispatchQueue,
	}
	runtime.SetFinalizer(d, func(self *VirtioSocketDevice) {
		self.Release()
	})
	return d
}

type VirtioSocketListener struct {
	pointer
}

type ShouldAcceptNewConnectionFunc func(listener *VirtioSocketListener, conn *VirtioSocketConnection, device *VirtioSocketDevice) bool

var listeners = map[string]ShouldAcceptNewConnectionFunc{}

func NewVirtioSocketListener(f ShouldAcceptNewConnectionFunc) *VirtioSocketListener {
	id := xid.New().String()
	cs := charWithGoString(id)
	listeners[id] = f
	defer cs.Free()
	l := &VirtioSocketListener{
		pointer: pointer{
			ptr: C.newVZVirtioSocketListener(cs.CString()),
		},
	}
	runtime.SetFinalizer(l, func(self *VirtioSocketListener) {
		self.Release()
	})
	return l
}

//export listenerShouldAcceptNewConnectionFromSocketDevice
func listenerShouldAcceptNewConnectionFromSocketDevice(listener, conn, device, dispatchQueue unsafe.Pointer, listenerID *C.char) C.int {
	lID := (*char)(listenerID)
	f := listeners[lID.String()]
	if f == nil {
		return C.int(0)
	}
	l := &VirtioSocketListener{
		pointer: pointer{
			ptr: listener,
		},
	}
	runtime.SetFinalizer(l, func(self *VirtioSocketListener) {
		self.Release()
	})
	c := newVirtioSocketConnection(conn)
	d := newVirtioSocketDevice(device, dispatchQueue)
	if f(l, c, d) {
		return C.int(1)
	}
	return C.int(0)
}

func (d *VirtioSocketDevice) SetSocketListenerForPort(l *VirtioSocketListener, port uint32) {
	C.setSocketListenerForPortVZVirtioSocketDevice(d.Ptr(), d.dispatchQueue, l.Ptr(), C.uint32_t(port))
}

func (d *VirtioSocketDevice) RemoveSocketListenerForPort(port uint32) {
	C.removeSocketListenerForPortVZVirtioSocketDevice(d.Ptr(), d.dispatchQueue, C.uint32_t(port))
}

type ConnectionAttemptFunc func(conn *VirtioSocketConnection, err error)

var connections = map[string]ConnectionAttemptFunc{}

//export connectToPortForSocketDeviceHandler
func connectToPortForSocketDeviceHandler(rawFnID *C.char, vzConn, nsErr unsafe.Pointer) {
	fnID := (*char)(rawFnID)
	fn := connections[fnID.String()]
	if fn == nil {
		return
	}
	conn := newVirtioSocketConnection(vzConn)
	err := newNSError(nsErr)
	fn(conn, err)
	delete(connections, fnID.String())
}

func (d *VirtioSocketDevice) ConnectToPort(port uint32, completionHandler ConnectionAttemptFunc) {
	id := xid.New().String()
	cs := charWithGoString(id)
	connections[id] = completionHandler
	defer cs.Free()
	C.connectToPortVZVirtioSocketDevice(d.Ptr(), d.dispatchQueue, C.uint32_t(port), cs.CString())
}

type VirtioSocketConnection struct {
	pointer
}

func newVirtioSocketConnection(ptr unsafe.Pointer) *VirtioSocketConnection {
	conn := &VirtioSocketConnection{
		pointer: pointer{
			ptr: ptr,
		},
	}
	runtime.SetFinalizer(conn, func(self *VirtioSocketConnection) {
		self.Release()
	})
	return conn
}

func (c *VirtioSocketConnection) SourcePort() uint32 {
	return uint32(C.getVZVirtioSocketConnectionSourcePort(c.Ptr()))
}

func (c *VirtioSocketConnection) DestinationPort() uint32 {
	return uint32(C.getVZVirtioSocketConnectionDestinationPort(c.Ptr()))
}

func (c *VirtioSocketConnection) Fd() uintptr {
	return uintptr(C.getVZVirtioSocketConnectionFileDescriptor(c.Ptr()))
}

func (c *VirtioSocketConnection) Close() {
	C.closeVZVirtioSocketConnection(c.Ptr())
}
