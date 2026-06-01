#!/usr/bin/env python3
"""Code Olympics GitHub Constraint Auditor — Python Rosetta Stone (stdlib only)."""
import json
import os
import re
import sys
import urllib.request

API = "https://api.github.com"
MAXF, MAXB, TOUT = 80, 120_000, 25
SKIP = ("node_modules/", "vendor/", ".git/", "dist/", "build/", "__pycache__/")
RE_EXT = re.compile(r"\.(go|py|js|ts|jsx|tsx|java|c|cpp|cc|h|cs|rb|php|rs|sh|bash)$")


def fetch(url: str, token: str, accept: str = "application/vnd.github+json") -> bytes:
    req = urllib.request.Request(url)
    req.add_header("Accept", accept)
    req.add_header("User-Agent", "co-check-py")
    if token:
        req.add_header("Authorization", f"Bearer {token}")
    with urllib.request.urlopen(req, timeout=TOUT) as r:
        return r.read()


def parse_url(raw: str):
    m = re.search(r"github\.com[/:]([^/\s#?]+)/([^/\s#?]+)", raw, re.I)
    if not m:
        raise ValueError("invalid GitHub URL")
    owner, repo = m.group(1), m.group(2)
    if repo.endswith(".git"):
        repo = repo[:-4]
    return owner, repo


def lines(code: str) -> int:
    n, block = 0, False
    for ln in code.replace("\r\n", "\n").split("\n"):
        t = ln.strip()
        if not t:
            continue
        if block:
            if "*/" in t:
                block = False
            continue
        if t.startswith("/*"):
            if t.find("*/") < 0:
                block = True
            continue
        if t.startswith(("//", "#", "--", "<!--", "*", ";")):
            continue
        n += 1
    return n


def analyze(owner: str, repo: str, token: str):
    meta = json.loads(fetch(f"{API}/repos/{owner}/{repo}", token))
    branch = meta.get("default_branch") or "main"
    sha = json.loads(
        fetch(f"{API}/repos/{owner}/{repo}/branches/{branch}", token)
    )["commit"]["sha"]
    tree = json.loads(
        fetch(f"{API}/repos/{owner}/{repo}/git/trees/{sha}?recursive=1", token)
    ).get("tree", [])

    parts, fcount, scount = [], 0, 0
    for ent in tree:
        if ent.get("type") != "blob":
            continue
        p = ent.get("path", "")
        if not RE_EXT.search(p) or ent.get("size", 0) > MAXB:
            continue
        if any(d in p for d in SKIP):
            scount += 1
            continue
        try:
            body = fetch(
                f"{API}/repos/{owner}/{repo}/git/blobs/{ent['sha']}",
                token,
                accept="application/vnd.github.raw",
            ).decode("utf-8", errors="replace")
        except Exception:
            scount += 1
            continue
        parts.append(body)
        fcount += 1
        if fcount >= MAXF:
            break

    code = "\n".join(parts)
    if not code.strip():
        raise RuntimeError("no source files found")

    cl = lines(code)
    imps = len(re.findall(r"(?m)^\s*(import\s+\S|from\s+\S+\s+import)", code))
    fns = len(set(m[0] or m[1] for m in re.findall(r"(?m)\bfunction\s+([A-Za-z_]\w*)|\bdef\s+([A-Za-z_]\w*)", code)))
    loops = len(re.findall(r"\b(for|while|do)\s*[({]", code))
    print(f"\n[Scan] {fcount} files, {scount} skipped | Lines: {cl}")
    print(f"Imports: {imps} | Functions: {fns} | Loops: {loops}")
    print("\nDone.")


def main():
    print("co-check-rosetta — Python Rosetta Stone (stdlib)")
    u = sys.argv[1] if len(sys.argv) > 1 else input("\nGitHub URL: ").strip()
    if not u:
        return
    tok = os.environ.get("GITHUB_TOKEN", "")
    try:
        analyze(*parse_url(u), tok)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
