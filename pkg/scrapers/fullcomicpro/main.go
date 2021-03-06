package fullcomicpro

import (
	"errors"
	"log"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/doctorbarber/laverna/pkg/comic"
)

func New() comic.Scraper {
	return scraper{}
}

type scraper struct {
}

func (d scraper) Domain() string {
	return "fullcomic.pro"
}

// We already populate the ImageUrl, so just return nil.
func (d scraper) FindImageUrl(p *comic.Page) error {
	return nil
}

// Find the comic book for fullcomic.pro with the given URL and "hydrate" a struct to return.
func (d scraper) GetBook(Url string) (*comic.Book, error) {
	u, err := url.Parse(Url)
	if err != nil {
		return nil, err
	}

	book := comic.Book{Url: Url, Scraper: d}

	doc, err := comic.FetchDocument(Url)
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
	chps := doc.Find(".scroll-eps > a")
	nodes := chps.Nodes
	// Reverse chapter order
	for i := len(nodes)/2 - 1; i >= 0; i-- {
		opp := len(nodes) - 1 - i
		nodes[i], nodes[opp] = nodes[opp], nodes[i]
	}
	chps.Nodes = nodes
	chps.EachWithBreak(func(i int, selection *goquery.Selection) bool {
		chpUrl, ok := selection.Attr("href")
		if !ok {
			return true
		}
		chp := comic.Chapter{Url: u.Scheme + "://" + u.Hostname() + chpUrl, Number: i + 1, Book: &book}

		// Get Pages for each chapter, with the ImageUrl
		doc, err = comic.FetchDocument(chp.Url + "?readType=1")
		if err != nil {
			outsideErr = err
			return false
		}

		doc.Find("[id=imgPages] > img").Each(func(i int, selection *goquery.Selection) {
			pageUrl, exists := selection.Attr("src")
			pageIndx := i + 1
			if exists {
				chp.Pages = append(chp.Pages, &comic.Page{Url: pageUrl, ImageUrl: pageUrl, Number: pageIndx, Chapter: &chp})
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
