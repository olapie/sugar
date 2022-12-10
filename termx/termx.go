package termx

import (
	"os"

	"golang.org/x/term"
)

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
