package comic

type Comic struct {
	Title    string
	Author   string
	Url      string
	Chapters []*Chapter
}

type Chapter struct {
	Number int
	Url    string
	Pages  []*Page
}

type Page struct {
	Number   int
	ImageUrl string
}

//TODO
func (cm *Comic) Download(dir string) error {
	for _, chp := range cm.Chapters {
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
