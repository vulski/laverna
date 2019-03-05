package fullcomicpro

import (
	"laverna/bus"
	"laverna/thek"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var FullComicProDownloader struct {
}

func (d *FullComicProDownloader) Domain() string {
	return "fullcomic.pro"
}

func (d *FullComicProDownloader) GetChapters(string) []string {
	doc := thek.FetchDocument(url)

	chapters := make([]string)
	doc.Find(".scroll-eps > a").Each(func(i int, selection *goquery.Selection) {
		//ComicStats.TotalChapters++
		chapters = append(chapters, selection.AttrOr("href", "http://example.com"))
	})

	return chapters
}

func (d *FullComicProDownloader) DownloadChapter(chapter bus.Chapter) {
	doc := thek.FetchDocument(chapter.Uri + "?readType=1")

	bus.Stats.PushEvent("Downloading Chapter " + strconv.Itoa(chapter.ChapterIdx))
	doc.Find("[id=imgPages] > img").Each(func(i int, selection *goquery.Selection) {
		pageUrl, exists := selection.Attr("src")
		pageIdx := strconv.Itoa(i + 1)

		//log.Println(pageUrl, pageIdx)

		if exists {
			bus.ImageWaitGroup.Add(1)

			go func() {
				bus.Images <- bus.Image{
					PageUrl:          pageUrl,
					Chapter:          chapter,
					PageIdx:          pageIdx,
					DownloadFunction: DownloadImage,
				}
			}()
			bus.Stats.TotalPages++
			//TotalImages++
		}
	})
}
