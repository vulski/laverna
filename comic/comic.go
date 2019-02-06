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

// Globals
const WorkerCount = 5

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

func getComicName(url string) string {
	nameParts := strings.Split(url, "/")
	namePartsLength := len(nameParts)
	name := nameParts[namePartsLength-1]

	return name
}

func getChapters(comic_url string) []string {
	chaptersRes := make([]string, 0)
	finalChapters := make([]string, 0)
	//run := true
	counter := 1
	//for run {
	chaptersRes = fullcomicpro.GetChapters(comic_url + "?page=" + strconv.Itoa(counter))
	//run = len(chaptersRes) > 0
	counter++

	//if run {
	finalChapters = append(chaptersRes, finalChapters...)
	bus.Stats.TotalChapters += len(finalChapters)
	//}
	//}

	return finalChapters
}

func Download(url string) {
	//if(!startedUpdating) {
	//	go Update()
	//}
	bus.Stats.PushEvent("Downloading Yo")

	url = strings.Trim(url, " ")
	err := os.MkdirAll(bus.DownloadDirectory, 0777)

	if err != nil {
		return
	}

	finalChapters := getChapters(url)
	//log.Println(finalChapters)
	//return

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
		time.Sleep(150*time.Millisecond)
	}
}