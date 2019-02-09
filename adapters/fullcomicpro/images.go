package fullcomicpro

import (
	"github.com/pkg/errors"
	"laverna/bus"
	"laverna/thek"
	"log"
	"os"
	"strconv"
)

func DownloadImage(image bus.Image) {
	downloadPath, err := GetDownloadPath(image)
	bus.Messages <- "Downloading " + downloadPath

	if err != nil {
		bus.Stats.DownloadedPages++
	} else {
		thek.DownloadPage(thek.Page{
			Uri:      image.PageUrl,
			FilePath: downloadPath,
		})
		bus.Stats.DownloadedPages++
	}
}

func GetDownloadPath(image bus.Image) (string, error) {
	chapterIdxString := strconv.Itoa(image.Chapter.ChapterIdx)
	chapterDownloadDirectory := bus.DownloadDirectory + image.Chapter.ComicName + "/" + chapterIdxString + "/"
	//splitUrl := strings.Split(image.PageUrl, ".")
	//extension := splitUrl[len(splitUrl) - 1]
	log.Println(image.PageUrl)
	filePath := chapterDownloadDirectory + image.PageIdx

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return "", errors.New("File already exists : " + filePath)
	}

	return filePath, nil
}
