package history

import (
	"bufio"
	"fmt"
	"os"
)

const (
	rootDir = ".goFileNsync"
	// LocalHistoryPath is the directory to store the history
	LocalHistoryPath = rootDir + "/localHistory"
)

// GetMetafiles return a list of the metafiles in the specified path
func GetMetafiles(path string) ([]Metafile, error) {
	files, err := lsDir(path)

	if err != nil {
		return nil, err
	}

	return fileInfoToMetafile(files), nil
}

// SaveHistory save the list of Metafile in the specified directory
func SaveHistory(metafiles []Metafile, path string) error {
	prepareMetaFolder(rootDir)

	err := saveMetafiles(metafiles, path)

	if err != nil {
		return fmt.Errorf("Error saving history. %s", err)
	}
	return nil
}

// LoadHistory load the historical list of Metafiles in the specified file
func LoadHistory(path string) ([]Metafile, error) {
	metafiles, err := readMetafiles(path)

	if err != nil {
		return nil, fmt.Errorf("Error loading the metafiles history. %s", err)
	}
	return metafiles, nil
}

func saveMetafiles(mfiles []Metafile, path string) error {
	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		return fmt.Errorf("Error creating file with path: %s", path)
	}

	for _, mf := range mfiles {
		_, err = f.Write(mf.ToBytes())

		if err != nil {
			return fmt.Errorf("Error writing file with name: %s's metadata", mf.Name)
		}
	}

	f.Sync()

	return nil
}

func readMetafiles(path string) ([]Metafile, error) {
	f, err := os.Open(path)
	defer f.Close()

	if err != nil {
		return nil, fmt.Errorf("Error opening file with path: %s", path)
	}

	scanner := bufio.NewScanner(f)

	var metafiles []Metafile

	for scanner.Scan() {
		data := scanner.Text()
		mf, err := StringToMetafile(data)
		if err != nil {
			return nil, fmt.Errorf("Error parsing %s to Metafile. %s", data, err)
		}
		metafiles = append(metafiles, mf)
	}

	return metafiles, nil
}
