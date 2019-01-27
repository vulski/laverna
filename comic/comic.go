package comic

import (
	"laverna/thek"
	"log"
	"os"
	"strconv"
	"strings"
)

// Globals
const WorkerCount = 5

type Stats struct {
	RunningWorkers string
	Messages []string
	DownloadedPages int
	TotalPages int
	TotalChapters int
}

var ComicStats Stats

func Init() {
	thek.Init()

	ChapterInit()
	ImagesInit()
}

func Wait() {
	thek.Wait()

	ChapterWait()
	ImagesWait()
}

func getComicName(url string) string {
	nameParts := strings.Split(url, "/")
	namePartsLength := len(nameParts)
	name := nameParts[namePartsLength-1]

	return name
}

func GetPagination(url string) {

}

func Download(url string) {
	url = strings.Trim(url, " ")
	err := os.MkdirAll(DownloadDirectory, 0777)

	if err != nil {
		//log.Println("Could not make downloads directory")
		return
	}

	CE.UpdateResults("Start Getting to it")

	chaptersRes := make([]string, 0)
	finalChapters := make([]string, 0)
	run := true
	counter := 1
	for run {
		chaptersRes = GetChapters(url + "?page=" + strconv.Itoa(counter))
		run = len(chaptersRes) > 0
		counter++

		if run {
			finalChapters = append(chaptersRes, finalChapters...)
		}
	}

	CE.UpdateResults("Pushing chapters to channel")
	for idx, chapter := range finalChapters {
		name := getComicName(url)
		chapters <- Chapter{
			Uri:        chapter,
			ChapterIdx: idx + 1,
			ComicName:  name,
		}
		chapterWaitGroup.Add(1)
	}

	//log.Println("Here")

	chapterCount := strconv.Itoa(len(finalChapters))
	CE.UpdateResults("Downloading " + chapterCount + " chapters")

	if chapterCount == "0" {
		log.Println("No chapters found - Possibly invalid url", url)
		//CE.UpdateResults("Oh no! " + url)

	}
}
