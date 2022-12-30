package xerror

import (
	"errors"
	"fmt"
)

// IsOf returns true if err can be unwrapped to T
func IsOf[T error](err error) bool {
	var zero T
	if errors.Is(err, zero) {
		return true
	}

	for {
		if _, ok := err.(T); ok {
			return true
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			break
		}
		err = u.Unwrap()
	}
	return false
}

func Not(err, target error) bool {
	return !errors.Is(err, target)
}

// Cause returns the root cause error
func Cause(err error) error {
	for {
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			break
		}
		err = u.Unwrap()
	}
	return err
}

// CauseOf returns error of type T if it's a cause
func CauseOf[T error](err error) (T, bool) {
	for {
		if v, ok := err.(T); ok {
			return v, true
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			break
		}
		err = u.Unwrap()
	}
	var zero T
	return zero, false
}

func Wrapf(err error, format string, a ...any) error {
	if err == nil {
		return nil
	}
	a = append(a, err)
	return fmt.Errorf(format+":%w", a...)
}
