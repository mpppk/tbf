package csv

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"net/http"

	"io"

	"hash/crc32"
	"io/ioutil"

	"encoding/json"

	"github.com/mpppk/tbf/tbf"
	"github.com/pkg/errors"
)

type CircleCSV struct {
	filePath string
	headers  []string
}

type Meta struct {
	Checksum uint32 `json:"checksum"`
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

func IsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func DownloadCSVIfChanged(csvURL, csvMetaURL, filePath string) (bool, error) {
	if !IsExist(filePath) {
		return downloadCSV(csvURL, filePath)
	}

	meta, err := readCSVMetaFromHTTP(csvMetaURL)
	if err != nil {
		return false, errors.Wrap(
			err, fmt.Sprintf("failed to download csv meta data from %s", csvMetaURL))
	}

	checksum, err := getFileCheckSum(filePath)
	if err != nil {
		return false, errors.Wrap(err, "failed to read csv file")
	}

	if checksum == meta.Checksum {
		return false, nil
	}

	fmt.Fprintf(
		os.Stderr,
		"csv file will be downloaded becase checksums are different between meta(%v) and local file(%v)\n",
		meta.Checksum,
		checksum)
	return downloadCSV(csvURL, filePath)
}

func readCSVMetaFromHTTP(csvMetaURL string) (*Meta, error) {
	res, err := http.Get(csvMetaURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get csv from URL: "+csvMetaURL)
	}

	if res.StatusCode != 200 {
		return nil, errors.New(
			fmt.Sprintf("failed to fetch csv from %s: invalid statuscode: %v", csvMetaURL, res.Status))
	}

	csvMetaJsonBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read http response")
	}

	meta := &Meta{}
	if err = json.Unmarshal(csvMetaJsonBytes, meta); err != nil {
		return nil, errors.Wrap(err,
			fmt.Sprintf(
				"failed to unmarshal latest csv json from %s, contents: %s",
				csvMetaURL,
				string(csvMetaJsonBytes)))
	}
	return meta, nil
}

func downloadCSV(csvURL, filePath string) (bool, error) {
	res, err := http.Get(csvURL)
	if err != nil {
		return false, errors.Wrap(err, "failed to download CSV from "+csvURL)
	}

	if res.StatusCode != 200 {
		return false, errors.New(
			fmt.Sprintf("failed to fetch csv from %s: %v", csvURL, res.Status))
	}

	file, err := os.Create(filePath)
	if err != nil {
		return false, errors.Wrap(err, "failed to create csv file to "+filePath)
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return false, errors.Wrap(err, "failed to write to downloaded csv to "+filePath)
	}

	return true, nil
}

func getFileCheckSum(filePath string) (uint32, error) {
	if !IsExist(filePath) {
		return 0, errors.New("csv file not found: " + filePath)
	}

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, errors.Wrap(err, "failed to read csv file")
	}

	return crc32.ChecksumIEEE(contents), nil
}
