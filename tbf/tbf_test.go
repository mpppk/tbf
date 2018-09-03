package tbf_test

import (
	"reflect"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/mpppk/tbf/tbf"
)

var methodName = "LineToCircleDetail"
var dummyDetailURL = "dummyDetailURL"
var dummySpace = "dummySpace"
var dummyName = "dummyName"
var dummyPenname = "dummyPenname"
var dummyGenre = "dummyGenre"
var dummyImageURL = "dummyImageURL"
var dummyWebURL = "dummyWebURL"
var dummyGenreFreeFormat = "dummyGenreFreeFormat"

func generateCircleDetailHeaders() []string {
	return []string{
		"DetailURL",
		"Space",
		"Name",
		"Penname",
		"Genre",
		"ImageURL",
		"WebURL",
		"GenreFreeFormat",
	}
}

func generateDummyLine() []string {
	return []string{
		dummyDetailURL,
		dummySpace,
		dummyName,
		dummyPenname,
		dummyGenre,
		dummyImageURL,
		dummyWebURL,
		dummyGenreFreeFormat,
	}
}

func generateDummyCircle() *tbf.Circle {
	return &tbf.Circle{
		DetailURL: dummyDetailURL,
		Space:     dummySpace,
		Name:      dummyName,
		Penname:   dummyPenname,
		Genre:     dummyGenre,
	}
}

func generateDummyCircleDetail() *tbf.CircleDetail {
	return &tbf.CircleDetail{
		Circle:          *generateDummyCircle(),
		ImageURL:        dummyImageURL,
		WebURL:          dummyWebURL,
		GenreFreeFormat: dummyGenreFreeFormat,
	}
}

func generateDummyCircleDetailMap() map[string]string {
	return map[string]string{
		"DetailURL":       dummyDetailURL,
		"Space":           dummySpace,
		"Name":            dummyName,
		"Penname":         dummyPenname,
		"Genre":           dummyGenre,
		"ImageURL":        dummyImageURL,
		"WebURL":          dummyWebURL,
		"GenreFreeFormat": dummyGenreFreeFormat,
	}
}

func TestLineToCircleDetail(t *testing.T) {
	circleDetailHeaders := generateCircleDetailHeaders()
	dummyLine := generateDummyLine()

	tooFewLine := make([]string, len(dummyLine))
	copy(tooFewLine, dummyLine)
	tooFewLine = tooFewLine[1:]

	tooMuchLine := make([]string, len(dummyLine))
	copy(tooMuchLine, dummyLine)
	tooMuchLine = append(tooMuchLine, "dummy elm")

	dummyCircleDetail := generateDummyCircleDetail()

	cases := []struct {
		headers     []string
		line        []string
		expected    *tbf.CircleDetail
		willBeError bool
	}{
		{
			headers:     circleDetailHeaders,
			line:        dummyLine,
			expected:    dummyCircleDetail,
			willBeError: false,
		},
		{
			headers:     circleDetailHeaders,
			line:        tooFewLine,
			expected:    dummyCircleDetail,
			willBeError: true,
		},
		{
			headers:     circleDetailHeaders,
			line:        tooMuchLine,
			expected:    dummyCircleDetail,
			willBeError: true,
		},
	}

	for _, c := range cases {
		actual, err := tbf.LineToCircleDetail(c.headers, c.line)
		if err != nil && !c.willBeError {
			t.Fatalf("Unexpected error occured in : %s", err)
		}

		if c.willBeError {
			if err == nil {
				t.Fatalf("%s is expected to be error if headers %q and line %q are given.", methodName, c.headers, c.line)
			} else {
				continue
			}
		}

		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("%s is expected to return %#v when {headers: %q, line: %q} are given ,but actually return %v\ndiff: %s", methodName, c.expected, c.headers, c.line, actual, pretty.Compare(actual, c.expected))

		}
	}
}

func TestNewCircleDetailFromMap(t *testing.T) {
	methodName := "TestNewCircleDetailFromMap"
	cases := []struct {
		m           map[string]string
		expected    *tbf.CircleDetail
		willBeError bool
	}{
		{
			m:           generateDummyCircleDetailMap(),
			expected:    generateDummyCircleDetail(),
			willBeError: false,
		},
	}

	for _, c := range cases {
		actual, err := tbf.NewCircleDetailFromMap(c.m)
		if err != nil && !c.willBeError {
			t.Fatalf("Unexpected error occured in : %s", err)
		}

		if c.willBeError {
			if err == nil {
				t.Fatalf("%s is expected to be error if map %v is given.", methodName, c.m)
			} else {
				continue
			}
		}

		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("%s is expected to return %#v when map %v are given ,but actually return %v\ndiff: %s", methodName, c.expected, c.m, actual, pretty.Compare(actual, c.expected))

		}
	}
}

func TestCircleDetailToLine(t *testing.T) {
	methodName := "TestCircleDetailToLine"
	cases := []struct {
		headers      []string
		circleDetail *tbf.CircleDetail
		expected     []string
		willBeError  bool
	}{
		{
			headers:      generateCircleDetailHeaders(),
			circleDetail: generateDummyCircleDetail(),
			expected:     generateDummyLine(),
			willBeError:  false,
		},
	}

	for _, c := range cases {
		actual, err := tbf.CircleDetailToLine(c.headers, c.circleDetail)
		if err != nil && !c.willBeError {
			t.Errorf("Unexpected error occured in : %s", err)
			continue
		}

		if c.willBeError {
			if err == nil {
				t.Errorf("%s is expected to be error if headers %v and circleDetail %#v are given.",
					methodName,
					c.headers,
					c.circleDetail)
			}
			continue
		}

		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("%s is expected to return %#v when headers %v and circleDetail %#v are given,"+
				" but actually return %v\ndiff: %s",
				methodName,
				c.expected,
				c.headers,
				c.circleDetail,
				actual,
				pretty.Compare(actual, c.expected))
		}
	}
}
