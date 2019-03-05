package bus

import "sync"

var WorkerCount = 5
var Chapters = make(chan Chapter, 0)
var ChapterWaitGroup = sync.WaitGroup{}

type stats struct {
	RunningWorkers  string
	Messages        []string
	DownloadedPages int
	TotalPages      int
	TotalChapters   int
}

func (s *stats) PushEvent(msg string) {
	s.Messages = append(s.Messages, msg)
}

var Stats stats

type Chapter struct {
	Uri              string
	ChapterIdx       int
	ComicName        string
	DownloadFunction func(chapter Chapter)
}

func ChapterWorker() {
	for {
		select {
		case chapter := <-Chapters:
			chapter.DownloadFunction(chapter)
			ChapterWaitGroup.Done()
		}
	}
}

func ChapterInit() {
	for i := 0; i < WorkerCount; i++ {
		go ChapterWorker()
	}
}

func ChapterWait() {
	ChapterWaitGroup.Wait()
}
