package crawl

import (
	"context"
	"log"

	"path"

	"fmt"

	"github.com/mpppk/chromedp"
	"github.com/mpppk/tbf/tbf"
	"github.com/pkg/errors"
)

type TBFCrawler struct {
	browser *chromedp.CDP
	baseURL string
}

func NewTBFCrawler(ctx context.Context, baseURL string) (*TBFCrawler, error) {
	c, err := chromedp.New(ctx, chromedp.WithLog(log.Printf))
	if err != nil {
		return nil, errors.Wrap(err, "chromedep new error:")
	}
	return &TBFCrawler{
		browser: c,
		baseURL: baseURL,
	}, nil
}

func (t *TBFCrawler) FetchCircles(ctx context.Context, circlesURL string) ([]*tbf.Circle, error) {
	tasks, res := circlesFetchingTasks(circlesURL)
	err := t.browser.Run(ctx, tasks)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute circles fetching tasks from "+circlesURL)
	}

	circles, err := fetchResultToCircles(res)
	return circles, errors.Wrap(err, "error occurred after circles are fetched")
}

func (t *TBFCrawler) FetchCircleDetail(ctx context.Context, circle *tbf.Circle) (*tbf.CircleDetail, error) {
	tasks, circleDetail := circlesDetailFetchingTasks(fmt.Sprintf("%s/%s", t.baseURL, circle.DetailURL))
	if err := t.browser.Run(ctx, tasks); err != nil {
		return nil, errors.Wrapf(err, "failed to navigate to %s", circle.DetailURL)
	}

	circleDetail.DetailURL = path.Join(t.baseURL, circle.DetailURL)
	return circleDetail, nil
}

func (t *TBFCrawler) Shutdown(ctx context.Context) error {
	return t.browser.Shutdown(ctx)
}

func (t *TBFCrawler) Wait() error {
	return t.browser.Wait()
}
