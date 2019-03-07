package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vulski/laverna/pkg/scraper"
	"github.com/vulski/laverna/pkg/scrapers/fullcomicpro"
	"github.com/vulski/laverna/pkg/scrapers/xoxocomics"
)

func init() {
	// Register your scraper
	scraper.RegisterScraper(fullcomicpro.Scraper{})
	scraper.RegisterScraper(xoxocomics.Scraper{})
}

func downloadBook(url string) {
	scrp, err := scraper.CreateScraper(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	book, err := scrp.GetBook(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Finished hydrating book, downloading: " + book.Title)
	book.Download("comics")
	fmt.Println("Finished downloading: " + book.Title)
}

func main() {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

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

		go downloadBook(comicUrl)
		fmt.Println("Sent that sucka")
	}
}
