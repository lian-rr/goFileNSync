package metafile

func WritingRouting(path string) {
	metafiles, err := lsMetafiles(path)

	check(err)

	prepareMetaFolder(rootDir)

	saveHistory(metafiles, localHistory)
}

// func ReadingRoutine() {
// 	fds := loadDescFile(localHistory)

// 	for _, fd := range fds {
// 		fmt.Printf("%s %d %s\n", fd.Name, fd.Size, fd.ModTime.String())
// 	}
// }

func check(e error) {
	if e != nil {
		panic(e)
	}
}
