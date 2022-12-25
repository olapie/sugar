package fsx

type FileInfo interface {
	UUID() string
	Name() string
	Type() string
	IsDir() bool
	Size() int64
	ModTime() int64
}
