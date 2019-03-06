package scrapers

import (
	"fmt"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"gitlab.com/PaperStreetHouse/laverna/pkg/comic"
	"gitlab.com/PaperStreetHouse/laverna/pkg/scraper"
)

type FullComicProScraper struct {
}

func (d FullComicProScraper) GetBook(Url string) (*comic.Book, error) {
	book := comic.Book{Url: Url}

	u, err := url.Parse(Url)

	if err != nil {
		return nil, err
	}

	doc := scraper.FetchDocument(Url)

	doc.Find(".scroll-eps > a").Each(func(i int, selection *goquery.Selection) {
		chp := comic.Chapter{Url: u.Scheme + "://" + u.Hostname() + selection.AttrOr("href", "http://example.com"), Book: &book}
		doc = scraper.FetchDocument(chp.Url + "?readType=1")

		doc.Find("[id=imgPages] > img").Each(func(i int, selection *goquery.Selection) {
			pageUrl, exists := selection.Attr("src")
			pageIndx := i + 1
			fmt.Println(pageUrl)
			if exists {
				fmt.Println("Found!")
				chp.Pages = append(chp.Pages, &comic.Page{Url: pageUrl, Number: pageIndx, Chapter: &chp})
			}
		})

		book.Chapters = append(book.Chapters, &chp)
	})

	return &book, nil
}

func (d FullComicProScraper) Domain() string {
	return "fullcomic.pro"
}
