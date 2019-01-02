package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var downloadDir = "./downloads/"
var chapterWorkers = 10
var imageWorkers = 5
var totalPages = 0

func fetchDocument(url string) *goquery.Document {
	r, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	doc, docerr := goquery.NewDocumentFromReader(r.Body)

	if docerr != nil {
		log.Fatalln(docerr)
	}
	return doc
}

func getChapters(url string) []string {
	doc := fetchDocument(url)

	chapters := make([]string, 0)
	doc.Find("div.chapter > a").Each(func(i int, selection *goquery.Selection) {
		chapters = append(chapters, selection.AttrOr("href", "http://example.com"))
	})

	return chapters
}

func downloadChapter(chapter Chapter) {

	log.Println("Preparing chapter ", chapter.ChapterIdx, " for downloading...")
-
	doc := fetchDocument(chapter.Uri)

	// Get Last Page
	doc.Find("[id=selectPage] > option").Each(func(i int, selection *goquery.Selection) {
		pageUrl, exists := selection.Attr("value")
		pageIdx := selection.Text()

		if exists {

			totalPages++
			wg.Add(1)
			imageDownloadChannel <- imageDownloader{
				pageUrl: pageUrl,
				chapter: chapter,
				pageIdx: pageIdx,
			}
		}
	})
}

func downloadPage(pageUrl string, idx string, chapter Chapter) {
	chapterIdxString := strconv.Itoa(chapter.ChapterIdx)
	chapterDownloadDirectory := downloadDir + chapter.Name + "/" + chapterIdxString + "/"
	filePath := chapterDownloadDirectory + idx + ".jpg"

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		//log.Println("File already exists : ", filePath)
		return
	}

	doc := fetchDocument(pageUrl)
	imgSel := doc.Find(".page-chapter > img").First()

	imgUrl, exists := imgSel.Attr("src")
	if exists {
		log.Println("Loading Page")

		r, err := http.Get(imgUrl)

		if err != nil {
			log.Fatalln(err)
		}

		defer r.Body.Close()

		bytes, readerr := ioutil.ReadAll(r.Body)

		if readerr != nil {
			log.Fatalln(readerr)
		}

		mkerr := os.MkdirAll(chapterDownloadDirectory, 0777)

		if mkerr != nil {
			log.Println("Failed making download directory")
			return
		}

		f, openerr := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE|os.O_WRONLY, 0777)

		defer f.Close()

		if openerr != nil {
			log.Fatalln(openerr)
		}

		_, writeerr := f.Write(bytes)

		if writeerr != nil {
			log.Fatalln(writeerr)
		}

		time.Sleep(1 * time.Second)
	}

}

type Chapter struct {
	Uri        string
	ChapterIdx int
	Name       string
}

type imageDownloader struct {
	pageUrl string
	pageIdx string
	chapter Chapter
}

var chapterUrls = make(chan Chapter, 0)
var imageDownloadChannel = make(chan imageDownloader, 10000)

func imageWorker() {
	for {
		select {
		case download := <-imageDownloadChannel:
			downloadPage(download.pageUrl, download.pageIdx, download.chapter)
			wg.Done()
		}
	}
}

func chapterWorker() {
	for {
		select {
		case chapter := <-chapterUrls:
			downloadChapter(chapter)
			wg.Done()
		}
	}
}

var wg = sync.WaitGroup{}

func main() {

	for i := 0; i < chapterWorkers; i++ {
		go chapterWorker()
	}

	for i := 0; i < imageWorkers; i++ {
		go imageWorker()
	}

	_ = os.MkdirAll(downloadDir, 0777)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter Comic Url: ")
	comicUrl := ""
	scanner.Scan()
	comicUrl = scanner.Text()

	if scanner.Err() != nil {
		log.Println(scanner.Err())
	}

	comicUrl = strings.Trim(comicUrl, " ")
	log.Println(comicUrl)
	chapters := getChapters(comicUrl)

	nameParts := strings.Split(comicUrl, "/")
	namePartsLength := len(nameParts)
	name := nameParts[namePartsLength-1]

	for idx, chapter := range chapters {
		wg.Add(1)
		chapterUrls <- Chapter{
			Uri:        chapter,
			ChapterIdx: idx + 1,
			Name:       name,
		}
	}

	if len(chapters) == 0 {
		log.Println("No chapters found - Possibly invalid url")
	}

	//wg.Wait()
	//wg.Add(totalPages)

	//for i := 0; i < imageWorkers; i++ {
	//	go imageWorker()
	//}

	wg.Wait()

}
