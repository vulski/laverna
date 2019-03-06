package scrapers

import (
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"gitlab.com/PaperStreetHouse/laverna/pkg/comic"
	"gitlab.com/PaperStreetHouse/laverna/pkg/scraper"
)

type FullComicProScraper struct {
}

func (d FullComicProScraper) GetComic(url string) (*comic.Comic, error) {
	book := comic.Book{Url: url}
	doc := scraper.FetchDocument(url)

	doc.Find(".scroll-eps > a").Each(func(i int, selection *goquery.Selection) {
		chp := comic.Chapter{Url: d.Domain() + selection.AttrOr("href", "http://example.com"), Book: &book}

		doc = scraper.FetchDocument(chp.Url + "?readType=1")

		doc.Find("[id=imgPages] > img").Each(func(i int, selection *goquery.Selection) {
			pageUrl, exists := selection.Attr("src")
			pageIndx := strconv.Itoa(i + 1)
			if exists {
				chp.Pages = append(chp.Pages, &Page{Url: pageUrl, Number: pageIndx, Chapter: &chp})
			}
		})

		book.Chapters = append(book.Chapters, &chp)
	})

	return &book, nil
}

func (d FullComicProScraper) Domain() string {
	return "fullcomic.pro"
}
