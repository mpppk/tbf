package tbf

import (
	"fmt"

	"path"

	"strings"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

var URLMap = map[string]string{
	"latest": "https://raw.githubusercontent.com/mpppk/tbf/master/data/latest_circles.csv",
	"tbf4":   "https://raw.githubusercontent.com/mpppk/tbf/master/data/tbf4_circles.csv",
}

type Circle struct {
	DetailURL string
	Space     string
	Name      string
	Penname   string
	Genre     string
}

type CircleDetail struct {
	Circle          `structs:",flatten" mapstructure:",squash"`
	ImageURL        string
	WebURL          string
	GenreFreeFormat string
}

func NewCircleDetailFromMap(m map[string]string) (*CircleDetail, error) {
	var circleDetail CircleDetail
	err := mapstructure.Decode(m, &circleDetail)
	if err != nil {
		return nil, errors.Wrap(err, "field to decode circle detail map to CircleDetail struct")
	}
	return &circleDetail, nil
}

func lineToMap(headers, line []string) (m map[string]string, err error) {
	if len(headers) != len(line) {
		return nil, errors.New(fmt.Sprintf("headers and line length are must be same (len(headers):%d, len(line):%d)", len(headers), len(line)))
	}

	m = map[string]string{}
	for i, v := range line {
		m[headers[i]] = v
	}
	return m, nil
}

func LineToCircleDetail(headers, line []string) (*CircleDetail, error) {
	m, err := lineToMap(headers, line)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert line to map")
	}
	return NewCircleDetailFromMap(m)
}

func circleDetailToMap(circleDetail *CircleDetail) map[string]string {
	m := structs.Map(circleDetail)
	m2 := map[string]string{}
	for k, v := range m {
		m2[k] = fmt.Sprint(v)
	}

	return m2
}

func CircleDetailToHeaders(circleDetail *CircleDetail) (headers []string) {
	m := circleDetailToMap(circleDetail)
	for header := range m {
		headers = append(headers, header)
	}
	return
}

func mapToLine(headers []string, m map[string]string) (line []string, err error) {
	for _, header := range headers {
		v, ok := m[header]
		if !ok {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to convert circle detail map to line, %s not found in map", header))
		}
		line = append(line, v)
	}
	return
}

func CircleDetailToLine(headers []string, circleDetail *CircleDetail) ([]string, error) {
	m := circleDetailToMap(circleDetail)
	line, err := mapToLine(headers, m)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert circle detail struct to line")
	}
	return line, nil
}

type Source struct {
	Url      string
	FileName string
}

func NewSource(sourcePath string) *Source {
	url, ok := GetCSVURL(sourcePath)
	fileName := path.Base(sourcePath)
	if ok {
		fileName = path.Base(url)
	}

	return &Source{
		Url:      url,
		FileName: fileName,
	}
}

func GetCSVURL(source string) (string, bool) {
	if strings.Contains(source, "http") {
		return source, true
	}
	u, ok := URLMap[source]
	if ok {
		return u, true
	}
	return "", false
}
