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

func main() {
	// Register your scraper
	scraper.RegisterScraper(fullcomicpro.Scraper{})
	scraper.RegisterScraper(xoxocomics.Scraper{})

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
	comic, err := scrp.GetBook(comicUrl)
	err = comic.Download("comics")
	if err != nil {
		panic(err)
	}
}
