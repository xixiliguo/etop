package util

import (
	"testing"
	"time"
)

func TestGetHumanSize(t *testing.T) {
	testCases := []struct {
		size     int
		expected string
	}{
		{579, "579.0 B"},
		{1024, "1.0 KB"},
		{4125, "4.0 KB"},
		{103000, "100.6 KB"},
		{10300000, "9.8 MB"},
		{1030000000000, "959.3 GB"},
		{1030000000000000, "936.8 TB"},
		{1030000000000000000, "914.8 PB"},
		{9030000000000000000, "8020.3 PB"},
	}
	for _, testCase := range testCases {
		if actual := GetHumanSize(testCase.size); actual != testCase.expected {
			t.Errorf("expected '%v' but got '%v'", testCase.expected, actual)
		}
	}
}

func BenchmarkGetHumanSize(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		GetHumanSize(579)
	}

}

func TestConvertToTime(t *testing.T) {
	y, m, d := time.Now().Date()
	testCases := []struct {
		intput      string
		shouldError bool
		expected    int64
	}{
		{
			intput:      "1hago",
			shouldError: false,
			expected:    time.Now().Add(-time.Hour).Unix(),
		},
		{
			intput:      "2h ago",
			shouldError: false,
			expected:    time.Now().Add(-2 * time.Hour).Unix(),
		},
		{
			intput:      "3h  ago",
			shouldError: false,
			expected:    time.Now().Add(-3 * time.Hour).Unix(),
		},
		{
			intput:      "4h ago ",
			shouldError: true,
			expected:    0,
		},
		{
			intput:      "123hago",
			shouldError: false,
			expected:    time.Now().Add(-123 * time.Hour).Unix(),
		},
		{
			intput:      "5d3h1sago",
			shouldError: true,
			expected:    time.Now().Add(-123*time.Hour - 1*time.Second).Unix(),
		},
		{
			intput:      "21:45",
			shouldError: false,
			expected:    time.Date(y, m, d, 21, 45, 0, 0, time.Local).Unix(),
		},
		{
			intput:      "45",
			shouldError: true,
			expected:    0,
		},
		{
			intput:      "2006-01-02 15:04",
			shouldError: false,
			expected:    time.Date(2006, 1, 2, 15, 4, 0, 0, time.Local).Unix(),
		},
		{
			intput:      "2006-01-02 15:04:05",
			shouldError: true,
			expected:    0,
		},
		{
			intput:      "1970-01-01 22:23",
			shouldError: false,
			expected:    time.Date(1970, 1, 1, 22, 23, 0, 0, time.Local).Unix(),
		},
		{
			intput:      "01-01 22:23",
			shouldError: false,
			expected:    time.Date(time.Now().Year(), 1, 1, 22, 23, 0, 0, time.Local).Unix(),
		},
	}

	for _, testCase := range testCases {
		actual, err := ConvertToTime(testCase.intput)
		if err != nil && testCase.shouldError == false {
			t.Errorf("input: %q expected no error but got %s", testCase.intput, err)
			continue
		}
		if err == nil && testCase.shouldError == true {
			t.Errorf("input: %q expected error but got no", testCase.intput)
			continue
		}
		if testCase.shouldError == false && actual != testCase.expected {
			t.Errorf("input: %q expected %d but got %d", testCase.intput, testCase.expected, actual)
			continue
		}

	}
}
