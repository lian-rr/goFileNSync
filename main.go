package main

import (
	mh "./history"
)

func main() {

	metafiles, err := mh.GetMetafiles("test/")

	if err != nil {
		panic(err)
	}

	mh.SaveHistory(metafiles, mh.LocalHistoryPath)

	// metafiles, err = fh.LoadHistory(fh.LocalHistoryPath)

}
