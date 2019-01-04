package comic

import (
	"comicArchiver/thek"
	"log"
	"os"
	"strconv"
	"strings"
)

// Globals
const WorkerCount = 5

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

func Download(url string) {
	url = strings.Trim(url, " ")
	err := os.MkdirAll(DownloadDirectory, 0777)

	if err != nil {
		//log.Println("Could not make downloads directory")
		return
	}


	CE.UpdateResults("Start Getting to it")

	chaptersRes := GetChapters(url)

	CE.UpdateResults("Pushing chapters to channel")
	for idx, chapter := range chaptersRes {
		name := getComicName(url)
		chapters <- Chapter{
			Uri:        chapter,
			ChapterIdx: idx + 1,
			ComicName:       name,
		}
		chapterWaitGroup.Add(1)
	}

	//log.Println("Here")

	chapterCount := strconv.Itoa(len(chaptersRes))
	CE.UpdateResults("Downloading " + chapterCount + " chapters")

	if chapterCount == "0" {
		log.Println("No chapters found - Possibly invalid url", url)
		//CE.UpdateResults("Oh no! " + url)

	}
}