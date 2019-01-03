package comic

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"strings"
)

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

var g *gocui.Gui

type Label struct {
	name string
	x, y int
	w, h int
	body string
}

func NewLabel(name string, x, y int, body string) *Label {
	lines := strings.Split(body, "\n")

	w := 0
	for _, l := range lines {
		if len(l) > w {
			w = len(l)
		}
	}
	h := len(lines) + 1
	w = w + 1

	return &Label{name: name, x: x, y: y, w: w, h: h, body: body}
}

func (l *Label) Layout(g *gocui.Gui) error {
	v, err := g.SetView(l.name, l.x, l.y, l.x+l.w, l.y+l.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprint(v, l.body)
	}
	return nil
}

type Input struct {
	name      string
	x, y      int
	w         int
	maxLength int
}

func NewInput(name string, x, y, w, maxLength int) *Input {
	return &Input{name: name, x: x, y: y, w: w, maxLength: maxLength}
}

func (i *Input) Layout(g *gocui.Gui) error {
	v, err := g.SetView(i.name, i.x, i.y, i.x+i.w, i.y+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editor = i
		v.Editable = true
	}
	return nil
}

func (i *Input) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	cx, _ := v.Cursor()
	ox, _ := v.Origin()
	limit := ox+cx+1 > i.maxLength
	switch {
	case ch != 0 && mod == 0 && !limit:
		v.EditWrite(ch)
	case key == gocui.KeySpace && !limit:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyEnter:
		input := v.Buffer()

		commandParts := strings.Split(input, " ")

		fmt.Println(commandParts, len(commandParts))

		if len(commandParts) >= 2 {

			command := commandParts[0]

			log.Println(command)

			switch command {
				case "get":
					//NewLabel("hello", 9, 6, "Hello World")
					//g.Update(SetFocus("hello"))
					Download(commandParts[1])
			}

		}

		v.Clear()
		_ = v.SetCursor(0, 0)
	}
}

func SetFocus(name string) func(g *gocui.Gui) error {
	return func(g *gocui.Gui) error {
		_, err := g.SetCurrentView(name)
		return err
	}
}

func InitUi() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	_, height := g.Size()

	label := NewLabel("label", 1, height - 3, "Command")
	input := NewInput("input", 9, height - 3, 100, 1000)
	focus := gocui.ManagerFunc(SetFocus("input"))
	g.SetManager(label, input, focus)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}