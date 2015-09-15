package mdFormat

import "testing"

func TestIsCommentedLine(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"APA", false},
		{"[[//]]:", true},
		{"# ", false},
		{"skadlk [[//]]:", true},
		{"[[//]]: # kdjaslk", true},
	}
	for _, c := range cases {
		got := IsCommentedLine(c.in)
		if got != c.want {
			t.Errorf("IsCommentedLine(%q) == %d, want %d", c.in, got, c.want)
		}
	}
}

func TestExtractCommentedDataField(t *testing.T) {
	cases := []struct {
		in, label, want string
	}{
		{"[[//]]: # Century: \"Test\"", "Century", "Test"}, // With comment
		{"Century: \"Test\"", "Century", "Test"},           //Beginning of line
		{"Background: \"Test\"", "Century", ""},            // Not found
		{"Century: Test", "Century", ""},                   //No comment
	}
	for _, c := range cases {
		got := ExtractCommentDataFiled(c.in, c.label)
		if got != c.want {
			t.Errorf("ExtractCommentDataFiled(%q, %q) == %q, want %q", c.in, c.label, got, c.want)
		}
	}
}

func TestMarkdownToHTML(t *testing.T) {
	t.Log("No tests are run")
}
