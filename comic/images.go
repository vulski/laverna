package comic

import (
	"comicArchiver/thek"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"sync"
)

var images = make(chan Image, 0)
var imageWaitGroup = sync.WaitGroup{}
const DownloadDirectory = "./downloads/"


type Image struct {
	pageUrl string
	pageIdx string
	chapter Chapter
}

func ImagesInit() {
	for i := 0; i < WorkerCount; i++ {
		go imageWorker()
	}
}

func ImagesWait() {
	imageWaitGroup.Wait()
}

func getImage(pageUrl string) (string, bool) {
	doc := thek.FetchDocument(pageUrl)
	imgSel := doc.Find(".page-chapter > img").First()

	imgUrl, exists := imgSel.Attr("src")


	return imgUrl, exists
}

func getDownloadPath(image Image) (string, error) {
	chapterIdxString := strconv.Itoa(image.chapter.ChapterIdx)
	chapterDownloadDirectory := DownloadDirectory + image.chapter.ComicName + "/" + chapterIdxString + "/"
	filePath := chapterDownloadDirectory + image.pageIdx + ".jpg"

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return "", errors.New("File already exists : " + filePath)
	}

	return filePath, nil
}

func imageWorker() {
	for {
		select {
		case image := <-images:
			img, exists := getImage(image.pageUrl)

			if exists {

				downloadPath, err := getDownloadPath(image)

				if err != nil {
					return
				}

				thek.DownloadPage(thek.Page{
					Uri:img,
					FilePath:downloadPath,
				})
				imageWaitGroup.Done()
			}
		}
	}
}