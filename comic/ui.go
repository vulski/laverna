package comic

import (
	"fmt"
	"laverna/bus"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jroimartin/gocui"
)

const NumGoroutines = 10

var (
	done = make(chan struct{})
	wg   sync.WaitGroup

	mu  sync.Mutex // protects ctr
	ctr = 0
)

func InitUi() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

type CommandEditor struct {
	Insert bool
	g *gocui.Gui
}

func AddMessage(msg string) {
	bus.Stats.Messages = append(bus.Stats.Messages, msg)
}

func (ce *CommandEditor) UpdateResults() {
	ce.g.Update(func(g *gocui.Gui) error {
		v, err := g.View("ctr")
		if err != nil {
			return err
		}
		v.Clear()

		fmt.Fprintln(v, "Pages: " + strconv.Itoa(bus.Stats.DownloadedPages) + "/" + strconv.Itoa(bus.Stats.TotalPages))
		fmt.Fprintln(v, "Total Chapters: " + strconv.Itoa(bus.Stats.TotalChapters))

		fmt.Fprintln(v,"----------------------------------------")

		// Build Message from Stats
		if len(bus.Stats.Messages) > 5 {
			for _, msg := range bus.Stats.Messages[len(bus.Stats.Messages)-6 : len(bus.Stats.Messages)] {
				fmt.Fprintln(v, msg)
			}
		} else {
			for _, msg := range bus.Stats.Messages {
				fmt.Fprintln(v, msg)
			}
		}
		return nil
	})
}

func (ce *CommandEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	cx, _ := v.Cursor()
	ox, _ := v.Origin()
	limit := ox+cx+1 > 50000

	switch {

	case key == gocui.KeyEnter:
		input := v.Buffer()

		commandParts := strings.Split(input, " ")

		//fmt.Println(commandParts, len(commandParts))

		if len(commandParts) >= 2 {
			command := commandParts[0]

			//log.Println(command)

			switch command {
			case "get":
				//NewLabel("hello", 9, 6, "Hello World")
				//g.Update(SetFocus("hello"))
				url := strings.TrimSpace(commandParts[1])
				go Download(url)
				//ce.UpdateResults("|" + url + "|")
			}

		}

		v.Clear()
		_ = v.SetCursor(0, 0)
		break

	case ch != 0 && mod == 0 && !limit:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)

	}
}

var CE = &CommandEditor{}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("ctr", 2, 2, maxX - 5, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Command: get [url]")
	}

	if v, err := g.SetView("input", 0, maxY - 3, maxX - 5, maxY - 1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		CE.g = g
		v.Editor = CE

		//fmt.Fprintln(v, "Hello ")
	}

	g.SetCurrentView("input")

	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	close(done)
	return gocui.ErrQuit
}

func counter(g *gocui.Gui) {
	defer wg.Done()

	for {
		select {
		case <-done:
			return
		case <-time.After(500 * time.Millisecond):
			mu.Lock()
			n := ctr
			ctr++
			mu.Unlock()

			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("ctr")
				if err != nil {
					return err
				}
				v.Clear()
				fmt.Fprintln(v, n)
				return nil
			})
		}
	}
}
