package nomobile

import (
	"errors"

	"code.olapie.com/sugar/v2/xerror"
)

type Result[T any] struct {
	val T
	err *xerror.Error
}

func (r *Result[T]) Value() T {
	return r.val
}

func (r *Result[T]) ErrorCode() int {
	if r.err == nil {
		return 0
	}
	return r.err.Code()
}

func (r *Result[T]) ErrorMessage() string {
	if r.err == nil {
		return ""
	}
	return r.err.Message()
}

func (r *Result[T]) SetValue(v T) {
	r.val = v
}

func (r *Result[T]) SetErrorCode(code int) {
	if r.err == nil {
		r.err = xerror.New(code, "")
	} else {
		r.err = xerror.New(code, r.err.Message())
	}
}

func (r *Result[T]) SetMessage(message string) {
	if r.err == nil {
		r.err = xerror.New(0, message)
	} else {
		r.err = xerror.New(r.err.Code(), message)
	}
}

func (r *Result[T]) SetError(err error) {
	res := ErrorResult[T](err)
	*r = *res
}

func ValueResult[T any](v T) *Result[T] {
	return &Result[T]{
		val: v,
	}
}

func ErrorResult[T any](err error) *Result[T] {
	res := new(Result[T])
	if err == nil {
		return res
	}

	if errors.As(err, &res.err) {
		return res
	}

	res.err = xerror.New(xerror.GetCode(err), err.Error())
	return res
}
