package comic

import (
	"laverna/adapters/fullcomicpro"
	"laverna/bus"
	"laverna/thek"
	"os"
	"strconv"
	"strings"
	"time"
)

const WorkerCount = 5

var ComicDownloaders = []*ComicDownloader{&fullcomicpro.FullComicProDownloader}

func Init() {
	thek.Init()

	bus.ChapterInit()
	bus.ImagesInit()
}

func Wait() {
	thek.Wait()

	bus.ChapterWait()
	bus.ImagesWait()
}

// This looks like its fullcomic.pro specific, should :100: be in the interface
func getComicName(url string) string {
	nameParts := strings.Split(url, "/")

	return nameParts[len(nameParts)-1]
}

func getChapters(comic_url string) []string {
	chaptersRes := make([]string, 0)
	finalChapters := make([]string, 0)
	counter := 1
	chaptersRes = fullcomicpro.GetChapters(comic_url + "?page=" + strconv.Itoa(counter))
	counter++

	finalChapters = append(chaptersRes, finalChapters...)
	bus.Stats.TotalChapters += len(finalChapters)

	return finalChapters
}

func Download(url string) {
	bus.Stats.PushEvent("Downloading Yo")

	url = strings.Trim(url, " ")
	err := os.MkdirAll(bus.DownloadDirectory, 0777)

	if err != nil {
		return
	}

	finalChapters := getChapters(url)

	for idx, chapter := range finalChapters {
		name := getComicName(url)
		bus.Chapters <- bus.Chapter{
			Uri:              "http://fullcomic.pro" + chapter,
			ChapterIdx:       idx + 1,
			ComicName:        name,
			DownloadFunction: fullcomicpro.DownloadChapter,
		}
		bus.ChapterWaitGroup.Add(1)
	}
}

var running = true
var startedUpdating = false

func Update() {
	startedUpdating = true
	for running {
		CE.UpdateResults()
		time.Sleep(150 * time.Millisecond)
	}
}
