package mobilex

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type DirInfo struct {
	Document  string
	Cache     string
	Temporary string
}

func (d *DirInfo) MustMakeDirs() {
	if d.Document != "" {
		MustMkdir(d.Document)
	} else {
		fmt.Println("Document directory is not specified")
	}

	if d.Cache != "" {
		MustMkdir(d.Cache)
	} else {
		fmt.Println("Cache directory is not specified")
	}

	if d.Temporary != "" {
		MustMkdir(d.Temporary)
	} else {
		fmt.Println("Temporary directory is not specified")
	}
}

func NewDirInfo() *DirInfo {
	return new(DirInfo)
}

func NewTestDirInfo() *DirInfo {
	return &DirInfo{
		Document:  "testdata/document",
		Cache:     "testdata/cache",
		Temporary: "testdata/temporary",
	}
}

func GetDiskSize(path string) *Int64E {
	res := new(Int64E)
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			res.Val += info.Size()
		}
		return err
	})
	if err != nil {
		res.Err = ToError(err)
	}
	return res
}

func MustMkdir(dir string) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatalf("Make dir: %s, %v\n", dir, err)
	}
}
