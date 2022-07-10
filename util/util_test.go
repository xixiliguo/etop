package util

import "testing"

// assertEquals(t, "1kB", HumanSize(1000))
// assertEquals(t, "1.024kB", HumanSize(1024))
// assertEquals(t, "1MB", HumanSize(1000000))
// assertEquals(t, "1.049MB", HumanSize(1048576))
// assertEquals(t, "2MB", HumanSize(2*MB))

func TestGetHumanSize(t *testing.T) {
	testCases := []struct {
		size     int
		expected string
	}{
		{579, "579B"},
		{1024, "1KB"},
		{4125, "4KB"},
		{103000, "100KB"},
		{10300000, "9MB"},
		{1030000000000, "959GB"},
		{1030000000000000, "936TB"},
		{1030000000000000000, "914PB"},
		{9030000000000000000, "8020PB"},
	}
	for _, testCase := range testCases {
		if actual := GetHumanSize(testCase.size); actual != testCase.expected {
			t.Errorf("Expected '%v' but got '%v'", testCase.expected, actual)
		}
	}
}
