package vz

/*
#cgo darwin CFLAGS: -x objective-c -fno-objc-arc
#cgo darwin LDFLAGS: -lobjc -framework Foundation -framework Virtualization
# include "virtualization.h"
*/
import "C"

// DirectorySharingDeviceConfiguration for a directory sharing device configuration.
type DirectorySharingDeviceConfiguration interface {
	NSObject

	directorySharingDeviceConfiguration()
}

type baseDirectorySharingDeviceConfiguration struct{}

func (*baseDirectorySharingDeviceConfiguration) directorySharingDeviceConfiguration() {}

type VirtioFileSystemDeviceConfiguration struct {
	pointer

	*baseDirectorySharingDeviceConfiguration
}

/*
func NewVirtioFileSystemDeviceConfiguration(tag string) *VirtioFileSystemDeviceConfiguration {
	tagChars := charWithGoString(tag)
	defer tagChars.Free()

	d := &VirtioFileSystemDeviceConfiguration{
		pointer: pointer{
			ptr: C.newVZVirtioFileSystemDeviceConfiguration(tagChars.CString()),
		},
	}
	runtime.SetFinalizer(d, func(self *VirtioFileSystemDeviceConfiguration) {
		self.Release()
	})
	return d
}
*/
