package comic

import (
	"errors"
	"log"
	"os"
	"path/filepath"
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

func download(book *Book, dir string, chps chan *Chapter, errs chan error, quit chan int) error {
	for {
		select {
		case chp := <-chps:
			go func() {
				err := chp.Download(dir)
				if err != nil {
					errs <- err
				}
			}()
		case err := <-errs:
			return err
		case <-quit:
			log.Println("Finished downloading " + book.Title)
			return nil
		}
	}
}

func (book *Book) Download(dir string) error {
	errs := make(chan error)
	chps := make(chan *Chapter, len(book.Chapters))
	quit := make(chan int)

	log.Println("Downloading " + book.Title)
	dir = filepath.Join(dir, book.Title)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	go func() {
		for _, chp := range book.Chapters {
			chps <- chp
		}
		quit <- 0
	}()

	return download(book, dir, chps, errs, quit)
}
