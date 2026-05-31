# Demo Walkthrough

## 1. Build & Verify Zero Warnings

```bash
go vet ./...
gofmt -l .
go test -v ./...
go build -o co-check .
```

All four commands should exit cleanly (zero warnings, zero test failures, binary produced).

## 2. Run on a Public Repo

```bash
export GITHUB_TOKEN="ghp_xxx"  # optional but recommended
echo "https://github.com/golang/example" | ./co-check
```

Expected output includes:
- **D1 Core Constraints** PASS/FAIL with evidence counts
- **D2 Line Budget** comparison against all 8 tiers
- **D3 Domain** keyword matching
- **D4 Language** detection
- **Code Olympics Score** (0-100)
- **Best Fit Tier** recommendation

## 3. Run with Assigned Constraints

When prompted, enter:
- D1: `Single-Function Master`
- D2: `Feature-Rich Dev`
- D3: `Basic Tools`
- D4: `Go`

You will see a targeted verdict:

```
CHALLENGE VERDICT: ALL ASSIGNED CONSTRAINTS MET
```

(or a gap report if any fail).

## 4. Inspect the Gap Report

If any D1 rules fail, the tool prints `What To Fix` with actionable advice for each failed rule.

## 5. Run the Rosetta Stone (Python)

```bash
python co-check.py
# Paste the same GitHub URL and compare output architecture.
```

Observe how Python's natural decomposition into `fetch`, `parse_url`, `lines`, and `analyze` contrasts with Go's single-function imperative scanner.

## 6. Read the Cross-Constraint Combo

Open `README.md` and scroll to **Cross-Constraint Combo** to see how the collision of Single-Function Master × Feature-Rich Dev × Go produced the "Global Pattern Library + Sequential Imperative Scanner" emergent architecture.
