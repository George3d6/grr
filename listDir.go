package main

import (
	"io/ioutil"
	"log"
	"os"
)

type FsItem struct {
	name string
	dir  bool
}

func listDir(dirname string) []FsItem {
	fileInfoArr, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}
	var names []FsItem
	for _, file := range fileInfoArr {
		path := dirname + "/" + file.Name()
		fileInfo, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if fileInfo.IsDir() {
			names = append(names, FsItem{file.Name(), true})
		} else {
			names = append(names, FsItem{file.Name(), false})
		}
	}
	return names
}
