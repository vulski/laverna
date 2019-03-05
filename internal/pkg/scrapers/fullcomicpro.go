package scrapers

import (
	"github.com/PuerkitoBio/goquery"
	"gitlab.com/PaperStreetHouse/laverna/pkg/comic"
	"gitlab.com/PaperStreetHouse/laverna/pkg/scraper"
)

type FullComicProScraper struct {
}

func (d FullComicProScraper) GetComic(url string) (*comic.Comic, error) {
	cmc := comic.Comic{Url: url}
	doc := scraper.FetchDocument(url)

	doc.Find(".scroll-eps > a").Each(func(i int, selection *goquery.Selection) {
		cmc.Chapters = append(cmc.Chapters, &comic.Chapter{Url: d.Domain() + selection.AttrOr("href", "http://example.com"), Comic: &cmc})
	})

	return &cmc, nil
}

func (d FullComicProScraper) Domain() string {
	return "fullcomic.pro"
}
