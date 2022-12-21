package termx

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"code.olapie.com/sugar/must"
	"golang.org/x/crypto/ssh/terminal"
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

func ReadPassword(msg ...any) string {
	var pass []byte
	for len(pass) == 0 {
		fmt.Print(msg...)
		fmt.Print(": ")
		pass = must.Get(terminal.ReadPassword(syscall.Stdin))
		fmt.Println()
	}
	return string(pass)
}

func ReadConfirmedPassword(prompt1, prompt2 string) *string {
	for i := 0; i < 3; i++ {
		pass1 := ReadPassword(prompt1)
		pass2 := ReadPassword(prompt2)
		if pass1 == pass2 {
			return &pass1
		}
	}
	return nil
}

func ConfirmInput(answer string) bool {
	answer = strings.TrimSpace(answer)
	if answer == "" {
		panic("answer cannot be empty")
	}
	fmt.Printf("Enter '%s' to confirm: ", answer)
	var actual string
	fmt.Scanln(&actual)
	return actual == answer
}
