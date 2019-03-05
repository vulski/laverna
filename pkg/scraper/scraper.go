package scraper

import (
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"

	"gitlab.com/PaperStreetHouse/laverna/pkg/comic"
)

type Scraper interface {
	// Domain the scraper should be used for.
	Domain() string

	GetComic(string) (*comic.Comic, error)
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

func FetchDocument(url string) *goquery.Document {
	r, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	doc, docerr := goquery.NewDocumentFromReader(r.Body)

	if docerr != nil {
		log.Fatalln(docerr)
	}
	return doc
}
