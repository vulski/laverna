package scraper

import (
	"errors"
	"net/url"

	"gitlab.com/PaperStreetHouse/laverna/pkg/comic"
)

type Scraper interface {
	// Domain the scraper should be used for.
	Domain() string

	GetComic(string) (*comic.Comic, error)
}

var scrapers = []Scraper{FullComicProScraper{}}

func CreateScraper(URL string) (Scraper, error) {
	u, err := url.Parse(URL)

	if err != nil {
		return nil, err
	}

	domain := u.Hostname()
	for _, scraper := range scrapers {
		if scraper.Domain() == domain {
			return scraper, nil
		}
	}

	return nil, errors.New("Couldn't find a scraper for that url.")
}

type FullComicProScraper struct {
}

func (d FullComicProScraper) GetComic(url string) (*comic.Comic, error) {
	return nil, nil
}

func (d FullComicProScraper) Domain() string {
	return "fullcomic.pro"
}
