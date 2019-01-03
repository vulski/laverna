package comic

import (
	"comicArchiver/thek"
	"log"
	"os"
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
	err := os.MkdirAll(DownloadDirectory, 0777)

	if err != nil {
		log.Println("Could not make downloads directory")
		return
	}

	chaptersRes := GetChapters(url)

	for idx, chapter := range chaptersRes {
		name := getComicName(url)
		chapters <- Chapter{
			Uri:        chapter,
			ChapterIdx: idx + 1,
			ComicName:       name,
		}
	}

	if len(chapters) == 0 {
		log.Println("No chapters found - Possibly invalid url")
	}
}