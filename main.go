package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const rootDir = ".goFileNsync"
const localHistory = rootDir + "/localHistory"

type FileDesc struct {
	Name    string
	ModTime time.Time
	Size    int
}

func main() {

	writingRouting()

	fds := loadDescFile(localHistory)

	for _, fd := range fds {
		fmt.Printf("%s %d %s\n", fd.Name, fd.Size, fd.ModTime.String())
	}

}

func writingRouting() {
	dir := "test/"
	files, err := lsDir(dir)
	check(err)

	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		mkDir()
	}

	fds := toFileDesc(files)

	fdsBytes := fdToBytes(fds)

	saveDescFile(localHistory, fdsBytes)
}

func New(f os.FileInfo) FileDesc {
	return FileDesc{
		Name:    f.Name(),
		ModTime: f.ModTime(),
		Size:    int(f.Size()),
	}
}

func (fd FileDesc) toBytes() []byte {
	return []byte(fmt.Sprintf("%s %d %d\n", fd.Name, fd.Size, fd.ModTime.UnixNano()))
}

func stringToFileDesc(data string) FileDesc {
	var name string
	var size int
	var modTime int64

	_, err := fmt.Sscanf(data, "%s %d %d\n", &name, &size, &modTime)

	check(err)

	return FileDesc{
		Name:    name,
		Size:    size,
		ModTime: time.Unix(0, modTime),
	}
}

func toFileDesc(files []os.FileInfo) []FileDesc {
	var fds []FileDesc
	for _, file := range files {
		fds = append(fds, New(file))
	}

	return fds
}

func fdToBytes(fds []FileDesc) [][]byte {
	var fdsBytes [][]byte

	for _, fd := range fds {
		fdsBytes = append(fdsBytes, fd.toBytes())
	}

	return fdsBytes
}

func mkDir() {
	os.Mkdir(rootDir, os.ModePerm)
}

func lsDir(p string) ([]os.FileInfo, error) {
	var files []os.FileInfo
	err := filepath.Walk(p, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() && f.Name() == rootDir {
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

func saveDescFile(dir string, data [][]byte) {
	f, err := os.Create(dir)
	check(err)

	defer f.Close()

	for _, d := range data {
		_, err = f.Write(d)
		check(err)
	}

	f.Sync()
}

func loadDescFile(dir string) []FileDesc {

	f, err := os.Open(dir)
	check(err)

	defer f.Close()

	scanner := bufio.NewScanner(f)

	var fds []FileDesc

	for scanner.Scan() {
		fd := stringToFileDesc(scanner.Text())
		fds = append(fds, fd)
	}

	return fds
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
