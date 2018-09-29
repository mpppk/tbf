package crawl

import (
	"strings"

	"github.com/mpppk/tbf/tbf"
)

func joinSelectors(selectors ...string) string {
	return strings.Join(selectors, " ")
}

func FilterCircles(circles []*tbf.Circle, circleDetailMap map[string]*tbf.CircleDetail) (filteredCircles []*tbf.Circle) {
	for _, c := range circles {
		if _, ok := circleDetailMap[c.Space]; !ok {
			filteredCircles = append(filteredCircles, c)
		}
	}
	return
}
