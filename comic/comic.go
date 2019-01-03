package comic

import (
	"comicArchiver/thek"
	"log"
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