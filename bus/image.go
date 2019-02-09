package bus

import "sync"

var Images = make(chan Image, 0)
var Messages = make(chan string, 0)
var ImageWaitGroup = sync.WaitGroup{}
const DownloadDirectory = "./downloads/"


type Image struct {
	PageUrl string
	PageIdx string
	Chapter Chapter
	DownloadFunction func(image Image)
}

func ImageWorker() {
	for {
		select {
		case image := <-Images:
			image.DownloadFunction(image)
			ImageWaitGroup.Done()
		}
	}
}

func ImagesInit() {
	for i := 0; i < WorkerCount; i++ {
		go ImageWorker()
	}
}

func ImagesWait() {
	ImageWaitGroup.Wait()
}
