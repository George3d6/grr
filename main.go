package main

import (
	ui "github.com/gizak/termui"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {

	/*
	* Intialize the termui lib
	*Display the current directory
	*Set up some enviornmental variables used to determine on which page
	* and in which directory we are + the index of the file currentl selected
	 */

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	var tout time.Duration = time.Duration(ui.TermHeight() * ui.TermWidth() / 50)
	var active int = 0

	dir, _ := os.Getwd()

	var page int = 0
	var middle int = ui.TermHeight()/2 - 3

	index_1 := 0
	index_2 := 0

	fdir := listDir(dir)
	sdir := listDir(dir)

	if fdir[index_1].dir {
		sdir = listDir(dir + "/" + fdir[index_1].name)
	}

	ls0 := getList(fdir, index_1)
	ls1 := getList(sdir, index_2)

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, ls0),
			ui.NewCol(6, 0, ls1),
		),
		ui.NewRow(
			ui.NewCol(12, 0, helpList()),
		),
	)
	ui.Body.Width = ui.TermWidth()
	ui.Body.Align()
	ui.Render(ui.Body)

	/*
	* Add some event listener for navigation
	* Not in a separate folder because they acess a lot of
	* vars from main
	 */

	//Change some env variables on resize
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		if fdir[index_1].dir {
			sdir = listDir(dir + "/" + fdir[index_1].name)
		} else {
			sdir = []FsItem{{"Is a file", false}}
		}

		ls0 = getList(fdir, index_1, index_1-middle)
		ls1 = getList(sdir, index_2, index_2-middle)

		ui.Body.Rows[0] = ui.NewRow(
			ui.NewCol(6, 0, ls0),
			ui.NewCol(6, 0, ls1),
		)

		tout = time.Duration(ui.TermHeight() * ui.TermWidth() / 50)
		middle = ui.TermHeight()/2 - 3
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Render(ui.Body)
	})
	//Quit program
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		//Stop loop should intialize the closing of the ui
		ui.StopLoop()
	})
	//

	//Move down
	ui.Handle("/sys/kbd/s", func(ui.Event) {
		if active != 0 {
			return
		}
		active = 1
		if page == 0 {
			if index_1 < (len(fdir) - 1) {
				index_1++
			}
		} else {
			if index_2 < (len(sdir) - 1) {
				index_2++
			}
		}

		if fdir[index_1].dir {
			sdir = listDir(dir + "/" + fdir[index_1].name)
		} else {
			sdir = []FsItem{{"Is a file", false}}
		}

		ls0 = getList(fdir, index_1, index_1-middle)
		ls1 = getList(sdir, index_2, index_2-middle)

		ui.Body.Rows[0] = ui.NewRow(
			ui.NewCol(6, 0, ls0),
			ui.NewCol(6, 0, ls1),
		)
		ui.Body.Align()
		ui.Render(ui.Body)
		time.AfterFunc(tout*time.Millisecond, func() {
			active = 0
		})
	})
	//

	//Move up
	ui.Handle("/sys/kbd/w", func(ui.Event) {
		if active != 0 {
			return
		}
		active = 1
		if page == 0 {
			if index_1 > 0 {
				index_1--
			}
		} else {
			if index_2 > 0 {
				index_2--
			}
		}

		if fdir[index_1].dir {
			sdir = listDir(dir + "/" + fdir[index_1].name)
		} else {
			sdir = []FsItem{{"Is a file", false}}
		}

		ls0 = getList(fdir, index_1, index_1-middle)
		ls1 = getList(sdir, index_2, index_2-middle)

		ui.Body.Rows[0] = ui.NewRow(
			ui.NewCol(6, 0, ls0),
			ui.NewCol(6, 0, ls1),
		)
		ui.Body.Align()
		ui.Render(ui.Body)
		time.AfterFunc(tout*time.Millisecond, func() {
			active = 0
		})
	})
	//

	//Move forward
	ui.Handle("/sys/kbd/d", func(ui.Event) {
		if !fdir[index_1].dir || len(sdir) < 1 {
			err := exec.Command("xdg-open", dir+"/"+fdir[index_1].name).Start()
			if err != nil {
				panic(err)
			}
			return
		}
		dir = dir + "/" + fdir[index_1].name
		fdir = sdir
		index_1 = 0
		index_2 = 0

		if fdir[index_1].dir {
			sdir = listDir(dir + "/" + fdir[index_1].name)
		}
		ls0 = getList(fdir, index_1)
		ls1 = getList(sdir, index_2)
		ui.Body.Rows[0] = ui.NewRow(
			ui.NewCol(6, 0, ls0),
			ui.NewCol(6, 0, ls1),
		)
		ui.Body.Align()
		ui.Render(ui.Body)
	})
	//

	//Move backward
	ui.Handle("/sys/kbd/a", func(ui.Event) {
		findir := strings.Split(dir, "/")
		dir = ""
		for ind, val := range findir {
			if ind < (len(findir) - 1) {
				dir = dir + "/" + val
			}
		}
		sdir = fdir
		fdir = listDir(dir)
		index_1 = 0
		index_2 = 0
		if fdir[index_1].dir {
			sdir = listDir(dir + "/" + fdir[index_1].name)
		}

		ls0 = getList(fdir, index_1)
		ls1 = getList(sdir, index_2)

		ui.Body.Rows[0] = ui.NewRow(
			ui.NewCol(6, 0, ls0),
			ui.NewCol(6, 0, ls1),
		)
		ui.Body.Align()
		ui.Render(ui.Body)
	})
	//

	ui.Handle("/sys/kbd/e", func(ui.Event) {
		if page == 0 {
			page = 1
		} else if page == 1 {
			page = 0
		} else {
			log.Fatal("Page number was not 0 or 1 it was: " + strconv.Itoa(page))
		}
	})

	var shellWasOpened bool = false
	ui.Handle("/sys/kbd/b", func(ui.Event) {
		if shellWasOpened {
			ui.Body.Rows[2] = ui.NewRow(
				ui.NewCol(12, 0, shellList(1, dir+" shell: ")),
			)
		} else {
			shellWasOpened = true
			ui.Body.AddRows(
				ui.NewRow(
					ui.NewCol(12, 0, shellList(1, dir+" shell: ")),
				),
			)
		}
		ui.Body.Align()
		ui.Render(ui.Body)

		var input string = ""
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
		var b []byte = make([]byte, 1)
		for {
			os.Stdin.Read(b)

			if string(b) == "\x0D" {
				dshell, _ := exec.Command("bash", "-c", "echo $SHELL").Output()
				tshold := string(dshell)
				var defaultShell string = tshold[strings.LastIndex(tshold, "/")+1 : len(tshold)-1]

				out, err := exec.Command(defaultShell, "-c", "cd "+dir+" && "+input).Output()
				if err != nil {
					ui.Body.Rows[2] = ui.NewRow(
						ui.NewCol(12, 0, shellList(1, "Invalid command")),
					)
					ui.Body.Align()
					ui.Render(ui.Body)
					break
				}
				ui.Body.Rows[2] = ui.NewRow(
					ui.NewCol(12, 0, shellList(2, string(out))),
				)
				ui.Body.Align()
				ui.Render(ui.Body)
				break
			}
			if string(b) == "\x1b" {
				ui.Body.Rows[2] = ui.NewRow(
					ui.NewCol(12, 0, shellList(1, "------------------------------------------------------------------------------------------------------------------------------------------------------")),
				)
				ui.Body.Align()
				ui.Render(ui.Body)
				break
			}
			if string(b) == "\x7f" {
				//This is ASCII for backspace
				//Which actually is ASCII for delete but on most
				//distros its ASCII for backspace\
				//computers are weird
				if len(input) == 0 {
				} else {
					input = input[0 : len(input)-1]
				}
			} else {
				input += string(b)
			}
			ui.Body.Rows[2] = ui.NewRow(
				ui.NewCol(12, 0, shellList(1, dir+" shell: "+input)),
			)
			ui.Body.Align()
			ui.Render(ui.Body)
		}
	})

	ui.Loop()

}
