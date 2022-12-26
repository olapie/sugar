package nomobile

type FileEntry interface {
	GetID() string
	Name() string
	IsDir() bool
	Size() int64
	ModTime() int64
	MIMEType() string
	SubIDs() []string
}
