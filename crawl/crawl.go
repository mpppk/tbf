package crawl

import (
	"context"

	"fmt"

	"path"

	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/mpppk/tbf/tbf"
	"github.com/pkg/errors"
)

type TBFCrawler struct {
	browser *chromedp.CDP
	baseURL string
}

type circlesTasksResult struct {
	detailUrls []string
	spaces     []string
	names      []string
	penNames   []string
	genres     []string
}

func (c *circlesTasksResult) validate() error {
	detailUrlsLen := len(c.detailUrls)
	spacesLen := len(c.spaces)
	namesLen := len(c.names)
	penNamesLen := len(c.penNames)
	genresLen := len(c.genres)
	if detailUrlsLen != spacesLen ||
		detailUrlsLen != namesLen ||
		detailUrlsLen != penNamesLen ||
		detailUrlsLen != genresLen {
		return fmt.Errorf("invalid circles information. "+
			"len(detailUrls):%d, len(spaces):%d, len(names):%d, len(penNames):%d, len(genres):%d",
			detailUrlsLen, spacesLen, namesLen, penNamesLen, genresLen)
	}
	return nil
}

func NewCirclesTasksResult() *circlesTasksResult {
	return &circlesTasksResult{
		detailUrls: []string{},
		spaces:     []string{},
		names:      []string{},
		penNames:   []string{},
		genres:     []string{},
	}
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

func (t *TBFCrawler) getFetchCirclesTasks(circlesURL string) (chromedp.Tasks, *circlesTasksResult) {
	circlesTasksResult := NewCirclesTasksResult()
	circleListItemSel := "li.circle-list-item"
	detailUrlsSel := joinSelectors(circleListItemSel, "a.circle-list-item-link")
	circleSpacesSel := joinSelectors(circleListItemSel, "span.circle-space-label")
	circleNamesSel := joinSelectors(circleListItemSel, "span.circle-name")
	penNamesSel := joinSelectors(circleListItemSel, "p.circle-list-item-penname")
	genresSel := joinSelectors(circleListItemSel, "p.circle-list-item-penname")

	return chromedp.Tasks{
		chromedp.Navigate(circlesURL),
		chromedp.WaitVisible(`li.circle-list-item`, chromedp.AtLeast(100)), // FIXME
		AttributeValueAll(detailUrlsSel, "href", &(circlesTasksResult.detailUrls), nil, chromedp.ByQueryAll),
		Texts(circleSpacesSel, &(circlesTasksResult.spaces), chromedp.ByQueryAll),
		Texts(circleNamesSel, &(circlesTasksResult.names), chromedp.ByQueryAll),
		Texts(penNamesSel, &(circlesTasksResult.penNames), chromedp.ByQueryAll),
		Texts(genresSel, &(circlesTasksResult.genres), chromedp.ByQueryAll),
	}, circlesTasksResult
}

func (t *TBFCrawler) FetchCircles(ctx context.Context, circlesURL string) ([]*tbf.Circle, error) {
	tasks, res := t.getFetchCirclesTasks(circlesURL)
	err := t.browser.Run(ctx, tasks)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute circles fetching tasks from "+circlesURL)
	}

	circles, err := fetchResultToCircles(res)
	return circles, errors.Wrap(err, "error occurred after circles are fetched")
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

func (t *TBFCrawler) fetchAttributeValues(ctx context.Context, sel interface{}, name string) ([]string, error) {
	var attributeMaps []map[string]string
	if err := t.browser.Run(ctx, chromedp.AttributesAll(sel, &attributeMaps, chromedp.ByQueryAll)); err != nil {
		return nil, errors.Wrapf(err, "failed to fetch attributes from %q", sel)
	}

	values, err := mapsToSliceByKey(attributeMaps, name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to extract attribute(%s) from nodes that fetched by %q", name, sel)
	}
	return values, nil
}

func (t *TBFCrawler) fetchText(ctx context.Context, sel interface{}) (text string, err error) {
	if err := t.browser.Run(ctx, chromedp.Text(sel, &text)); err != nil {
		return "", errors.Wrapf(err, "failed to fetch text from %q", sel)
	}
	return
}

func (t *TBFCrawler) fetchNodes(ctx context.Context, sel interface{}) (nodes []*cdp.Node, err error) {
	if err := t.browser.Run(ctx,
		chromedp.Nodes(sel, &nodes, chromedp.ByQueryAll),
	); err != nil {
		return nil, errors.Wrapf(err, "failed to fetch node from %q", sel)
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

func mapsToSliceByKey(ms []map[string]string, key string) (values []string, err error) {
	for _, m := range ms {
		v, ok := m[key]
		if !ok {
			return nil, errors.New(fmt.Sprintf("failed to get value of %q from %v", key, m))
		}
		values = append(values, v)
	}
	return
}

func fetchResultToCircles(res *circlesTasksResult) (circles []*tbf.Circle, err error) {
	if err := res.validate(); err != nil {
		return nil, errors.Wrap(err, "failed to convert circle fetching tasks result to circles")
	}
	for i := range res.detailUrls {
		circles = append(circles, &tbf.Circle{
			DetailURL: res.detailUrls[i],
			Space:     res.spaces[i],
			Name:      res.names[i],
			Penname:   res.penNames[i],
			Genre:     res.genres[i],
		})
	}
	return
}
