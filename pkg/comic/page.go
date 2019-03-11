package comic

import (
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Page struct {
	// The site's default reading page url.
	Url string
	// Page number.
	Number  int
	Chapter *Chapter
	// The actual URL to the image.
	ImageUrl string
}

func (p *Page) Scraper() Scraper {
	return p.Chapter.Book.Scraper
}

func (p *Page) Download(dir string) error {
	log.Println("Downloading page: " + strconv.Itoa(p.Number))

	err := p.Scraper().FindImageUrl(p)
	if err != nil {
		return err
	}

	r, err := http.Get(p.ImageUrl)

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

	log.Println("Finished downloading page.")
	return nil
}
