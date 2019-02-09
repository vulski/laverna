package thek

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Structs
type Page struct {
	Uri string
	FilePath string
}

// Globals
const WorkerCount  = 5
const SleepTime = 500

var pages = make(chan Page, 0)
var pageWaitGroup = sync.WaitGroup{}

func Init() {
	for i := 0; i < WorkerCount; i++ {
		go pageWorker()
	}
}

func Wait() {
	pageWaitGroup.Wait()
}

// Public API

func DownloadPage(page Page) {
	//log.Println("Downloading Page")
	pageWaitGroup.Add(1)
	pages <- page
}

func FetchDocument(url string) *goquery.Document {
	r, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	doc, docerr := goquery.NewDocumentFromReader(r.Body)

	if docerr != nil {
		log.Fatalln(docerr)
	}
	return doc
}

// Private Methods

func savePage(page Page) error {

	//log.Println("Loading Page")

	r, err := http.Get(page.Uri)

	//log.Println("Page Loaded")

	if err != nil {
		//log.Println(err)
		//log.Println("Couldn't get page")
		return err
	}

	defer r.Body.Close()

	//log.Println("Reading into memory")
	bytes, readerr := ioutil.ReadAll(r.Body)

	contentType := http.DetectContentType(bytes[0:512])

	extension := "jpg"
	if contentType == "image/jpeg" {
		extension = "jpg"
	} else if contentType == "image/png" {
		extension = "png"
	}

	if readerr != nil {
		log.Fatalln(readerr)
	}

	page.FilePath = page.FilePath + "." + extension

	dirParts := strings.Split(page.FilePath, "/")
	dir := strings.Join(dirParts[0:len(dirParts) - 1], "/")

	//log.Println("Make download directory", dir)
	mkerr := os.MkdirAll(dir, 0777)

	if mkerr != nil {
		return errors.New("Failed making download directory")
	}

	//log.Println("Save file to disk")
	f, openerr := os.OpenFile(page.FilePath, os.O_RDONLY|os.O_CREATE|os.O_WRONLY, 0775)

	defer f.Close()

	if openerr != nil {
		log.Fatalln(openerr)
	}

	_, writeerr := f.Write(bytes)

	if writeerr != nil {
		log.Fatalln(writeerr)
	}

	//log.Println("Saved - now pausing")
	time.Sleep(time.Duration(SleepTime * time.Millisecond))

	return nil

}

func pageWorker() {
	for {
		select {
		case page := <-pages:
			saveerr := savePage(page)

			if saveerr != nil {
				log.Println(saveerr)
			}
			pageWaitGroup.Done()
		}
	}
}
