package main

import (
	ui "github.com/gizak/termui"
)

func getList(files []FsItem, index ...int) *ui.List {
	var fileNames []string
	for _, f := range files {
		fileNames = append(fileNames, f.name)
	}
	if index[0] >= 0 && index[0] < len(files) {
		fileNames[index[0]] = "[" + fileNames[index[0]] + "](fg-green,bg-black)"
	}
	for in, _ := range fileNames {
		if in != index[0] {
			fileNames[in] = "[" + fileNames[in] + "](fg-black,bg-white)"
		}
	}
	ls := ui.NewList()
	if len(fileNames) > ui.TermHeight()-5 && len(index) > 1 && index[1] > -1 {
		ls.Items = fileNames[index[1]:]
	} else {
		ls.Items = fileNames
	}
	ls.Bg = ui.ColorWhite
	ls.Border = false
	ls.Height = ui.TermHeight() - 3
	ls.Y = 0
	ls.X = 0
	return ls
}

func helpList() *ui.List {
	ls := ui.NewList()
	ls.Items = []string{"[Move back an forth with w/a/s/d  |  Run shell script with b  | Quit with q | Change the page with e](fg-white,bg-blue)"}
	ls.Border = false
	ls.Bg = ui.ColorBlue
	ls.Overflow = "wrap"
	ls.BorderLabel = "Help"
	ls.Height = 2
	ls.Y = 0
	ls.X = 0
	return ls
}

func shellList(height int, content string) *ui.List {
	ls := ui.NewList()
	ls.Items = []string{"[" + content + "](fg-green,bg-black)"}
	ls.Border = false
	ls.Bg = ui.ColorBlack
	ls.Overflow = "wrap"
	ls.BorderLabel = "Help"
	ls.Height = height
	ls.Y = 0
	ls.X = 0
	return ls
}
