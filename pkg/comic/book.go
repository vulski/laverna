package comic

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Book struct {
	Title    string
	Url      string
	Author   string
	Chapters []*Chapter
	Scraper  Scraper
}

func (book *Book) GetChapter(number int) (*Chapter, error) {
	for _, chp := range book.Chapters {
		if chp.Number == number {
			return chp, nil
		}
	}

	return nil, errors.New("Chapter not found.")
}

func (book *Book) Download(dir string) error {
	errs := make(chan error)
	var wg sync.WaitGroup

	log.Println("Downloading " + book.Title)
	dir = filepath.Join(dir, book.Title)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	for _, chp := range book.Chapters {
		wg.Add(1)
		go func() {
			err := chp.Download(dir)
			if err != nil {
				errs <- err
				log.Println("Had trouble downloading chapter #" + chp.Number)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if err = <-errs; err != nil {
		log.Println(err)
	}

	log.Println("Finished downloading " + book.Title)
	return nil
}
