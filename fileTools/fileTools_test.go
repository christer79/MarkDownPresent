package fileTools

import "testing"

func TestFindFiles(t *testing.T) {
	cases := []struct {
		in    string
		match string
		want  int
	}{
		{"test/1", "*.md", 2},
		{"test/1", "*.txt", 1},
		{"test/2", "*.txt", 1},
		{"test/3", "apa?.h", 2},
	}
	for _, c := range cases {
		got := FindFiles(c.in, c.match)

		if len(got) != c.want {
			t.Errorf("FindeFiles(%q) == %d, want %d", c.in, len(got), c.want)
		}
	}
}
