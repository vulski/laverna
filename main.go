package main

import (
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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

func downloadChapter(chapter Chapter) {
	doc := fetchDocument(chapter.Uri)

	// Get Last Page
	doc.Find("[id=selectPage] > option").Each(func(i int, selection *goquery.Selection) {
		pageUrl, exists := selection.Attr("value")

		if exists {
			downloadPage(pageUrl, chapter.ChapterIdx)
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

func downloadPage(pageUrl string, chapterIdx int) {
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

		chapterIdxString := strconv.Itoa(chapterIdx)

		mkerr := os.MkdirAll("/tmp/chapter-downloads/"+chapterIdxString, 0777)

		if mkerr != nil {
			log.Println("Failed making download directory")
			return
		}

		f, openerr := os.OpenFile("/tmp/chapter-downloads/" + chapterIdxString + "/" + idx + ".jpg", os.O_RDONLY|os.O_CREATE|os.O_WRONLY, 0777)

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

type Chapter struct {
	Uri string
	ChapterIdx int
}

var chapterUrls = make(chan Chapter, 0)

func chapterWorker() {
	for {
		select {
		case chapter := <- chapterUrls:
			downloadChapter(chapter)
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

	for idx, chapter := range chapters {
		wg.Add(1)
		chapterUrls <- Chapter{
			Uri : chapter,
			ChapterIdx : idx + 1,
		}
	}

	wg.Wait()

}
