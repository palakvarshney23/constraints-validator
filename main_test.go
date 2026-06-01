package main

import (
	"regexp"
	"strings"
	"testing"
)

func TestGitHubPattern(t *testing.T) {
	gh := regexp.MustCompile(`(?i)github\.com[/:]([^/\s#?]+)/([^/\s#?]+)`)
	cases := []struct {
		in string
		ok bool
	}{
		{"https://github.com/golang/go", true},
		{"https://github.com/a/b.git", true},
		{"https://gitlab.com/a/b", false},
		{"not-a-url", false},
	}
	for _, c := range cases {
		if got := gh.MatchString(c.in); got != c.ok {
			t.Errorf("github URL %q: got %v want %v", c.in, got, c.ok)
		}
	}
}

func TestExtensionPattern(t *testing.T) {
	ext := regexp.MustCompile(`\.(go|py|js|ts|jsx|tsx|java|c|cpp|cc|h|cs|rb|php|rs|sh|bash)$`)
	for _, p := range []string{"main.go", "app.py", "lib.rs", "script.sh"} {
		if !ext.MatchString(p) {
			t.Errorf("should match %q", p)
		}
	}
	for _, p := range []string{"readme.md", "style.css"} {
		if ext.MatchString(p) {
			t.Errorf("should not match %q", p)
		}
	}
}

func testCountLines(code string) int {
	n, inBlock := 0, false
	for _, ln := range strings.Split(strings.ReplaceAll(code, "\r\n", "\n"), "\n") {
		t := strings.TrimSpace(ln)
		if t == "" {
			continue
		}
		if inBlock {
			if strings.Contains(t, "*/") {
				inBlock = false
			}
			continue
		}
		if strings.HasPrefix(t, "/*") {
			if strings.Index(t, "*/") < 0 {
				inBlock = true
			}
			continue
		}
		if strings.HasPrefix(t, "//") || strings.HasPrefix(t, "#") || strings.HasPrefix(t, "--") ||
			strings.HasPrefix(t, "<!--") || strings.HasPrefix(t, "*") || strings.HasPrefix(t, ";") {
			continue
		}
		n++
	}
	return n
}

func TestCountableLines(t *testing.T) {
	code := "package main\n// c\nimport \"fmt\"\n/* x */\nfunc main(){\n}\n"
	if n := testCountLines(code); n != 4 {
		t.Errorf("countable lines = %d, want 4", n)
	}
	block := "a\n/* start\n middle\n end */\nb\n"
	if n := testCountLines(block); n != 2 {
		t.Errorf("block comment lines = %d, want 2", n)
	}
}
