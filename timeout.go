package serial

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

type CommTimeouts struct {
	ReadIntervalTimeout         uint32
	ReadTotalTimeoutMultiplier  uint32
	ReadTotalTimeoutConstant    uint32
	WriteTotalTimeoutMultiplier uint32
	WriteTotalTimeoutConstant   uint32
}

var (
	procSetCommTimeouts = kernel32.NewProc("SetCommTimeouts")
)

func (timeouts *CommTimeouts) Configure(handle windows.Handle, readTimeout uint32, writeTimeout uint32) error {
	timeouts.ReadIntervalTimeout = 2<<31 - 1
	timeouts.ReadTotalTimeoutConstant = readTimeout
	timeouts.ReadTotalTimeoutMultiplier = 0
	timeouts.WriteTotalTimeoutConstant = writeTimeout
	timeouts.WriteTotalTimeoutMultiplier = 0
	if r, _, err := procSetCommTimeouts.Call(uintptr(handle), uintptr(unsafe.Pointer(timeouts))); r == 0 {
		return err
	}
	return nil
}
