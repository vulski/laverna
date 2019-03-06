package xoxocomics

import (
	"errors"
	"fmt"
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
	_, err := url.Parse(Url)
	if err != nil {
		return nil, err
	}

	book := comic.Book{Url: Url}

	doc, err := scraper.FetchDocument(Url)
	if err != nil {
		return nil, err
	}

	doc.Find(".title-detail").Each(func(i int, selection *goquery.Selection) {
		book.Title = selection.Text()
	})

	if book.Title == "" {
		log.Println("Couldn't find book.")
		return nil, errors.New("Couldn't find book")
	}

	log.Println("Found book " + book.Title)

	var outsideErr error
	// Get Chapters
	chps := doc.Find(".col-xs-9 > a")
	nodes := chps.Nodes
	// Reverse chapter order
	for i := len(nodes)/2 - 1; i >= 0; i-- {
		opp := len(nodes) - 1 - i
		nodes[i], nodes[opp] = nodes[opp], nodes[i]
	}
	chps.Nodes = nodes
	chps.EachWithBreak(func(i int, selection *goquery.Selection) bool {
		chp := comic.Chapter{Url: selection.AttrOr("href", "http://example.com"), Number: i + 1, Book: &book}

		// Get Pages for each chapter
		doc, err = scraper.FetchDocument(chp.Url)
		if err != nil {
			outsideErr = err
			return false
		}

		doc.Find("[id=selectPage] > option").EachWithBreak(func(i int, selection *goquery.Selection) bool {
			pageUrl, exists := selection.Attr("value")
			if !exists {
				return true
			}

			doc, err = scraper.FetchDocument(pageUrl)
			if err != nil {
				outsideErr = err
				return false
			}

			doc.Find("[id=page_" + strconv.Itoa(i+1) + "] > img").EachWithBreak(func(i int, selection *goquery.Selection) bool {
				fmt.Println("found page for chapter: " + strconv.Itoa(chp.Number))
				pageUrl, exists := selection.Attr("src")
				pageIndx := i + 1

				if exists {
					chp.Pages = append(chp.Pages, &comic.Page{Url: pageUrl, Number: pageIndx, Chapter: &chp})
				}

				return true
			})
			return true
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
