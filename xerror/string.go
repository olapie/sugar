package xerror

type String string

func (s String) Error() string {
	return string(s)
}

const (
	NotExist String = "not exist"
)
