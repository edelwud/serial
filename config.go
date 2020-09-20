package serial

import "strconv"

type Config struct {
	BaudRate       uint32
	ByteSize       uint8
	Parity         uint8
	StopBits       uint8
	MaxReadBuffer  uint32
	MaxWriteBuffer uint32
	ReadTimeout    uint32
	WriteTimeout   uint32
}

func (config *Config) Serialize() map[string]string {
	var parity string
	switch config.Parity {
	case 0:
		parity = "NO PARITY"
		break
	case 1:
		parity = "ODD PARITY"
		break
	case 2:
		parity = "EVEN PARITY"
		break
	case 3:
		parity = "MARK PARITY"
		break
	case 4:
		parity = "SPACE PARITY"
		break
	}

	var stopBits string
	switch config.StopBits {
	case 0:
		stopBits = "1"
		break
	case 1:
		stopBits = "1.5"
		break
	case 2:
		stopBits = "2"
		break
	}

	return map[string]string{
		"Baud rate":             strconv.Itoa(int(config.BaudRate)) + " baud",
		"Byte size":             strconv.Itoa(int(config.ByteSize)) + " bit",
		"Parity":                parity,
		"Stop bits":             stopBits + " bit",
		"Max read buffer size":  strconv.Itoa(int(config.MaxReadBuffer)) + " bytes",
		"Max write buffer size": strconv.Itoa(int(config.MaxWriteBuffer)) + " bytes",
		"Timeout read":          strconv.Itoa(int(config.ReadTimeout)) + " msec",
		"Timeout write":         strconv.Itoa(int(config.WriteTimeout)) + " msec",
	}
}
