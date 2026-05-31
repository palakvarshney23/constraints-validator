package main

import "testing"

func TestGitHubPattern(t *testing.T) {
	cases := []struct {
		in string
		ok bool
	}{
		{"https://github.com/golang/go", true},
		{"https://github.com/a/b.git", true},
		{"https://gitlab.com/a/b", false},
		{"not-a-url", false},
		{"", false},
	}
	for _, c := range cases {
		if got := ghPat.MatchString(c.in); got != c.ok {
			t.Errorf("ghPat(%q)=%v want %v", c.in, got, c.ok)
		}
	}
}

func TestExtensionPattern(t *testing.T) {
	for _, p := range []string{"main.go", "app.py", "lib.rs", "script.sh"} {
		if !srcExt.MatchString(p) {
			t.Errorf("srcExt should match %q", p)
		}
	}
	for _, p := range []string{"readme.md", "style.css", "notes.txt"} {
		if srcExt.MatchString(p) {
			t.Errorf("srcExt should NOT match %q", p)
		}
	}
}

func TestCountLines(t *testing.T) {
	code := "package main\n// c\nimport \"fmt\"\n/* x */\nfunc main(){\n}\n"
	if n := countLines(code); n != 4 {
		t.Errorf("countLines=%d want 4", n)
	}
	block := "a\n/* start\n middle\n end */\nb\n"
	if n := countLines(block); n != 2 {
		t.Errorf("block comment countLines=%d want 2", n)
	}
}

func TestConstants(t *testing.T) {
	if maxFiles != 80 || maxBlob != 120000 {
		t.Error("unexpected constant value")
	}
	if len(extMap) == 0 || len(skip) == 0 {
		t.Error("maps should not be empty")
	}
}
