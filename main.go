package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"io"
)

const ignore = ".goFileNsync"
const lhistory_path = ignore + "/lhistory"

type FileDesc struct {
	Name    string
	ModTime time.Time
	Size    int
}

func main() {

	dir := "./test"

	files, err := lsDir(dir)

	if err != nil {
		panic(err)
	}

	var fds []FileDesc

	for _, file := range files {
		fmt.Println("Name: ", file.Name())
		fmt.Println("Mode Time: ", file.ModTime())
		fmt.Println("Size: ", file.Size())
		fds = append(fds, New(file))
	}

	mkDir()
	err = saveHistory(fds)

	if err != nil {
		panic(err)
	}

	lfds, err := loadHistory()

	if err != nil {
		panic(err)
	}

	for _, file := range lfds {
		fmt.Println("Name: ", file.Name)
		fmt.Println("Mode Time: ", file.ModTime)
		fmt.Println("Size: ", file.Size)
	}
}

func saveHistory(files []FileDesc) error {

	f, err := os.Create(lhistory_path)

	if err != nil {
		return err
	}

	for _, fd := range files {
		s := fmt.Sprintf("%s %s %d\n", fd.Name, fd.ModTime, fd.Size)
		_, err := f.WriteString(s)

		if err != nil {
			fmt.Println("Error writing file named %s data", fd.Name)
		}
	}

	return nil
}

func loadHistory() ([]FileDesc, error) {

	var files []FileDesc

	file, err := os.Open(lhistory_path)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)

	var line string

	var name string
	var modTime time.Time
	var size int

	for {
		line, err = reader.ReadString('\n')

		fmt.Sscan(line, &name, &modTime, &size)
		fd := FileDesc{
			Name:    name,
			ModTime: modTime,
			Size:    size,
		}
		files = append(files, fd)

		if err != nil {
			break
		}
	}

	if err != io.EOF {
        return nil, err
    }

	return files, nil
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
