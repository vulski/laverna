package scraper

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"

	"github.com/vulski/laverna/pkg/comic"
)

type Scraper interface {
	// Domain the scraper should be used for.
	Domain() string

	// Get the book with the given URL.
	GetBook(string) (*comic.Book, error)
}

var Scrapers = []Scraper{}

func RegisterScraper(scraper Scraper) {
	Scrapers = append(Scrapers, scraper)
}

func CreateScraper(URL string) (Scraper, error) {
	u, err := url.Parse(URL)

	if err != nil {
		return nil, err
	}

	domain := u.Hostname()
	for _, scraper := range Scrapers {
		if scraper.Domain() == domain {
			return scraper, nil
		}
	}

	return nil, errors.New("Couldn't find a scraper for that url.")
}

func FetchDocument(url string) (*goquery.Document, error) {
	r, err := http.Get(url)

	if err != nil {
		return nil, nil
	}

	doc, docerr := goquery.NewDocumentFromReader(r.Body)

	if docerr != nil {
		return nil, nil
	}
	return doc, nil
}
