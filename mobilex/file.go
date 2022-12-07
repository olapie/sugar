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

func GetDiskSize(path string) int64 {
	var sum int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			sum += info.Size()
		}
		return err
	})
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return sum
}

func MustMkdir(dir string) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatalf("Make dir: %s, %v\n", dir, err)
	}
}
