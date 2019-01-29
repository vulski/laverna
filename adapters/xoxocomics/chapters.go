package xoxocomics

import (
	"github.com/PuerkitoBio/goquery"
	"laverna/bus"
	"laverna/thek"
)

func GetChapters(url string) []string {
	doc := thek.FetchDocument(url)

	chapters := make([]string, 0)
	doc.Find("div.chapter > a").Each(func(i int, selection *goquery.Selection) {
		//ComicStats.TotalChapters++
		chapters = append(chapters, selection.AttrOr("href", "http://example.com"))
	})

	return chapters
}

func DownloadChapter(chapter bus.Chapter) {
	doc := thek.FetchDocument(chapter.Uri)
	doc.Find("[id=selectPage] > option").Each(func(i int, selection *goquery.Selection) {
		pageUrl, exists := selection.Attr("value")
		pageIdx := selection.Text()

		if exists {
			bus.ImageWaitGroup.Add(1)

			go func() {
				bus.Images <- bus.Image{
					PageUrl: pageUrl,
					Chapter: chapter,
					PageIdx: pageIdx,
					DownloadFunction:DownloadImage,
				}
			}()
			bus.Stats.TotalPages++
		}
	})
}
