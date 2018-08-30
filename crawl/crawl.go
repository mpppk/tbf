package crawl

import (
	"context"

	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type Circle struct {
	Href    string
	Space   string
	Name    string
	Penname string
	Genre   string
}

type CircleDetail struct {
	Circle          `structs:",flatten" mapstructure:",squash"`
	ImageURL        string
	WebURL          string
	GenreFreeFormat string
}

type TBFCrawler struct {
	browser    *chromedp.CDP
	baseURL    string
	circlesURL string
}

func NewTBFCrawler(ctx context.Context) (*TBFCrawler, error) {
	c, err := chromedp.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "chromedep new error:")
	}
	return &TBFCrawler{
		browser:    c,
		baseURL:    `https://techbookfest.org`,
		circlesURL: `https://techbookfest.org/event/tbf04/circle`,
	}, nil
}

func NewCircleDetailFromMap(m map[string]string) (*CircleDetail, error) {
	var circleDetail CircleDetail
	err := mapstructure.Decode(m, &circleDetail)
	if err != nil {
		return nil, errors.Wrap(err, "field to decode circle detail map to CircleDetail struct")
	}
	return &circleDetail, nil
}

func (t *TBFCrawler) FetchCircles(ctx context.Context) ([]*Circle, error) {
	var circles []*Circle
	err := t.browser.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(t.circlesURL),
		chromedp.WaitVisible(`li.circle-list-item`),
		chromedp.Evaluate(
			`Array.from(document.querySelectorAll('li.circle-list-item')).map((l) => ({href: l.querySelector('a.circle-list-item-link').getAttribute('href'), space: l.querySelector('span.circle-space-label').textContent, name: l.querySelector('span.circle-name').textContent, penname: l.querySelector('p.circle-list-item-penname').textContent, genre: l.querySelector('p.circle-list-item-genre').textContent}))`,
			&circles,
		),
	})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to fetch circles from %v", t.circlesURL))
	}

	return circles, nil
}

func (t *TBFCrawler) FetchCircleDetail(ctx context.Context, circle *Circle) (*CircleDetail, error) {
	jsCmd := `
		(() => {
			const mat             = document.querySelector('mat-card.circle-detail-card');
			const imageURL        = mat.querySelector('div.circle-detail-image>img').getAttribute('src');
			const table           = mat.querySelector('tbody');
			const name            = table.querySelector('span.circle-name').textContent;
			const space           = table.querySelector('tr:nth-of-type(2)>td:nth-of-type(2)').textContent;
			const penname         = table.querySelector('tr:nth-of-type(3)>td:nth-of-type(2)').textContent;
			const webURL          = table.querySelector('tr:nth-of-type(4)>td:nth-of-type(2) a').getAttribute('href');
			const genre           = table.querySelector('tr:nth-of-type(5)>td:nth-of-type(2)').textContent;
			const genreFreeFormat = table.querySelector('tr:nth-of-type(6)>td:nth-of-type(2)').textContent;
			return {imageURL, name, space, penname, webURL, genre, genreFreeFormat};
		})();
	`
	var circleDetail *CircleDetail
	err := t.browser.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(fmt.Sprintf("%s/%s", t.baseURL, circle.Href)),
		chromedp.WaitVisible(`mat-card-content.mat-card-content`),
		chromedp.Evaluate(
			jsCmd,
			&circleDetail,
		),
	})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to fetch circle details from %s", circle.Href))
	}
	return circleDetail, err
}

func (t *TBFCrawler) Shutdown(ctx context.Context) error {
	return t.browser.Shutdown(ctx)
}

func (t *TBFCrawler) Wait() error {
	return t.browser.Wait()
}
