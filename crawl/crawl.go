package crawl

import (
	"context"

	"fmt"

	"path"

	"strings"

	"github.com/chromedp/chromedp"
	"github.com/mpppk/tbf/tbf"
	"github.com/pkg/errors"
)

type TBFCrawler struct {
	browser *chromedp.CDP
	baseURL string
}

func NewTBFCrawler(ctx context.Context, baseURL string) (*TBFCrawler, error) {
	c, err := chromedp.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "chromedep new error:")
	}
	return &TBFCrawler{
		browser: c,
		baseURL: baseURL,
	}, nil
}

func (t *TBFCrawler) FetchCircles(ctx context.Context, circlesURL string) ([]*tbf.Circle, error) {
	var circles []*tbf.Circle
	err := t.browser.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(circlesURL),
		chromedp.WaitVisible(`li.circle-list-item`),
		chromedp.Evaluate(
			`Array.from(document.querySelectorAll('li.circle-list-item')).map((l) => ({detailUrl: l.querySelector('a.circle-list-item-link').getAttribute('href'), space: l.querySelector('span.circle-space-label').textContent, name: l.querySelector('span.circle-name').textContent, penname: l.querySelector('p.circle-list-item-penname').textContent, genre: l.querySelector('p.circle-list-item-genre').textContent}))`,
			&circles,
		),
	})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to fetch circles from %v", circlesURL))
	}

	return circles, nil
}

func (t *TBFCrawler) fetchAttributeValue(ctx context.Context, sel interface{}, name string) (attributeValue string, err error) {
	var ok bool
	if err := t.browser.Run(ctx, chromedp.AttributeValue(sel, "src", &attributeValue, &ok)); err != nil {
		return "", errors.Wrapf(err, "failed to fetch image src from %s", sel)
	}
	if !ok {
		return "", fmt.Errorf("target DOM not found: %s", sel)
	}
	return
}

func (t *TBFCrawler) fetchText(ctx context.Context, sel interface{}) (text string, err error) {
	if err := t.browser.Run(ctx, chromedp.Text(sel, &text)); err != nil {
		return "", errors.Wrapf(err, "failed to fetch text from %q", sel)
	}
	return
}

func (t *TBFCrawler) getNavigateToCircleDetailTasks(circleDetailURL string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(fmt.Sprintf("%s/%s", t.baseURL, circleDetailURL)),
		chromedp.WaitVisible(`mat-card-content.mat-card-content`),
	}
}

func (t *TBFCrawler) FetchCircleDetail(ctx context.Context, circle *tbf.Circle) (*tbf.CircleDetail, error) {
	if err := t.browser.Run(ctx, t.getNavigateToCircleDetailTasks(circle.DetailURL)); err != nil {
		return nil, errors.Wrapf(err, "failed to navigate to %s", circle.DetailURL)
	}

	circleDetailCardSel := "mat-card.circle-detail-card"
	circleDetailTableSel := joinSelectors(circleDetailCardSel, "tbody")
	tableQueryTmpl := "tr:nth-of-type(%d)>td:nth-of-type(2)"

	circleImageSel := joinSelectors(circleDetailCardSel, "div.circle-detail-image>img")
	imageURL, err := t.fetchAttributeValue(ctx, circleImageSel, "src")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch detail of circle: %#v", circle)
	}

	circleNameSel := joinSelectors(circleDetailTableSel, "span.circle-name")
	name, err := t.fetchText(ctx, circleNameSel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch circle name from %q", circleNameSel)
	}

	circleSpaceSel := joinSelectors(circleDetailTableSel, fmt.Sprintf(tableQueryTmpl, 2))
	space, err := t.fetchText(ctx, circleSpaceSel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch circle space from %q", circleSpaceSel)
	}

	circlePennameSel := joinSelectors(circleDetailTableSel, fmt.Sprintf(tableQueryTmpl, 3))
	penname, err := t.fetchText(ctx, circlePennameSel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch circle penname from %q", circlePennameSel)
	}

	circleWebURLSel := joinSelectors(circleDetailTableSel, fmt.Sprintf(tableQueryTmpl, 4), "a")
	webURL, err := t.fetchText(ctx, circleWebURLSel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch circle web URL from %q", circleWebURLSel)
	}

	circleGenreSel := joinSelectors(circleDetailTableSel, fmt.Sprintf(tableQueryTmpl, 5))
	genre, err := t.fetchText(ctx, circleGenreSel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch circle genre from %q", circleGenreSel)
	}

	circleGenreFreeFormatSel := joinSelectors(circleDetailTableSel, fmt.Sprintf(tableQueryTmpl, 6))
	genreFreeFormat, err := t.fetchText(ctx, circleGenreFreeFormatSel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch circle genre free format from %q", circleGenreFreeFormatSel)
	}

	circleDetail := &tbf.CircleDetail{
		Circle: tbf.Circle{
			Name:      name,
			Space:     space,
			DetailURL: path.Join(t.baseURL, circle.DetailURL),
			Penname:   penname,
			Genre:     genre,
		},
		ImageURL:        imageURL,
		WebURL:          webURL,
		GenreFreeFormat: genreFreeFormat,
	}

	return circleDetail, err
}

func (t *TBFCrawler) Shutdown(ctx context.Context) error {
	return t.browser.Shutdown(ctx)
}

func (t *TBFCrawler) Wait() error {
	return t.browser.Wait()
}

func joinSelectors(selectors ...string) string {
	return strings.Join(selectors, " ")
}
