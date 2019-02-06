package comic

import (
	"fmt"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	"laverna/bus"
	"log"
	"net/http"

	//"log"
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
	// Initialize astilectron
	var a, err = astilectron.New(astilectron.Options{
		AppName: "Laverna",
		AppIconDefaultPath: "<your .png icon>", // If path is relative, it must be relative to the data directory
		AppIconDarwinPath:  "<your .icns icon>", // Same here
		BaseDirectoryPath: "./electron",
	})
	defer a.Close()


	if err != nil {
		log.Fatalln(err)
	}
	// Start astilectron
	err = a.Start()

	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		http.ListenAndServe("127.0.0.1:4000", http.FileServer(http.Dir("./resources/app")))
	}()

	// Create a new window
	w, err := a.NewWindow("http://127.0.0.1:4000", &astilectron.WindowOptions{
		Center: astilectron.PtrBool(true),
		Height: astilectron.PtrInt(600),
		Width:  astilectron.PtrInt(600),
	})

	if err != nil {
		log.Fatalln(err)
	}

	err = w.Create()

	if err != nil {
		log.Fatalln(err)
	}

	// Add a listener on the window
	w.On(astilectron.EventNameWindowEventResize, func(e astilectron.Event) (deleteListener bool) {
		astilog.Info("Window resized")
		return
	})

	// Init a new app menu
	// You can do the same thing with a window
	var m = a.NewMenu([]*astilectron.MenuItemOptions{
		{
			Label: astilectron.PtrStr("Separator"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Label: astilectron.PtrStr("Normal 1")},
				{
					Label: astilectron.PtrStr("Normal 2"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						astilog.Info("Normal 2 item has been clicked")
						return
					},
				},
				{Type: astilectron.MenuItemTypeSeparator},
				{Label: astilectron.PtrStr("Normal 3")},
			},
		},
		{
			Label: astilectron.PtrStr("Checkbox"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Checked: astilectron.PtrBool(true), Label: astilectron.PtrStr("Checkbox 1"), Type: astilectron.MenuItemTypeCheckbox},
				{Label: astilectron.PtrStr("Checkbox 2"), Type: astilectron.MenuItemTypeCheckbox},
				{Label: astilectron.PtrStr("Checkbox 3"), Type: astilectron.MenuItemTypeCheckbox},
			},
		},
		{
			Label: astilectron.PtrStr("Radio"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Checked: astilectron.PtrBool(true), Label: astilectron.PtrStr("Radio 1"), Type: astilectron.MenuItemTypeRadio},
				{Label: astilectron.PtrStr("Radio 2"), Type: astilectron.MenuItemTypeRadio},
				{Label: astilectron.PtrStr("Radio 3"), Type: astilectron.MenuItemTypeRadio},
			},
		},
		{
			Label: astilectron.PtrStr("Roles"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Label: astilectron.PtrStr("Minimize"), Role: astilectron.MenuItemRoleMinimize},
				{Label: astilectron.PtrStr("Close"), Role: astilectron.MenuItemRoleClose},
			},
		},
	})
	// Open dev tools
	w.OpenDevTools()

	m.Create()

	// This will send a message and execute a callback
	// Callbacks are optional

	// This will listen to messages sent by Javascript
	w.OnMessage(func(m *astilectron.EventMessage) interface{} {
		// Unmarshal
		var s string
		m.Unmarshal(&s)

		println(s)
		parts := strings.Split(s, " ")
		log.Println(parts)

		if len(parts) > 1 {
			Download(parts[1])
		}

		return "Pressed Yo"
	})

	// Blocking pattern
	a.Wait()

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
