package comic

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Chapter struct {
	Title  string
	Url    string
	Number int
	Pages  []*Page
	Book   *Book
}

func (c *Chapter) GetPage(number int) (*Page, error) {
	for _, page := range c.Pages {
		if page.Number == number {
			return page, nil
		}
	}

	return nil, errors.New("Page not found.")
}

func (c *Chapter) Download(dir string) error {
	log.Println("Downloading Chapter: " + strconv.Itoa(c.Number))
	var title string
	if c.Title == "" {
		title = strconv.Itoa(c.Number)
	} else {
		title = c.Title
	}
	dir = filepath.Join(dir, title)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	for _, page := range c.Pages {
		err := page.Download(dir)
		if err != nil {
			return err
		}
	}
	log.Println("Finished downloading chapter")
	return nil
}
