package comic

import (
	"errors"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Book struct {
	Title    string
	Author   string
	Url      string
	Chapters []*Chapter
}

type Chapter struct {
	Book   *Book
	Number int
	Url    string
	Pages  []*Page
}

type Page struct {
	Chapter *Chapter
	Number  int
	Url     string
}

func (book *Book) GetChapter(number int) (*Chapter, error) {
	for _, chp := range book.Chapters {
		if chp.Number == number {
			return &chp, nil
		}
	}

	return nil, errors.New("Chapter not found.")
}

func (book *Book) Download(dir string) error {
	log.Println("Downloading " + book.Title)
	dir = filepath.Join(dir, book.Title)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	for _, chp := range book.Chapters {
		err := chp.Download(dir)
		if err != nil {
			return err
		}
	}

	log.Println("Finished downloading " + book.Title)
	return nil
}

func (c *Chapter) GetPage(number int) (*Page, error) {
	for _, page := range c.Pages {
		if page.Number == number {
			return &page, nil
		}
	}

	return nil, errors.New("Page not found.")
}

func (c *Chapter) Download(dir string) error {
	dir = filepath.Join(dir, strconv.Itoa(c.Number))
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
	return nil
}

func (p *Page) Download(dir string) error {
	r, err := http.Get(p.Url)

	if err != nil {
		return err
	}
	defer r.Body.Close()

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// Why 512? Shouldn't it be len(bytes) ? s h r u g
	exts, err := mime.ExtensionsByType(http.DetectContentType(bytes[0:512]))
	if err != nil {
		return err
	}
	ext := exts[0]

	f, err := os.OpenFile(filepath.Join(dir, strconv.Itoa(p.Number)+ext), os.O_RDONLY|os.O_CREATE|os.O_WRONLY, 0775)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
