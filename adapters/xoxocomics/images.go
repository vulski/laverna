package xoxocomics

import (
	"github.com/pkg/errors"
	"laverna/bus"
	"laverna/thek"
	"os"
	"strconv"
)



func GetImage(pageUrl string) (string, bool) {
	doc := thek.FetchDocument(pageUrl)
	imgSel := doc.Find(".page-chapter > img").First()

	imgUrl, exists := imgSel.Attr("src")

	return imgUrl, exists
}

func DownloadImage(image bus.Image) {
	img, exists := GetImage(image.PageUrl)

	//log.Println("Here", exists)

	if exists {

		downloadPath, err := GetDownloadPath(image)

		if err != nil {
			//DownloadedImages++
			//ComicStats.DownloadedPages++
			//di := strconv.Itoa(DownloadedImages)
			//ti := strconv.Itoa(TotalImages)

			//log.Println(di, ti)

			//AddMessage(di + " / " + ti + " downloaded images - " + idstring)
			//CE.UpdateResults(di + " / " + ti + " downloaded images - " + idstring)

		} else {
			//CE.UpdateResults("Downloading Page...")
			thek.DownloadPage(thek.Page{
				Uri:      img,
				FilePath: downloadPath,
			})

			//DownloadedImages++
			//ComicStats.DownloadedPages++

			//di := strconv.Itoa(DownloadedImages)
			//ti := strconv.Itoa(TotalImages)

			//log.Println(di, ti)

			//CE.UpdateResults(di + " / " + ti + " downloaded images - " + idstring)
		}
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