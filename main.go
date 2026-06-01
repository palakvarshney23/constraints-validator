// Code Olympics Constraint Auditor | Single-Function Master | Feature-Rich Dev (400) | Go
package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	for _, l := range []string{
		"  ╔════════════════════════════════════════════════════════╗",
		"  ║  ┌────────────────────────────────────────────────┐  ║",
		"  ║  │  ◉  CODE OLYMPICS 2026  ·  co-check            │  ║",
		"  ║  │      GitHub 4D Constraint Auditor              │  ║",
		"  ║  │  D1 Core │ D2 Lines │ D3 Domain │ D4 Language  │  ║",
		"  ║  └────────────────────────────────────────────────┘  ║",
		"  ╚════════════════════════════════════════════════════════╝",
	} {
		fmt.Println(l)
	}
	if len(os.Args) > 1 {
		a := os.Args[1]
		if a == "-h" || a == "--help" {
			fmt.Println("\n  Usage:  co-check [github-url]")
			fmt.Println("  Env:    GITHUB_TOKEN  (optional, avoids API rate limits)")
			fmt.Println("  Output: co-check-report.txt")
			return
		}
		if err := analyze(a); err != nil {
			fmt.Println("\n  Error:", err)
		}
		return
	}
	fmt.Print("\n  GitHub URL: ")
	sc := bufio.NewScanner(os.Stdin)
	if !sc.Scan() {
		return
	}
	u := strings.TrimSpace(sc.Text())
	if u == "" {
		return
	}
	if err := analyze(u); err != nil {
		fmt.Println("\n  Error:", err)
	}
}
func analyze(rawURL string) error {
	cl := &http.Client{Timeout: 25 * time.Second}
	tok := os.Getenv("GITHUB_TOKEN")
	spin := []string{"⠋", "⠙", "⠹", "⠸"}
	srcExt := regexp.MustCompile(`\.(go|py|js|ts|jsx|tsx|java|c|cpp|cc|h|cs|rb|php|rs|sh|bash)$`)
	ghPat := regexp.MustCompile(`(?i)github\.com[/:]([^/\s#?]+)/([^/\s#?]+)`)
	m := ghPat.FindStringSubmatch(rawURL)
	if m == nil {
		return fmt.Errorf("invalid GitHub URL")
	}
	owner, repo := m[1], strings.TrimSuffix(m[2], ".git")
	var rep strings.Builder
	for _, l := range []string{
		"  ╔════════════════════════════════════════════════════════╗",
		"  ║  ┌────────────────────────────────────────────────┐  ║",
		"  ║  │  ◉  CODE OLYMPICS 2026  ·  co-check            │  ║",
		"  ║  │      GitHub 4D Constraint Auditor              │  ║",
		"  ║  │  D1 Core │ D2 Lines │ D3 Domain │ D4 Language  │  ║",
		"  ║  └────────────────────────────────────────────────┘  ║",
		"  ╚════════════════════════════════════════════════════════╝",
	} {
		rep.WriteString(l + "\n")
	}
	fmt.Fprintf(&rep, "\n  Target: %s/%s\n\n", owner, repo)
	lbl := "Connecting to GitHub..."
	for i := 0; i < 4; i++ {
		fmt.Printf("\r  %s  %s", spin[i%4], lbl)
		time.Sleep(40 * time.Millisecond)
	}
	rq, _ := http.NewRequest("GET", "https://api.github.com/repos/"+owner+"/"+repo, nil)
	rq.Header.Set("Accept", "application/vnd.github+json")
	rq.Header.Set("User-Agent", "co-check")
	if tok != "" {
		rq.Header.Set("Authorization", "Bearer "+tok)
	}
	rs, err := cl.Do(rq)
	if err != nil {
		return err
	}
	repoB, _ := io.ReadAll(rs.Body)
	rs.Body.Close()
	if rs.StatusCode >= 400 {
		return fmt.Errorf("GitHub API %s", rs.Status)
	}
	var meta struct {
		DefaultBranch string   `json:"default_branch"`
		Language      string   `json:"language"`
		Description   string   `json:"description"`
		Topics        []string `json:"topics"`
	}
	if json.Unmarshal(repoB, &meta) != nil {
		return fmt.Errorf("bad repo JSON")
	}
	msg := fmt.Sprintf("  ✓  %s\n", lbl)
	fmt.Print(msg)
	rep.WriteString(msg)
	branch := meta.DefaultBranch
	if branch == "" {
		branch = "main"
	}
	lbl = "Resolving branch..."
	for i := 0; i < 4; i++ {
		fmt.Printf("\r  %s  %s", spin[i%4], lbl)
		time.Sleep(40 * time.Millisecond)
	}
	rq, _ = http.NewRequest("GET", "https://api.github.com/repos/"+owner+"/"+repo+"/branches/"+branch, nil)
	rq.Header.Set("Accept", "application/vnd.github+json")
	rq.Header.Set("User-Agent", "co-check")
	if tok != "" {
		rq.Header.Set("Authorization", "Bearer "+tok)
	}
	rs, err = cl.Do(rq)
	if err != nil {
		return err
	}
	brB, _ := io.ReadAll(rs.Body)
	rs.Body.Close()
	if rs.StatusCode >= 400 {
		return fmt.Errorf("branch API %s", rs.Status)
	}
	var br struct {
		Commit struct {
			Sha string `json:"sha"`
		} `json:"commit"`
	}
	if json.Unmarshal(brB, &br) != nil || br.Commit.Sha == "" {
		return fmt.Errorf("bad branch JSON")
	}
	msg = fmt.Sprintf("  ✓  Branch: %s\n", branch)
	fmt.Print(msg)
	rep.WriteString(msg)
	lbl = "Fetching file tree..."
	for i := 0; i < 4; i++ {
		fmt.Printf("\r  %s  %s", spin[i%4], lbl)
		time.Sleep(40 * time.Millisecond)
	}
	treeURL := "https://api.github.com/repos/" + owner + "/" + repo + "/git/trees/" + br.Commit.Sha + "?recursive=1"
	rq, _ = http.NewRequest("GET", treeURL, nil)
	rq.Header.Set("Accept", "application/vnd.github+json")
	rq.Header.Set("User-Agent", "co-check")
	if tok != "" {
		rq.Header.Set("Authorization", "Bearer "+tok)
	}
	rs, err = cl.Do(rq)
	if err != nil {
		return err
	}
	treeB, _ := io.ReadAll(rs.Body)
	rs.Body.Close()
	if rs.StatusCode >= 400 {
		return fmt.Errorf("tree API %s", rs.Status)
	}
	var tree struct {
		Tree []struct {
			Path string `json:"path"`
			Type string `json:"type"`
			Sha  string `json:"sha"`
			Size int    `json:"size"`
		} `json:"tree"`
	}
	if json.Unmarshal(treeB, &tree) != nil {
		return fmt.Errorf("bad tree JSON")
	}
	msg = fmt.Sprintf("  ✓  Tree: %d entries\n", len(tree.Tree))
	fmt.Print(msg)
	rep.WriteString(msg)
	lbl = "Downloading source + README..."
	for i := 0; i < 4; i++ {
		fmt.Printf("\r  %s  %s", spin[i%4], lbl)
		time.Sleep(40 * time.Millisecond)
	}
	var buf, pathLow strings.Builder
	fileN, skipN := 0, 0
	langCnt := map[string]int{}
	readmeLow := ""
	metaHigh := strings.ToLower(meta.Description + " " + strings.Join(meta.Topics, " ") + " " + owner + " " + repo)
	rq, _ = http.NewRequest("GET", "https://api.github.com/repos/"+owner+"/"+repo+"/readme", nil)
	rq.Header.Set("Accept", "application/vnd.github+json")
	rq.Header.Set("User-Agent", "co-check")
	if tok != "" {
		rq.Header.Set("Authorization", "Bearer "+tok)
	}
	if rs2, e := cl.Do(rq); e == nil && rs2.StatusCode < 400 {
		rb, _ := io.ReadAll(rs2.Body)
		rs2.Body.Close()
		var rd struct {
			Content string `json:"content"`
		}
		if json.Unmarshal(rb, &rd) == nil && rd.Content != "" {
			if b, e := base64.StdEncoding.DecodeString(strings.ReplaceAll(rd.Content, "\n", "")); e == nil {
				readmeLow = strings.ToLower(string(b))
			}
		}
	}
	for _, ent := range tree.Tree {
		pathLow.WriteString(strings.ToLower(ent.Path))
		pathLow.WriteByte(' ')
		if ent.Type != "blob" || !srcExt.MatchString(ent.Path) || ent.Size > 120000 {
			continue
		}
		skip := false
		for _, d := range []string{"node_modules/", "vendor/", ".git/", "dist/", "build/", "__pycache__/"} {
			if strings.Contains(ent.Path, d) {
				skip = true
				break
			}
		}
		if skip {
			skipN++
			continue
		}
		bq, e := http.NewRequest("GET", "https://api.github.com/repos/"+owner+"/"+repo+"/git/blobs/"+ent.Sha, nil)
		if e != nil {
			skipN++
			continue
		}
		bq.Header.Set("Accept", "application/vnd.github.raw")
		bq.Header.Set("User-Agent", "co-check")
		if tok != "" {
			bq.Header.Set("Authorization", "Bearer "+tok)
		}
		brs, e := cl.Do(bq)
		if e != nil || brs.StatusCode >= 400 {
			if brs != nil {
				brs.Body.Close()
			}
			skipN++
			continue
		}
		body, e := io.ReadAll(brs.Body)
		brs.Body.Close()
		if e != nil {
			skipN++
			continue
		}
		langCnt[ent.Path[strings.LastIndex(ent.Path, ".")+1:]]++
		buf.Write(body)
		buf.WriteByte('\n')
		fileN++
		if fileN >= 80 {
			break
		}
	}
	code := buf.String()
	if strings.TrimSpace(code) == "" {
		return fmt.Errorf("no source files found")
	}
	msg = fmt.Sprintf("  ✓  Files: %d analyzed, %d skipped\n", fileN, skipN)
	fmt.Print(msg)
	rep.WriteString(msg)
	lbl = "Running constraint analysis..."
	for i := 0; i < 4; i++ {
		fmt.Printf("\r  %s  %s", spin[i%4], lbl)
		time.Sleep(40 * time.Millisecond)
	}
	codeLines, inBlock := 0, false
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
			if idx := strings.Index(t, "*/"); idx < 0 || idx <= 2 {
				inBlock = true
			}
			continue
		}
		if strings.HasPrefix(t, "//") || strings.HasPrefix(t, "#") || strings.HasPrefix(t, "--") ||
			strings.HasPrefix(t, "<!--") || strings.HasPrefix(t, "*") || strings.HasPrefix(t, ";") {
			continue
		}
		codeLines++
	}
	pyImp := len(regexp.MustCompile(`(?m)^\s*(import\s+\S|from\s+\S+\s+import)`).FindAllString(code, -1))
	jsImp := len(regexp.MustCompile(`(?m)^\s*(import\s+["'{]|require\s*\()`).FindAllString(code, -1))
	javaImp := len(regexp.MustCompile(`(?m)^\s*import\s+[a-z]`).FindAllString(code, -1))
	cppImp := len(regexp.MustCompile(`(?m)^\s*#include\s*[<"]`).FindAllString(code, -1))
	rubyImp := len(regexp.MustCompile(`(?m)^\s*require\s+['"]`).FindAllString(code, -1))
	goImp := 0
	for _, blk := range regexp.MustCompile(`(?s)import\s*\(([^)]+)\)`).FindAllStringSubmatch(code, -1) {
		for _, l := range strings.Split(blk[1], "\n") {
			if strings.TrimSpace(l) != "" {
				goImp++
			}
		}
	}
	goImp += len(regexp.MustCompile(`(?m)^\s*import\s+"`).FindAllString(code, -1))
	imports := pyImp + jsImp + javaImp + cppImp + rubyImp + goImp
	fnSet := map[string]bool{}
	for _, f := range regexp.MustCompile(`(?m)\bfunction\s+([A-Za-z_]\w*)|\bdef\s+([A-Za-z_]\w*)`).FindAllStringSubmatch(code, -1) {
		n := f[1]
		if n == "" {
			n = f[2]
		}
		if n != "" && n != "main" {
			fnSet[n] = true
		}
	}
	for _, f := range regexp.MustCompile(`(?m)^func\s+([A-Za-z_]\w*)\s*\(`).FindAllStringSubmatch(code, -1) {
		if f[1] != "main" && f[1] != "init" {
			fnSet[f[1]] = true
		}
	}
	for _, f := range regexp.MustCompile(`(?m)^func\s+\([^)]+\)\s+([A-Za-z_]\w*)\s*\(`).FindAllStringSubmatch(code, -1) {
		fnSet[f[1]] = true
	}
	for _, f := range regexp.MustCompile(`(?m)^\s+(?:public|private|protected|static|final|override)\s+\S+\s+([A-Za-z_]\w*)\s*\(`).FindAllStringSubmatch(code, -1) {
		if f[1] != "main" {
			fnSet[f[1]] = true
		}
	}
	for _, f := range regexp.MustCompile(`(?m)(?:const|let|var)\s+([A-Za-z_]\w*)\s*=\s*(?:async\s*)?\(`).FindAllStringSubmatch(code, -1) {
		fnSet[f[1]] = true
	}
	fnCount := len(fnSet)
	loops := len(regexp.MustCompile(`\b(for|while|do)\s*[({]`).FindAllString(code, -1))
	varSet := map[string]bool{}
	for _, v := range regexp.MustCompile(`(?m)^\s+([a-z][A-Za-z0-9]*)\s*:=`).FindAllStringSubmatch(code, -1) {
		varSet[v[1]] = true
	}
	for _, v := range regexp.MustCompile(`(?m)\b(?:let|const)\s+([A-Za-z_]\w*)`).FindAllStringSubmatch(code, -1) {
		varSet[v[1]] = true
	}
	for _, v := range regexp.MustCompile(`(?m)^\s{4}([a-z_][a-z0-9_]*)\s*=[^=<>!]`).FindAllStringSubmatch(code, -1) {
		varSet[v[1]] = true
	}
	varCount := len(varSet)
	shortFail := false
	for n := range varSet {
		if len(n) > 3 {
			shortFail = true
			break
		}
	}
	extMap := map[string]string{"go": "Go", "py": "Python", "js": "JavaScript", "ts": "TypeScript", "java": "Java", "rb": "Ruby", "rs": "Rust", "php": "PHP", "cs": "C#", "c": "C", "cpp": "C++", "sh": "Bash", "bash": "Bash"}
	topLang := meta.Language
	if topLang == "" {
		best, bestN := "", 0
		for ext, n := range langCnt {
			if n > bestN {
				bestN = n
				best = ext
			}
		}
		if best != "" {
			topLang = extMap[best]
			if topLang == "" {
				topLang = strings.ToUpper(best)
			}
		}
	}
	if topLang == "" {
		topLang = "n/a"
	}
	errH := len(regexp.MustCompile(`if\s+err\s*!=|try\s*\{|except[\s(:]|catch[\s({]|\.expect\(|\.unwrap_or`).FindAllString(code, -1))
	errS := len(regexp.MustCompile(`panic\(|\.unwrap\(\)|,\s*_\s*=`).FindAllString(code, -1))
	errOk := errH > 0 && errS == 0
	slpN := len(regexp.MustCompile(`time\.Sleep|Thread\.sleep|time\.sleep\s*\(`).FindAllString(code, -1))
	fastOk := slpN == 0
	caseN := len(regexp.MustCompile(`\bcase\b`).FindAllString(code, -1))
	stateKw := len(regexp.MustCompile(`(?i)\b(state|mode|phase|status)\s*[=:]`).FindAllString(code, -1))
	stateOk := (caseN >= 2 && caseN <= 6) || stateKw >= 2
	msg = "  ✓  Analysis complete\n"
	fmt.Print(msg)
	rep.WriteString(msg)
	d1Pass := []string{}
	if imports == 0 {
		d1Pass = append(d1Pass, "No-Import Rookie")
	}
	if varCount <= 8 {
		d1Pass = append(d1Pass, "Few-Variable Hero")
	}
	if fnCount <= 1 {
		d1Pass = append(d1Pass, "Single-Function Master")
	}
	if errOk {
		d1Pass = append(d1Pass, "Error-Proof Coder")
	}
	if loops <= 1 {
		d1Pass = append(d1Pass, "One-Loop Warrior")
	}
	if !shortFail {
		d1Pass = append(d1Pass, "Short-Name Ninja")
	}
	if fastOk {
		d1Pass = append(d1Pass, "Fast-Response Builder")
	}
	if stateOk {
		d1Pass = append(d1Pass, "Simple-State Creator")
	}
	budgets := []struct {
		name string
		max  int
	}{
		{"Tiny Scripter", 50}, {"Mini Builder", 100}, {"Compact Coder", 150},
		{"Standard Maker", 200}, {"Detailed Creator", 300}, {"Feature-Rich Dev", 400},
		{"Professional Builder", 500}, {"Enterprise Creator", 650},
	}
	d2Pass := []string{}
	for _, b := range budgets {
		if codeLines <= b.max {
			d2Pass = append(d2Pass, fmt.Sprintf("%s (%d lines)", b.name, codeLines))
			break
		}
	}
	if readmeLow != "" {
		metaHigh += " " + readmeLow
	}
	codeLow := strings.ToLower(code)
	paths := pathLow.String()
	domains := []struct {
		name string
		keys []string
	}{
		{"Basic Tools", []string{"convert", "calculator", "encoder", "generator", "tool", "audit"}},
		{"Simple Games", []string{"game", "tic", "hangman", "puzzle", "player", "score"}},
		{"Text Processing", []string{"text", "format", "parser", "search", "editor", "markdown"}},
		{"Number Crunching", []string{"math", "statistic", "algorithm", "solver", "numeric", "prime"}},
		{"File Management", []string{"file", "folder", "directory", "organize", "rename", "copy"}},
		{"Quiz Systems", []string{"quiz", "trivia", "flashcard", "assessment", "question", "answer"}},
		{"Visual Creation", []string{"ascii", "chart", "graphic", "visual", "draw", "render"}},
		{"Mini Databases", []string{"record", "inventory", "contact", "database", "sqlite", "store"}},
		{"Data Processing", []string{"data", "pipeline", "transform", "validate", "csv", "json"}},
		{"System Utilities", []string{"monitor", "clean", "health", "util", "daemon", "process"}},
	}
	type d3Rank struct {
		name  string
		score int
	}
	ranks := []d3Rank{}
	for _, d := range domains {
		sc := 0
		for _, k := range d.keys {
			wb := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(k) + `\b`)
			sc += 4 * len(wb.FindAllString(metaHigh, -1))
			sc += 2 * len(wb.FindAllString(paths, -1))
			sc += len(wb.FindAllString(codeLow, -1))
		}
		if sc > 0 {
			ranks = append(ranks, d3Rank{d.name, sc})
		}
	}
	for i := 0; i < len(ranks); i++ {
		for j := i + 1; j < len(ranks); j++ {
			if ranks[j].score > ranks[i].score {
				ranks[i], ranks[j] = ranks[j], ranks[i]
			}
		}
	}
	d3Total := 0
	for _, r := range ranks {
		d3Total += r.score
	}
	d3Pass := []string{}
	for i := 0; i < len(ranks) && i < 3; i++ {
		pct := 0
		if d3Total > 0 {
			pct = ranks[i].score * 100 / d3Total
		}
		d3Pass = append(d3Pass, fmt.Sprintf("%s (%d, %d%%)", ranks[i].name, ranks[i].score, pct))
	}
	d4Pass := []string{}
	if topLang != "n/a" {
		d4Pass = append(d4Pass, topLang)
	}
	tables := []struct {
		title string
		items []string
	}{
		{"D1 -- Core Constraints", d1Pass},
		{"D2 -- Line Budget (tightest tier)", d2Pass},
		{"D3 -- Project Domain (top 3, scored)", d3Pass},
		{"D4 -- Primary Language", d4Pass},
	}
	rep.WriteString("\n  ╔════════════════════════════════════════════════════════╗\n")
	rep.WriteString("  ║           PASSED CONSTRAINTS REPORT                  ║\n")
	rep.WriteString("  ╚════════════════════════════════════════════════════════╝\n")
	fmt.Println("\n  ╔════════════════════════════════════════════════════════╗")
	fmt.Println("  ║           PASSED CONSTRAINTS REPORT                  ║")
	fmt.Println("  ╚════════════════════════════════════════════════════════╝")
	sum := fmt.Sprintf("\n  Scan: %d files | Language: %s | Countable lines: %d\n", fileN, topLang, codeLines)
	fmt.Print(sum)
	rep.WriteString(sum)
	total := 0
	for _, tb := range tables {
		fmt.Println()
		fmt.Println("  " + tb.title)
		rep.WriteString("\n  " + tb.title + "\n")
		if len(tb.items) == 0 {
			fmt.Println("  +------------------------------------------------------------------+")
			fmt.Println("  |  (none passed)                                                   |")
			fmt.Println("  +------------------------------------------------------------------+")
			rep.WriteString("  +------------------------------------------------------------------+\n  |  (none passed)                                                   |\n  +------------------------------------------------------------------+\n")
			continue
		}
		sep := "+"
		for range tb.items {
			sep += "----------------------------------+"
		}
		fmt.Println("  " + sep)
		rep.WriteString("  " + sep + "\n  |")
		fmt.Print("  |")
		for _, it := range tb.items {
			cell := fmt.Sprintf(" %-32s|", it)
			fmt.Print(cell)
			rep.WriteString(cell)
		}
		fmt.Println()
		fmt.Println("  " + sep)
		rep.WriteString("\n  " + sep + "\n")
		total += len(tb.items)
	}
	foot := fmt.Sprintf("\n  Total passed: %d\n\n  Note: Heuristic scan; max 80 files; D3 %% = share of total domain score across all matches.\n  Report saved: co-check-report.txt\n\n  Done.\n", total)
	fmt.Print(foot)
	rep.WriteString(foot)
	os.WriteFile("co-check-report.txt", []byte(rep.String()), 0644)
	return nil
}
