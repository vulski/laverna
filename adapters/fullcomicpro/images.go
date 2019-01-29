package fullcomicpro

import (
	"github.com/pkg/errors"
	"laverna/bus"
	"laverna/thek"
	"os"
	"strconv"
)

func DownloadImage(image bus.Image) {
	downloadPath, err := GetDownloadPath(image)

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
	filePath := chapterDownloadDirectory + image.PageIdx + ".jpg"

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return "", errors.New("File already exists : " + filePath)
	}

	return filePath, nil
}
