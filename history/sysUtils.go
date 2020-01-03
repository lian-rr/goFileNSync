package history

import (
	"errors"
	"os"
	"path/filepath"
)

func fileInfoToMetafile(files []os.FileInfo) []Metafile {
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
