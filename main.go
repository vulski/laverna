package main

import (
	"bufio"
	_ "bufio"
	"fmt"
	_ "fmt"
	"log"
	"os"
	"strings"

	"gitlab.com/PaperStreetHouse/laverna/pkg/scraper"
)

func main() {

	// comic.Init()
	// defer comic.Wait()

	// comic.InitUi()

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
	scraper, err := scraper.CreateScraper(comicUrl)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Using scraper for host: " + scraper.Domain())
	comic, err := scraper.GetComic(comicUrl)
	fmt.Println(comic)
}
