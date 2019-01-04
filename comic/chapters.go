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
	CE.UpdateResults("Fetching Chapters")
	doc := thek.FetchDocument(url)

	CE.UpdateResults("Processing Chapters")
	chapters := make([]string, 0)
	doc.Find("div.chapter > a").Each(func(i int, selection *goquery.Selection) {
		chapters = append(chapters, selection.AttrOr("href", "http://example.com"))
	})
	CE.UpdateResults("Processed Chapters")

	return chapters
}

func DownloadChapter(chapter Chapter) {

	CE.UpdateResults("Fetching Chapter Page")
	doc := thek.FetchDocument(chapter.Uri)

	//CE.UpdateResults("Push found images to queue")
	doc.Find("[id=selectPage] > option").Each(func(i int, selection *goquery.Selection) {
		pageUrl, exists := selection.Attr("value")
		pageIdx := selection.Text()

		if exists {
			imageWaitGroup.Add(1)

			go func() {
				images <- Image{
					pageUrl: pageUrl,
					chapter: chapter,
					pageIdx: pageIdx,
				}
			}()
			TotalImages++
		}
	})
}