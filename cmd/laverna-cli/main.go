package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vulski/laverna/pkg/comic"
	"github.com/vulski/laverna/pkg/scraper"
	"github.com/vulski/laverna/pkg/scrapers/fullcomicpro"
	"github.com/vulski/laverna/pkg/scrapers/xoxocomics"
)

var hydrateChan chan string

func init() {
	hydrateChan = make(chan string, 100)

	// Register your scraper
	scraper.RegisterScraper(fullcomicpro.Scraper{})
	scraper.RegisterScraper(xoxocomics.Scraper{})
}

func createBook(urls chan string) {
	url := <-urls
	scrp, err := scraper.CreateScraper(url)
	if err != nil {
		panic(err)
	}

	book, err := scrp.GetBook(url)
	if err != nil {
		panic(err)
	}

	go func(*comic.Book) {
		book.Download("comics")
	}(book)
}

func main() {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	go createBook(hydrateChan)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter Comic Url: ")
		comicUrl := ""
		scanner.Scan()
		comicUrl = scanner.Text()

		if scanner.Err() != nil {
			panic(scanner.Err())
		}

		comicUrl = strings.Trim(comicUrl, " ")

		hydrateChan <- comicUrl
		fmt.Println("Sent that sucka")
	}
}
