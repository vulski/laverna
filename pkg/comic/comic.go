package comic

type Book struct {
	Title    string
	Author   string
	Url      string
	Chapters []*Chapter
}

type Chapter struct {
	Book   *Book
	Number int
	Url    string
	Pages  []*Page
}

type Page struct {
	Chapter *Chapter
	Number  int
	Url     string
}

//TODO
func (book *Book) Download(dir string) error {
	for _, chp := range book.Chapters {
		err := chp.Download(dir)
		if err != nil {
			return err
		}
	}

	return nil
}

//TODO
func (c *Chapter) Download(dir string) error {
	for _, page := range c.Pages {
		err := page.Download(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

//TODO
func (p *Page) Download(dir string) error {
	return nil
}
