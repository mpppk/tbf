package util

import (
	"encoding/csv"
	"os"
	"strings"

	"fmt"

	"bufio"

	"github.com/fatih/structs"
	"github.com/mpppk/tbf/crawl"
	"github.com/pkg/errors"
)

type CircleCSV struct {
	filePath string
	file     *os.File
	writer   *csv.Writer
	headers  []string
}

func NewCircleCSV(filePath string) (*CircleCSV, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open circle csv: "+filePath)
	}

	writer := csv.NewWriter(file)

	return &CircleCSV{
		filePath: filePath,
		file:     file,
		writer:   writer,
	}, nil
}

func (c *CircleCSV) getHeaders() ([]string, bool) {
	if c.headers != nil && len(c.headers) > 0 {
		return c.headers, true
	}

	scanner := bufio.NewScanner(c.file)
	for scanner.Scan() {
		line := scanner.Text()
		c.headers = strings.Split(line, ",") // FIXME
		return c.headers, true
	}
	return nil, false
}

func (c *CircleCSV) AppendCircleDetail(circleDetail *crawl.CircleDetail) error {
	headers, ok := c.getHeaders()
	if !ok {
		headers = circleDetailToHeaders(circleDetail)
		c.headers = headers       // FIXME
		c.writer.Write(c.headers) // FIXME
	}

	line, err := circleDetailToLine(headers, circleDetail)

	fmt.Println(line)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to append circle detail struct as line to circle csv: %#v", circleDetail))
	}
	if err := c.writer.Write(line); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write circle detail struct as line to csv: %#v", circleDetail))
	}
	return nil
}

func (c *CircleCSV) ToCircleDetailMap() (m map[string]*crawl.CircleDetail, err error) {
	m = map[string]*crawl.CircleDetail{}
	reader := csv.NewReader(c.file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to read circle csv from %s", c.filePath))
	}

	if len(lines) < 1 {
		return m, nil
	}

	headers := lines[0]

	for _, line := range lines[1:] {
		circleDetail, err := lineToCircleDetail(headers, line)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert csv line to CircleDetail struct")
		}
		m[circleDetail.Space] = circleDetail
	}
	return m, nil
}

func (c *CircleCSV) Flush() {
	c.writer.Flush()
}

func (c *CircleCSV) Close() error {
	if err := c.file.Close(); err != nil {
		return errors.Wrap(err, "failed to close circle csv: "+c.filePath)
	}
	return nil
}

func isExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func lineToMap(headers, line []string) (m map[string]string) {
	m = map[string]string{}
	for i, v := range line {
		m[headers[i]] = v
	}
	return m
}

func lineToCircleDetail(headers, line []string) (*crawl.CircleDetail, error) {
	m := lineToMap(headers, line)
	return crawl.NewCircleDetailFromMap(m)
}

func circleDetailToMap(circleDetail *crawl.CircleDetail) map[string]string {
	m := structs.Map(circleDetail)
	m2 := map[string]string{}
	for k, v := range m {
		m2[k] = fmt.Sprint(v)
	}

	return m2
}

func circleDetailToHeaders(circleDetail *crawl.CircleDetail) (headers []string) {
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

func circleDetailToLine(headers []string, circleDetail *crawl.CircleDetail) ([]string, error) {
	m := circleDetailToMap(circleDetail)
	line, err := mapToLine(headers, m)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert circle detail struct to line")
	}
	return line, nil
}
