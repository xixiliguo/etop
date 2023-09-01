package model

import "testing"

func TestRender(t *testing.T) {

	testCases := []struct {
		name     string
		f        Field
		input    any
		expected string
	}{
		{"Render int", Field{"Load1", Raw, 0, "", 10, false}, 1, "1"},
		{"Render float", Field{"User", Raw, 1, "%", 10, false}, 1.27, "1.3%"},
		{"Render int with suffix", Field{"MemTotal", Raw, 0, " KB", 10, false}, 1024, "1024 KB"},
		{"Render int with suffix and width", Field{"MemTotal", Raw, 0, " KB", 10, true}, 1024, "1024 KB   "},
		{"Render readablesize", Field{"ReadByte/s", HumanReadableSize, 1, "/s", 10, false}, 3457, "3.4KB/s"},
		{"Render string", Field{"Comm", Raw, 0, "", 16, false}, "systemd", "systemd"},
		{"Render short string", Field{"Comm", Raw, 0, "", 16, true}, "sh", "sh              "},
	}
	for _, testCase := range testCases {
		if actual := testCase.f.Render(testCase.input); actual != testCase.expected {
			t.Errorf("%s: expected '%v' but got '%v'", testCase.name, testCase.expected, actual)
		}
	}
}

func BenchmarkRender(b *testing.B) {
	testCases := []struct {
		name     string
		f        Field
		input    any
		expected string
	}{
		{"Render int", Field{"Load1", Raw, 0, "", 10, false}, 1, "1"},
		{"Render float", Field{"User", Raw, 1, "%", 10, false}, 1.27, "1.3%"},
		{"Render int with suffix", Field{"MemTotal", Raw, 0, " KB", 10, false}, 1024, "1024 KB"},
		{"Render int with suffix and width", Field{"MemTotal", Raw, 0, " KB", 10, true}, 1024, "1024 KB   "},
		{"Render readablesize", Field{"ReadByte/s", HumanReadableSize, 1, "/s", 10, false}, 3457, "3.4KB/s"},
		{"Render string", Field{"Comm", Raw, 0, "", 16, false}, "systemd", "systemd"},
	}

	for _, testCase := range testCases {

		b.Run(testCase.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				testCase.f.Render(testCase.input)
			}
		})
	}
}
