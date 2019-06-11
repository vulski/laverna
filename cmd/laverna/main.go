package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/doctorbarber/laverna/internal/cli"
	"github.com/doctorbarber/laverna/pkg/comic"
	"github.com/doctorbarber/laverna/pkg/scrapers/fullcomicpro"
	"github.com/doctorbarber/laverna/pkg/scrapers/xoxocomics"
)

func init() {
	// Register your scraper
	comic.RegisterScraper(fullcomicpro.New())
	comic.RegisterScraper(xoxocomics.New())
	initLogger()
}

func initLogger() {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// defer f.Close()
	log.SetOutput(f)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Which directory to you want to download to: [./comics]")
	scanner.Scan()
	dir := scanner.Text()
	if dir == "" {
		dir = "comics"
	}
	for {
		fmt.Print("Enter Comic Url: ")
		comicUrl := ""
		scanner.Scan()
		comicUrl = scanner.Text()

		if scanner.Err() != nil {
			panic(scanner.Err())
		}

		comicUrl = strings.Trim(comicUrl, " ")

		go cli.DownloadBook(comicUrl, dir)
	}
}
