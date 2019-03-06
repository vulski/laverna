package xoxocomics

import (
	"errors"
	"log"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/vulski/laverna/pkg/comic"
	"github.com/vulski/laverna/pkg/scraper"
)

type Scraper struct {
}

func (d Scraper) Domain() string {
	return "xoxocomics.com"
}

func (d Scraper) GetBook(Url string) (*comic.Book, error) {
	u, err := url.Parse(Url)
	if err != nil {
		return nil, err
	}

	book := comic.Book{Url: Url}

	doc, err := scraper.FetchDocument(Url)
	if err != nil {
		return nil, err
	}

	doc.Find(".title > a").Each(func(i int, selection *goquery.Selection) {
		book.Title = selection.Text()
	})

	if book.Title == "" {
		log.Println("Couldn't find book.")
		return nil, errors.New("Couldn't find book")
	}

	log.Println("Found book " + book.Title)

	var outsideErr error
	// Get Chapters
	doc.Find(".scroll-eps > a").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		chp := comic.Chapter{Url: u.Scheme + "://" + u.Hostname() + selection.AttrOr("href", "http://example.com"), Number: i + 1, Book: &book}

		// Get Pages for each chapter
		doc, err = scraper.FetchDocument(chp.Url + "?readType=1")
		if err != nil {
			outsideErr = err
			return false
		}

		doc.Find("[id=imgPages] > img").Each(func(i int, selection *goquery.Selection) {
			pageUrl, exists := selection.Attr("src")
			pageIndx := i + 1
			if exists {
				chp.Pages = append(chp.Pages, &comic.Page{Url: pageUrl, Number: pageIndx, Chapter: &chp})
			}
		})

		log.Println("Found Chapter " + strconv.Itoa(chp.Number) + " -- " + strconv.Itoa(len(chp.Pages)) + " pages.")
		book.Chapters = append(book.Chapters, &chp)
		return true
	})

	if nil != outsideErr {
		return nil, outsideErr
	}

	return &book, nil
}
