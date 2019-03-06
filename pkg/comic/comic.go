package comic

import (
	"io/ioutil"
	"mime"
	"net/http"
	"os"
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

func (book *Book) Download(dir string) error {
	for _, chp := range book.Chapters {
		err := chp.Download(dir)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Chapter) Download(dir string) error {
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

	bytes, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// Why 512? Shouldn't it be len(bytes) ? s h r u g
	extension, err = mime.ExtensionsByType(http.DetectContentType(bytes[0:512]))
	if err != nil {
		return err
	}
	extension = extension[0]

	f, err = os.OpenFile(page.FilePath, os.O_RDONLY|os.O_CREATE|os.O_WRONLY, 0775)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err := f.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
