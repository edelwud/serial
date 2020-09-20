package serial

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

type DCB struct {
	DCBlength uint32
	BaudRate  uint32
	Flags     uint32
	reserved  uint16
	XonLim    uint16
	XoffLim   uint16
	ByteSize  uint8
	Parity    uint8
	StopBits  uint8
	XonChar   byte
	XoffChar  byte
	ErrorChar byte
	EofChar   byte
	EvtChar   byte
	reserved1 uint16
}

var (
	procGetCommState = kernel32.NewProc("GetCommState")
	procSetCommState = kernel32.NewProc("SetCommState")
)

func (dcb *DCB) Build(handle windows.Handle, config *Config) error {
	dcb.DCBlength = uint32(unsafe.Sizeof(*dcb))
	r, _, err := procGetCommState.Call(uintptr(handle), uintptr(unsafe.Pointer(dcb)))
	if r == 0 {
		return err
	}
	dcb.BaudRate = config.BaudRate
	dcb.ByteSize = config.ByteSize
	dcb.Parity = config.Parity
	dcb.StopBits = config.StopBits

	parityBit := uint32(0)
	if config.Parity != 0 {
		parityBit = 1
	}

	dcb.Flags |= parityBit << 1    // Parity

	dcb.Flags |= 1 << 0    // Binary
	dcb.Flags |= 1 << 2    // OutxCtsFlow
	dcb.Flags |= 1 << 3    // OutxDsrFlows
	dcb.Flags |= 0x02 << 4 // DtrControl (HANDSHAKE)

	dcb.Flags |= 0 << 4 // DTR
	dcb.Flags |= 0 << 5

	dcb.Flags |= 0 << 8  // OutX
	dcb.Flags |= 0 << 9  // InX
	dcb.Flags |= 0 << 10 // ErrorChar

	dcb.Flags |= 0x02 << 12 // RtsControl (HANDSHAKE)

	dcb.Flags |= 1 << 14 // On error abort

	dcb.XonChar = 0
	dcb.XoffChar = 0xff
	dcb.XonLim = 128
	dcb.XoffLim = 128

	if r, _, err := procSetCommState.Call(uintptr(handle), uintptr(unsafe.Pointer(dcb))); r == 0 {
		return err
	}

	return nil
}
