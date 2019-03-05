package main

import (
	"bufio"
	_ "bufio"
	"fmt"
	_ "fmt"
	"log"
	"os"
	"strings"

	"gitlab.com/PaperStreetHouse/laverna/comic"
)

func main() {

	comic.Init()
	defer comic.Wait()

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
	comic.Download(comicUrl)
}
