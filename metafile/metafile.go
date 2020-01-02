package metafile

import (
	"fmt"
	"time"
)

// Metafile holds the information of the file
type Metafile struct {
	Name    string    // Name of the file
	ModTime time.Time // Last modification date
	Size    int       // Size of the file
}

// New creates a new Metafile
func New(name string, size int, modTime time.Time) Metafile {
	return Metafile{
		Name:    name,
		Size:    size,
		ModTime: modTime,
	}
}

func metafileToBytes(mf Metafile) []byte {
	return []byte(fmt.Sprintf("%s %d %d\n", mf.Name, mf.Size, mf.ModTime.UnixNano()))
}

func stringToMetafile(data string) (Metafile, error) {
	var name string
	var size int
	var modTime int64

	_, err := fmt.Sscanf(data, "%s %d %d\n", &name, &size, &modTime)

	if err != nil {
		return Metafile{}, fmt.Errorf("Error parsing string into Metafile")
	}

	return New(name, size, time.Unix(0, modTime)), nil
}
