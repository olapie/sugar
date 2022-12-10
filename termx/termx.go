package termx

import (
	"bytes"
	"os"

	"golang.org/x/term"
)

const (
	ControlC = 0x03
	ControlD = 0x04
	ControlZ = 0x1A
)

func Read(stop byte) ([]byte, error) {
	var buf bytes.Buffer
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}

	b := make([]byte, 1)
	for {
		_, err = os.Stdin.Read(b)
		if err != nil {
			term.Restore(int(os.Stdin.Fd()), oldState)
			return nil, err
		}
		switch b[0] {
		case stop, ControlC, ControlZ:
			break
		default:
			buf.WriteByte(b[0])
		}
	}
	term.Restore(int(os.Stdin.Fd()), oldState)
	return buf.Bytes(), nil
}

func ReadOne() (byte, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return 0, err
	}

	var b [1]byte
	_, err = os.Stdin.Read(b[:])
	term.Restore(int(os.Stdin.Fd()), oldState)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}
