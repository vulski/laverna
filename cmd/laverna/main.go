package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"gitlab.com/PaperStreetHouse/laverna/internal/pkg/scrapers"
	"gitlab.com/PaperStreetHouse/laverna/pkg/scraper"
)

func main() {
	// Register your scraper
	scraper.RegisterScraper(scrapers.FullComicProScraper{})

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter Comic Url: ")
	comicUrl := ""
	scanner.Scan()
	comicUrl = scanner.Text()

	if scanner.Err() != nil {
		log.Println(scanner.Err())
	}

	comicUrl = strings.Trim(comicUrl, " ")

	fmt.Println(comicUrl)

	scrp, err := scraper.CreateScraper(comicUrl)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Using scraper for host: " + scrp.Domain())
	comic, err := scrp.GetComic(comicUrl)
	fmt.Println(comic.Author)
	for _, chp := range comic.Chapters {
		fmt.Println(chp.Url)
	}
}
