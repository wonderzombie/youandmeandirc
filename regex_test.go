package youandmeandirc

import "testing"

func TestHasRegex(t *testing.T) {
	tests := []struct {
		message string
		want    *Replacement
	}{
		{
			"whoops I meant s/foo/bar/",
			&Replacement{"foo", "bar"},
		},
		{
			"this regex here /foo/bar/ has no leading s",
			&Replacement{},
		},
		{
			"no trailing slash here s/foo/bar",
			&Replacement{},
		},
		{
			"this will fail probably s/foo// but who knows",
			&Replacement{"foo", ""},
		},
		{
			"this one is misleading/confusing/evil s/evil/good/",
			&Replacement{"evil", "good"},
		},
		{
			"now s/foo/bar/ there are two s/evil/good/ regexen",
			&Replacement{"foo", "bar"},
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
