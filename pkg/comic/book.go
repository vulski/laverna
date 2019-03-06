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
	log.Println("Downloading " + book.Title)
	dir = filepath.Join(dir, book.Title)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	for _, chp := range book.Chapters {
		err := chp.Download(dir)
		if err != nil {
			return err
		}
	}

	log.Println("Finished downloading " + book.Title)
	return nil
}
