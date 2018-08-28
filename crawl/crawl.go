package crawl

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
)

type Circle struct {
	Href    string
	Space   string
	Name    string
	Penname string
	Genre   string
}

type TBFCrawler struct {
	browser *chromedp.CDP
}

func NewTBFCrawler(ctx context.Context) (*TBFCrawler, error) {
	c, err := chromedp.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "chromedep new error:")
	}
	return &TBFCrawler{
		browser: c,
	}, nil
}

func (t *TBFCrawler) FetchCircles(ctx context.Context) ([]*Circle, error) {
	var circles []*Circle
	err := t.browser.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(`https://techbookfest.org/event/tbf04/circle`),
		chromedp.WaitVisible(`li.circle-list-item`),
		chromedp.Evaluate(
			`Array.from(document.querySelectorAll('li.circle-list-item')).map((l) => ({href: l.querySelector('a.circle-list-item-link').getAttribute('href'), space: l.querySelector('span.circle-space-label').textContent, name: l.querySelector('span.circle-name').textContent, penname: l.querySelector('p.circle-list-item-penname').textContent, genre: l.querySelector('p.circle-list-item-genre').textContent}))`,
			&circles,
		),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch circles")
	}

	return circles, nil
}

func (t *TBFCrawler) Shutdown(ctx context.Context) error {
	return t.browser.Shutdown(ctx)
}

func (t *TBFCrawler) Wait() error {
	return t.browser.Wait()
}
