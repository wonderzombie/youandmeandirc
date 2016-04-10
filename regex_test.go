package youandmeandirc

import (
	"reflect"
	"testing"
)

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
			nil,
		},
		{
			"no trailing slash here s/foo/bar",
			nil,
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
		// {
		// 	"now s/foo and/bar and/ is the worst regex",
		// 	&Replacement{"foo and", "bar and"},
		// },
		{
			"",
			nil,
		},
	}

	for _, test := range tests {
		got := regex(test.message)
		want := test.want
		if got != nil && test.want == nil || !reflect.DeepEqual(got, want) {
			t.Errorf("regex(%q) => %q; want %q", test.message, got, want)
		}
	}
}

// func TestRewrite(t *testing.T) {
// 	tests := []struct {
// 		orig string
// 		r    string
// 		want string
// 	}{
// 		{"hello wrold", "s/wrold/world", "hello world"},
// 	}
// }
