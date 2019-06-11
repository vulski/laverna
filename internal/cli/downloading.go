package cli

import (
	"fmt"

	"github.com/doctorbarber/laverna/pkg/comic"
)

func DownloadBook(url string, dir string) {
	scrp, err := comic.CreateScraper(url)
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
	book.Download(dir)
	fmt.Println("Finished downloading: " + book.Title)
}
