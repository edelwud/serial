package serial

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

type Serial interface {
	Open(string) error
	Close() error
	Write([]byte) error
	Read([]byte) (uint32, error)
	GetConfig() Config
}

type SerialPort struct {
	Name   string
	DCB    *DCB
	Handle windows.Handle
	Config *Config
}

var (
	kernel32                = windows.NewLazyDLL("kernel32.dll")
	procSetupComm           = kernel32.NewProc("SetupComm")
	procGetOverlappedResult = kernel32.NewProc("GetOverlappedResult")
	procPurgeComm           = kernel32.NewProc("PurgeComm")
)

func (port *SerialPort) Open(com string) error {
	var err error
	port.Handle, err = windows.CreateFile(
		windows.StringToUTF16Ptr("\\\\.\\"+com),
		windows.GENERIC_WRITE|windows.GENERIC_READ,
		0,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL|windows.FILE_FLAG_OVERLAPPED,
		windows.InvalidHandle)
	if err != nil {
		return err
	}

	return nil
}

func (port *SerialPort) GetConfig() Config {
	return *port.Config
}

func (port *SerialPort) Close() error {
	if err := windows.CloseHandle(port.Handle); err != nil {
		return err
	}
	return nil
}

func (port *SerialPort) Clear() error {
	var flags uint32 = 0x0001 | 0x0002 | 0x004 | 0x0008
	if r, _, err := procPurgeComm.Call(uintptr(port.Handle), uintptr(flags)); r == 0 {
		return err
	}
	return nil
}

func (port *SerialPort) Write(buffer []byte) error {
	var written uint32
	var overlapped windows.Overlapped

	if err := port.Clear(); err != nil {
		return err
	}

	if err := windows.WriteFile(
		port.Handle,
		buffer,
		&written,
		&overlapped,
	); err != windows.ERROR_IO_PENDING {
		return err
	}
	return nil
}

func (port *SerialPort) Read(buffer []byte) (uint32, error) {
	var overlapped windows.Overlapped
	var err error
	overlapped.HEvent, err = windows.CreateEvent(nil, 1, 0, nil)
	if err != nil {
		return 0, err
	}

	var read uint32
	err = windows.ReadFile(port.Handle, buffer, &read, &overlapped)
	if err == nil {
		return 0, nil
	}
	if err != windows.ERROR_IO_PENDING {
		return 0, err
	}
	r, _, err := procGetOverlappedResult.Call(uintptr(port.Handle),
		uintptr(unsafe.Pointer(&overlapped)),
		uintptr(unsafe.Pointer(&read)), 1)
	if read == 0 {
		return 0, err
	}
	if r == 0 {
		return 0, err
	}

	if err := windows.CloseHandle(overlapped.HEvent); err != nil {
		return 0, err
	}

	return read, nil
}

func Open(com string, config *Config) (Serial, error) {
	serial := &SerialPort{}
	serial.Config = config
	if err := serial.Open(com); err != nil {
		return nil, err
	}

	if r, _, err := procSetupComm.Call(
		uintptr(serial.Handle),
		uintptr(config.MaxReadBuffer),
		uintptr(config.MaxWriteBuffer),
	); r == 0 {
		return nil, err
	}

	serial.DCB = &DCB{}
	if err := serial.DCB.Build(serial.Handle, config); err != nil {
		return nil, err
	}

	timeouts := &CommTimeouts{}
	if err := timeouts.Configure(serial.Handle, config.ReadTimeout, config.WriteTimeout); err != nil {
		return nil, err
	}

	return serial, nil
}
