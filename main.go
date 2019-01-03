package main

import (
	"bufio"
	_ "bufio"
	"comicArchiver/comic"
	"fmt"
	_ "fmt"
	"github.com/jroimartin/gocui"
	"log"
	"os"
	"strings"
	_ "strings"
)

var downloadDir = "./downloads/"

var Editor gocui.Editor

var msg string
var msgView *gocui.View

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("input", 0, maxY - 2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		_, err := g.SetCurrentView("input")

		if err != nil {
			return err
		}

		//logger.Logger.Println(" CHANGE:", "input", x, y, maxX, maxY)

		v.Editor = gocui.EditorFunc(simpleEditor)

		v.FgColor = gocui.Attribute(15 + 1)
		// v.BgColor = gocui.Attribute(0)
		v.BgColor = gocui.ColorDefault

		v.Autoscroll = false
		v.Editable = true
		v.Wrap = false
		v.Frame = false

	}

	if msgView, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(msgView, msg)
	}
	//return nil

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func simpleEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {

	case key == gocui.KeyEnter:
		msg = v.Buffer()

		//log.Println(msg)

		mainView, _ := g.View("hello")

		_, err := fmt.Fprintln(mainView, msg)

		if err != nil {
			log.Fatalln(err)
		}

		if len(msg) <= 0 {
			// return errors.New("input line empty")
			v.Clear()
			v.SetCursor(0, 0)
		}

	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	//case key == gocui.KeyEnter:


	case key == gocui.KeyArrowDown:

	case key == gocui.KeyArrowUp:

	case key == gocui.KeyArrowLeft:

	case key == gocui.KeyArrowRight:


	}
}

var g *gocui.Gui

func main() {

	comic.Init()
	defer comic.Wait()

	//g, err := gocui.NewGui(gocui.OutputNormal)
	//if err != nil {
	//	log.Panicln(err)
	//}
	//defer g.Close()
	//
	//g.SetManagerFunc(layout)
	//
	//if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
	//	log.Panicln(err)
	//}
	//
	////msg = "Not Hello World"
	//
	//if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
	//	log.Panicln(err)
	//}




	_ = os.MkdirAll(downloadDir, 0777)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter Comic Url: ")
	comicUrl := ""
	scanner.Scan()
	comicUrl = scanner.Text()

	if scanner.Err() != nil {
		log.Println(scanner.Err())
	}

	comicUrl = strings.Trim(comicUrl, " ")
	comic.Download(comicUrl)
}
