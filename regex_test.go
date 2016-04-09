package youandmeandirc

import "testing"

func TestHasRegex(t *testing.T) {
	tests := []struct {
		message string
		want    *RegexResult
	}{
		{
			"whoops I meant s/foo/bar/",
			&RegexResult{"foo", "bar"},
		},
		{
			"this regex here /foo/bar/ has no leading s",
			&RegexResult{},
		},
		{
			"no trailing slash here s/foo/bar",
			&RegexResult{},
		},
		{
			"this will fail probably s/foo// but who knows",
			&RegexResult{"foo", ""},
		},
	}

	for _, test := range tests {
		got := hasRegex(test.message)
		want := test.want
		if got.search != want.search && got.replace != want.replace {
			t.Errorf("hasRegex(%q) => %s; want %s", test.message, got, want)
		}
	}
}
