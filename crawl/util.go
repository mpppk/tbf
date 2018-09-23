package crawl

import (
	"strings"
)

func joinSelectors(selectors ...string) string {
	return strings.Join(selectors, " ")
}
