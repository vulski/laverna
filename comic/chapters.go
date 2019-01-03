package comic

import (
	"comicArchiver/thek"
	"github.com/PuerkitoBio/goquery"
	"sync"
)

// Globals
var chapters = make(chan Chapter, 0)
var chapterWaitGroup = sync.WaitGroup{}

type Chapter struct {
	Uri        string
	ChapterIdx int
	ComicName       string
}

func ChapterInit() {
	for i := 0; i < WorkerCount; i++ {
		go chapterWorker()
	}
}

func ChapterWait() {
	chapterWaitGroup.Wait()
}

func chapterWorker() {
	for {
		select {
		case chapter := <-chapters:
			DownloadChapter(chapter)
			chapterWaitGroup.Done()
		}
	}
}

func GetChapters(url string) []string {
	doc := thek.FetchDocument(url)

	chapters := make([]string, 0)
	doc.Find("div.chapter > a").Each(func(i int, selection *goquery.Selection) {
		chapters = append(chapters, selection.AttrOr("href", "http://example.com"))
	})

	return chapters
}

func DownloadChapter(chapter Chapter) {

	doc := thek.FetchDocument(chapter.Uri)

	// Get Last Page
	doc.Find("[id=selectPage] > option").Each(func(i int, selection *goquery.Selection) {
		pageUrl, exists := selection.Attr("value")
		pageIdx := selection.Text()

		if exists {
			imageWaitGroup.Add(1)
			images <- Image{
				pageUrl: pageUrl,
				chapter: chapter,
				pageIdx: pageIdx,
			}
		}
	})
}