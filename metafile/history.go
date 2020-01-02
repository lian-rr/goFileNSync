package metafile

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	rootDir      = ".goFileNsync"
	localHistory = rootDir + "/localHistory"
)

func lsMetafiles(path string) ([]Metafile, error) {
	files, err := lsDir(path)

	if err != nil {
		return nil, err
	}

	return fileInfoToMetafile(files), nil
}

func fileInfoToMetafile(files []FileInfo) []Metafile {
	var metafiles []Metafile
	for _, file := range files {
		metafiles = append(metafiles, Metafile{
			Name:    file.Name(),
			Size:    int(file.Size()),
			ModTime: file.ModTime(),
		})
	}

	return metafiles
}

func prepareMetaFolder(path string) {
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
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
		return nil, errors.New("Error listing the directory")
	}

	return files, nil
}

func saveHistory(mfiles []Metafile, path string) error {
	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		return fmt.Errorf("Error creating file with path: %s", path)
	}

	for _, mf := range mfiles {
		_, err = f.Write(metafileToBytes(mf))

		if err != nil {
			return fmt.Errorf("Error writing file with name: %s's metadata", mf.Name)
		}
	}

	f.Sync()

	return nil
}

func loadHistory(path string) ([]Metafile, error) {
	f, err := os.Open(path)
	defer f.Close()

	if err != nil {
		return nil, fmt.Errorf("Error opening file with path: %s", path)
	}

	scanner := bufio.NewScanner(f)

	var metafiles []Metafile

	for scanner.Scan() {
		mf := stringToFileDesc(scanner.Text())
		metafiles = append(metafiles, mf)
	}

	return metafiles, nil
}
