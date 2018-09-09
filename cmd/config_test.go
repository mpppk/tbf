package cmd

import (
	"testing"

	"github.com/mpppk/tbf/csv"
)

func TestCreateConfigFromSource(t *testing.T) {
	cases := []struct {
		source   string
		expected *listConfig
	}{
		{
			source: "tbf4",
			expected: &listConfig{
				url:      csv.URLMap["tbf4"],
				fileName: "tbf4_circles.csv",
			},
		},
		{
			source: "latest",
			expected: &listConfig{
				url:      csv.URLMap["latest"],
				fileName: "latest_circles.csv",
			},
		},
		{
			source: "http://example.com/test_circles.csv",
			expected: &listConfig{
				url:      "http://example.com/test_circles.csv",
				fileName: "test_circles.csv",
			},
		},
		{
			source: "test_circles.csv",
			expected: &listConfig{
				url:      "",
				fileName: "test_circles.csv",
			},
		},
	}

	for _, c := range cases {
		config := createConfigFromSource(c.source)

		if config.url != c.expected.url {
			t.Errorf("CreateConfigFromSource is expected to return url: %q when source %q is given, but actually has %q",
				c.expected, c.source, config)
		}

		if config.fileName != c.expected.fileName {
			t.Errorf("CreateConfigFromSource is expected to return fileName: %q when source %q is given, but actually has %q",
				c.expected, c.source, config)
		}
	}
}
