package tbf

import (
	"testing"
)

func TestNewSource(t *testing.T) {
	cases := []struct {
		source   string
		expected *Source
	}{
		{
			source: "tbf4",
			expected: &Source{
				Url:      URLMap["tbf4"],
				FileName: "tbf4_circles.csv",
			},
		},
		{
			source: "latest",
			expected: &Source{
				Url:      URLMap["latest"],
				FileName: "latest_circles.csv",
			},
		},
		{
			source: "http://example.com/test_circles.csv",
			expected: &Source{
				Url:      "http://example.com/test_circles.csv",
				FileName: "test_circles.csv",
			},
		},
		{
			source: "test_circles.csv",
			expected: &Source{
				Url:      "",
				FileName: "test_circles.csv",
			},
		},
	}

	for _, c := range cases {
		config := NewSource(c.source)

		if config.Url != c.expected.Url {
			t.Errorf("NewSource is expected to return Url: %q when source %q is given, but actually has %q",
				c.expected, c.source, config)
		}

		if config.FileName != c.expected.FileName {
			t.Errorf("NewSource is expected to return FleName: %q when source %q is given, but actually has %q",
				c.expected, c.source, config)
		}
	}
}
