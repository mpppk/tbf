package crawl

import (
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/mpppk/tbf/tbf"
	"github.com/pkg/errors"
)

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
func circlesFetchingTasks(circlesURL string) (chromedp.Tasks, *circlesTasksResult) {
	circlesTasksResult := NewCirclesTasksResult()
	circleListItemSel := "li.circle-list-item"
	detailUrlsSel := joinSelectors(circleListItemSel, "a.circle-list-item-link")
	circleSpacesSel := joinSelectors(circleListItemSel, "span.circle-space-label")
	circleNamesSel := joinSelectors(circleListItemSel, "span.circle-name")
	penNamesSel := joinSelectors(circleListItemSel, "p.circle-list-item-penname")
	genresSel := joinSelectors(circleListItemSel, "p.circle-list-item-penname")

	return chromedp.Tasks{
		chromedp.Navigate(circlesURL),
		chromedp.WaitVisible(`li.circle-list-item`, chromedp.AtLeast(200)), // FIXME
		AttributeValueAll(detailUrlsSel, "href", &(circlesTasksResult.detailUrls), nil, chromedp.ByQueryAll),
		Texts(circleSpacesSel, &(circlesTasksResult.spaces), chromedp.ByQueryAll),
		Texts(circleNamesSel, &(circlesTasksResult.names), chromedp.ByQueryAll),
		Texts(penNamesSel, &(circlesTasksResult.penNames), chromedp.ByQueryAll),
		Texts(genresSel, &(circlesTasksResult.genres), chromedp.ByQueryAll),
	}, circlesTasksResult
}

func circlesDetailFetchingTasks(fullCircleDetailURL string) (chromedp.Tasks, *tbf.CircleDetail) {
	circleDetail := &tbf.CircleDetail{}
	circleDetailCardSel := "mat-card.circle-detail-card"
	circleDetailTableSel := joinSelectors(circleDetailCardSel, "tbody")
	tableQueryTmpl := "tr:nth-of-type(%d)>td:nth-of-type(2)"

	circleImageSel := joinSelectors(circleDetailCardSel, "div.circle-detail-image>img")
	circleNameSel := joinSelectors(circleDetailTableSel, "span.circle-name")
	circleSpaceSel := joinSelectors(circleDetailTableSel, fmt.Sprintf(tableQueryTmpl, 2))
	circlePennameSel := joinSelectors(circleDetailTableSel, fmt.Sprintf(tableQueryTmpl, 3))
	circleWebURLSel := joinSelectors(circleDetailTableSel, fmt.Sprintf(tableQueryTmpl, 4), "a")
	circleGenreSel := joinSelectors(circleDetailTableSel, fmt.Sprintf(tableQueryTmpl, 5))
	circleGenreFreeFormatSel := joinSelectors(circleDetailTableSel, fmt.Sprintf(tableQueryTmpl, 6))

	return chromedp.Tasks{
		chromedp.Navigate(fullCircleDetailURL),
		chromedp.WaitVisible(`mat-card-content.mat-card-content`),
		chromedp.AttributeValue(circleImageSel, "src", &(circleDetail.ImageURL), nil, chromedp.ByQueryAll),
		chromedp.Text(circleNameSel, &(circleDetail.Circle.Name), chromedp.ByQueryAll),
		chromedp.Text(circleSpaceSel, &(circleDetail.Circle.Space), chromedp.ByQueryAll),
		chromedp.Text(circlePennameSel, &(circleDetail.Circle.Penname), chromedp.ByQueryAll),
		chromedp.Text(circleWebURLSel, &(circleDetail.WebURL), chromedp.ByQueryAll),
		chromedp.Text(circleGenreSel, &(circleDetail.Circle.Genre), chromedp.ByQueryAll),
		chromedp.Text(circleGenreFreeFormatSel, &(circleDetail.GenreFreeFormat), chromedp.ByQueryAll),
	}, circleDetail
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
