package xerror

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type unwrapError interface {
	Unwrap() error
}

type unwrapErrors interface {
	Unwrap() []error
}

// CauseOf returns true if err can be unwrapped to T
func CauseOf[T error](err error) (T, bool) {
	var zero T
	for {
		if v, ok := err.(T); ok {
			return v, true
		}

		if u, ok := err.(unwrapError); ok {
			err = u.Unwrap()
			continue
		}

		if u, ok := err.(unwrapErrors); ok {
			for _, e := range u.Unwrap() {
				if v, ok := CauseOf[T](e); ok {
					return v, true
				}
			}
			return zero, false
		}
	}
	return zero, false
}

func Not(err, target error) bool {
	return !errors.Is(err, target)
}

func Wrapf(err error, format string, a ...any) error {
	if err == nil {
		return nil
	}
	a = append(a, err)
	return fmt.Errorf(format+":%w", a...)
}

func IsNotExist(err error) bool {
	return errors.Is(err, NotExist) ||
		errors.Is(err, sql.ErrNoRows) ||
		errors.Is(err, os.ErrNotExist) ||
		GetCode(err) == http.StatusNotFound
}

func IsNilOrNotExist(err error) bool {
	return err == nil || err.Error() == "nil" || IsNotExist(err)
}

func Or(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func And(errs ...error) error {
	for _, err := range errs {
		if err == nil {
			return nil
		}
	}
	return errors.Join(errs...)
}

func Combine(errs ...error) error {
	for i := len(errs) - 1; i >= 0; i-- {
		if errs[i] == nil {
			errs = append(errs[:i], errs[i+1:]...)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
}

func OrFn(err error, errFns ...func() error) error {
	if err != nil {
		return err
	}

	for _, fn := range errFns {
		if e := fn(); e != nil {
			return e
		}
	}

	return nil
}
