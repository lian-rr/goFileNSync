package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const ignore = ".goFileNsync"

type FileDesc struct {
	Name    string
	ModTime time.Time
	Size    int
}

func main() {

	dir := "."

	files, err := lsDir(dir)

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fd := New(file)

		fmt.Println("Name: ", fd.Name)
		fmt.Println("Mode Time: ", fd.ModTime)
		fmt.Println("Size: ", fd.Size)
	}

}

func New(f os.FileInfo) FileDesc {
	return FileDesc{
		Name:    f.Name(),
		ModTime: f.ModTime(),
		Size:    int(f.Size()),
	}
}

func mkDir() {
	os.Mkdir(ignore, os.ModePerm)
}

func lsDir(p string) ([]os.FileInfo, error) {
	var files []os.FileInfo
	err := filepath.Walk(p, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() && f.Name() == ignore {
			return filepath.SkipDir
		}

		if !f.IsDir() {
			files = append(files, f)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error listing the directory")
	}

	return files, nil
}
