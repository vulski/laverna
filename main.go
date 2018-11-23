package main

import (
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

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

func downloadChapter(chapterUrl string) {
	doc := fetchDocument(chapterUrl)

	// Get Last Page
	doc.Find("[id=selectPage] > option").Each(func(i int, selection *goquery.Selection) {
		pageUrl, exists := selection.Attr("value")

		if exists {
			downloadPage(pageUrl)
		}
	})
	//lastPageString := sel.Text()
	//
	//lastPage, converr := strconv.Atoi(lastPageString)
	//if converr != nil {
	//	log.Fatalln(converr)
	//}
	//
	//for i := 1; i <= lastPage; i++ {
	//	pageUrl := chapterUrl + "/" + strconv.Itoa(i)
	//	downloadPage(pageUrl)
	//}
}

func downloadPage(pageUrl string) {
	doc := fetchDocument(pageUrl)
	imgSel := doc.Find(".page-chapter > img").First()

	imgUrl, exists := imgSel.Attr("src")
	idx, _ := imgSel.Attr("data-index")

	if exists {
		r, err := http.Get(imgUrl)

		if err != nil {
			log.Fatalln(err)
		}

		defer r.Body.Close()

		bytes, readerr := ioutil.ReadAll(r.Body)

		if readerr != nil {
			log.Fatalln(readerr)
		}

		f, openerr := os.OpenFile("/tmp/chapter-downloads/" + idx + ".jpg", os.O_RDONLY|os.O_CREATE|os.O_WRONLY, 0777)

		defer f.Close()

		if openerr != nil {
			log.Fatalln(openerr)
		}

		_, writeerr := f.Write(bytes)

		if writeerr != nil {
			log.Fatalln(writeerr)
		}

		log.Println("Downloaded : ", imgUrl)
	}

}

type ChapterUrl struct {
	Uri string
}

var chapterUrls = make(chan ChapterUrl, 0)

func chapterWorker() {
	for {
		select {
		case url := <- chapterUrls:
			downloadChapter(url.Uri)
			wg.Done()
		}
	}
}

var wg = sync.WaitGroup{}

func main() {
	//r, err := http.Get("http://xoxocomics.com/comic/miles-morales-ultimate-spider-man/issue-1/9821/1")
	//
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//defer r.Body.Close()
	//
	//doc, docErr := goquery.NewDocumentFromReader(r.Body)
	//
	//if docErr != nil {
	//	log.Fatalln(docErr)
	//}

	//doc.Find("img[data-original]").Each(func(i int, selection *goquery.Selection) {
	//	src, _ := selection.Attr("src")
	//
	//
	//	log.Println(src)
	//
	//})

	for i := 0; i < 5; i++ {
		go chapterWorker()
	}

	_ = os.MkdirAll("/tmp/chapter-downloads", 0777)
	chapters := getChapters("http://xoxocomics.com/comic/miles-morales-ultimate-spider-man")

	for _, chapter := range chapters[0:1] {
		wg.Add(1)
		chapterUrls <- ChapterUrl{
			Uri : chapter,
		}
		break
	}

	wg.Wait()

}
