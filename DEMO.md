# Demo Walkthrough

## 1. Build & verify

```bash
go vet ./...
gofmt -l .
go test -v ./...
go build -o co-check .
```

## 2. Run on a public repo

```bash
# optional — avoids GitHub API rate limits
export GITHUB_TOKEN="ghp_xxx"

./co-check https://github.com/owner/repo
# or interactive:
go run .
```

## 3. Expected output

- Progress spinners for each GitHub API phase
- **PASSED CONSTRAINTS REPORT** with horizontal tables (D1–D4)
- D2: tightest passing line tier
- D3: up to 3 domains with score and % share
- `co-check-report.txt` saved in the working directory

## 4. Landing page

Open `index.html` in a browser for the project overview (no build step).

## 5. Optional — Python Rosetta Stone

```bash
python co-check-rosetta.py https://github.com/owner/repo
```

Same audit idea with a decomposed Python architecture (stdlib only).
