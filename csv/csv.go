package csv

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/mpppk/tbf/tbf"
	"github.com/pkg/errors"
)

type CircleCSV struct {
	filePath string
	headers  []string
}

func NewCircleCSV(filePath string) (*CircleCSV, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open circle csv: "+filePath)
	}

	var headers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		headers = strings.Split(line, ",") // FIXME
		break
	}

	return &CircleCSV{
		filePath: filePath,
		headers:  headers,
	}, nil
}

func (c *CircleCSV) getHeaders() ([]string, bool) {
	if c.headers != nil && len(c.headers) > 0 {
		return c.headers, true
	}
	return nil, false
}

func (c *CircleCSV) AppendLine(line []string) error {
	file, err := os.OpenFile(c.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		return errors.Wrap(err, "failed to open circle csv: "+c.filePath)
	}
	writer := csv.NewWriter(file)
	if err := writer.Write(line); err != nil {
		return errors.Wrap(err, "failed to write to circle csv: "+c.filePath)
	}
	writer.Flush()
	return nil
}

func (c *CircleCSV) AppendCircleDetail(circleDetail *tbf.CircleDetail) error {
	headers, ok := c.getHeaders()
	if !ok {
		headers = tbf.CircleDetailToHeaders(circleDetail)
		c.headers = headers // FIXME
		if err := c.AppendLine(headers); err != nil {
			return errors.Wrap(err, "failed to append headers to "+c.filePath)
		}
	}

	line, err := tbf.CircleDetailToLine(headers, circleDetail)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to append circle detail struct as line to circle csv: %#v", circleDetail))
	}
	if err := c.AppendLine(line); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write circle detail struct as line to csv: %#v", circleDetail))
	}
	return nil
}

func (c *CircleCSV) ToCircleDetailMap() (m map[string]*tbf.CircleDetail, err error) {
	m = map[string]*tbf.CircleDetail{}
	file, err := os.OpenFile(c.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to open file: %s", c.filePath))
	}
	defer file.Close()
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to read circle csv from %s", c.filePath))
	}

	if len(lines) < 1 {
		return m, nil
	}

	headers := lines[0]

	for _, line := range lines[1:] {
		circleDetail, err := tbf.LineToCircleDetail(headers, line)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert csv line to CircleDetail struct")
		}
		m[circleDetail.Space] = circleDetail
	}
	return m, nil
}
