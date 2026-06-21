# CodeWeaver — Security Remediation & Skill-Enablement Plan

**Status:** Proposed
**Author:** Security review (consultant-facing)
**Date:** 2026-06-21
**Scope:** `main.go`, `main_test.go`, CI/CD, plus a new global "AI skill" wrapper

---

## 0. Guiding principle (from review feedback)

A file *path* is not a secret; a file's *contents* can be. CodeWeaver's job is to
produce a faithful map of a codebase, so it should keep **listing** sensitive files in
the tree and path-lists, but must never **extract their contents** into the output that
gets shipped to an LLM or the clipboard.

This reframes the headline fix from *exclusion* to **content-level redaction**.

---

## 1. Findings recap (priority order)

| ID | Issue | Fix theme | Priority |
|----|-------|-----------|----------|
| S1 | Secrets' **contents** embedded by default (`.env`, keys, creds) | Content redaction + optional `.gitignore` | P1 |
| S2 | Symlinked files are followed; target read even if outside root | Symlink-escape guard | P1 |
| S3 | Output files written world-readable (`0644`) | `0600` default | P2 |
| S4 | Whole files / whole corpus read into memory, no caps | Size limits | P2 |
| S5 | CI: `contents: write` on PR; mutable action tags | Least-privilege + SHA pins | P2 |
| S6 | Stale transitive `x/*` deps; `govulncheck` not run | Update + CI gate | P3 |
| B1 | `-version` prints nothing; ldflags vars undeclared (no-op) | Wire up version | P2 (bug) |

> Non-issue, recorded to prevent rework: user `-ignore`/`-include` regexes are **not**
> ReDoS-vulnerable. Go's `regexp` is RE2 (linear time). No change needed.

---

## 2. Detailed changes

### S1 — Redact contents of sensitive files (keep them listed) · P1

**Design:** introduce a third filter dimension distinct from ignore/include:
*content sensitivity*. A file can pass `shouldProcess` (so it appears in the tree and
`included-paths` list) yet have its body suppressed.

1. Add a curated default sensitivity denylist (case-insensitive, path-based regex), e.g.:
   - `(^|/)\.env(\..*)?$`, `(^|/)\.envrc$`
   - `\.pem$`, `\.key$`, `\.pfx$`, `\.p12$`, `\.keystore$`
   - `(^|/)id_(rsa|dsa|ecdsa|ed25519)$`
   - `(^|/)\.npmrc$`, `(^|/)\.pypirc$`, `(^|/)\.netrc$`
   - `(^|/)credentials$`, `(^|/)\.aws/`, `(^|/)\.ssh/`
   - `\.tfstate$`, `\.tfvars$`
   - `(^|/)secrets?\.(ya?ml|json|toml|ini)$`
2. In `buildContentString`, after a file passes `shouldProcess` and **before**
   `os.ReadFile`, test the relative path against the denylist. On match:
   - Append the file header + a fenced placeholder:
     `[content redacted by CodeWeaver — matched sensitive pattern: <pattern>]`
   - Still record it in `localProcessedPaths` and the extension summary (so the tree and
     listings stay complete).
   - **Do not read the file** (avoids loading secrets into memory at all).
3. New flags:
   - `-redact` (comma-separated regexes) — append to the default denylist.
   - `-no-default-redact` — drop the built-in denylist (explicit opt-out).
   - `-unsafe-include-secrets` — escape hatch to actually embed contents (off by default;
     name is intentionally alarming).
4. **Optional `.gitignore` honoring** as a broader net:
   - `-gitignore` flag: parse the input dir's `.gitignore` and treat matched paths as
     `-ignore` entries (these are fully excluded — gitignored build artifacts/secrets
     usually shouldn't even be listed). Keep off by default to preserve current behavior,
     but document it as the recommended safe mode.

**Why both:** redaction guarantees no secret *content* ever leaks even if `.gitignore`
is missing or incomplete; `.gitignore` mode additionally trims noise and locally-managed
secret files entirely. Layered defense.

### S2 — Symlink-escape guard · P1

- After computing `currentWalkPath`, resolve with `filepath.EvalSymlinks`. If the real
  path is **outside** `rootAbsPath` (compare via `filepath.Rel` → no `..` prefix), skip
  the entry and log a warning. Applies to both content reading and tree display.
- Add `-follow-symlinks` opt-in to restore following (still bounded to root unless an
  additional `-allow-symlink-escape` is set).
- Test: a symlink pointing to `/etc/hosts` (or a temp file outside root) is skipped.

### S3 — Tighten output permissions · P2

- `os.WriteFile(... , 0600)` in `writeOutput` and `savePathsToFile`. The aggregate is a
  derived artifact that may concentrate sensitive data; `0600` is the correct default.

### S4 — Size limits · P2

- `-max-file-size` (bytes, default e.g. `5_000_000`). Files above it are listed but their
  content replaced with `[content skipped — file exceeds max-file-size]`.
- Optional `-max-total-size` guard that aborts with a clear error before OOM.
- Prefer `os.Open` + `io.LimitReader` over `os.ReadFile` so oversized files aren't fully
  buffered.

### S5 — CI hardening · P2

- Split triggers: a lightweight build/test job on `pull_request` with
  `permissions: contents: read`; the `release` job runs only on `tags` with
  `contents: write`.
- Pin `actions/checkout`, `actions/setup-go`, `goreleaser/goreleaser-action` to full
  commit SHAs (keep the human-readable tag in a trailing comment).
- Add a `govulncheck` step (see S6).

### S6 — Dependencies · P3

- `go install golang.org/x/vuln/cmd/govulncheck@latest && govulncheck ./...` — record
  results; this is the authoritative gate before sign-off.
- `go get -u golang.org/x/sys golang.org/x/image golang.org/x/mobile && go mod tidy`.
- Add `govulncheck ./...` to CI so regressions are caught.

### B1 — Fix `-version` (functional bug) · P2

- Declare package-level `var version, commit, date = "dev", "none", "unknown"` so the
  goreleaser `-X` ldflags actually bind.
- In `parseFlags`/`main`, when `-version` is set, print `version (commit, date)` and exit 0.

---

## 3. Test plan (extend `main_test.go`)

- `.env` / `id_rsa` → appears in tree & included-paths, body is the redaction placeholder,
  raw secret string never present in output.
- `-no-default-redact` → content embedded (asserts opt-out works).
- Symlink to a file outside root → skipped, warning emitted, target contents absent.
- Output file mode is `0600`.
- File over `-max-file-size` → listed, body is the skip placeholder.
- `.gitignore` mode excludes a gitignored path entirely.
- Regression: existing tree/content/dynamic-fencing tests still pass.

---

## 4. Rollout

1. Branch `security/redaction-and-hardening`.
2. Land S1 + S2 + tests (P1) → review → merge.
3. Land S3/S4/B1, then S5/S6.
4. Tag a release; confirm `-version` stamping works end-to-end via goreleaser.

Backward-compatibility note: defaults change behavior (secret *contents* stop appearing).
This is the intended security improvement; call it out in the changelog as a notable change
and document the `-unsafe-include-secrets` / `-no-default-redact` escape hatches.

---

## 5. Make it a global, location-independent AI skill

Goal: invoke from anywhere in a file tree (like your `ai-model-scanner`) and have it
weave the *current* project regardless of CWD.

### 5a. Code changes that enable portability

- **`-input` already defaults to `.`** and is resolved with `filepath.Abs`, so "run
  wherever I am" works today. Add quality-of-life on top:
  - `-root-marker` (default `.git,go.mod,package.json,pyproject.toml`): walk **up** from
    CWD to the nearest ancestor containing any marker and use that as the input root, so
    the skill weaves the *whole project* even when called from a deep subdirectory.
    Add `-here` to force literal CWD instead.
  - `-output -` to stream Markdown to **stdout** (no temp file) — essential for piping
    into an agent without writing into the user's repo.
  - `-stdout-only` / respect when output is a pipe, so the human log goes to stderr
    (it already uses a logger; point it at `os.Stderr`) and only the Markdown hits stdout.
  - Emit a trailing **size + approx token count** summary (chars/4 heuristic) so the
    caller can gauge context-window fit.

These keep the binary self-contained and dependency-light, which matters for a tool you'll
drop onto many machines.

### 5b. Skill packaging (Claude Code skill convention)

Create `~/.claude/skills/codeweaver/` so it's available globally:

```
~/.claude/skills/codeweaver/
  SKILL.md            # name, description (trigger phrasing), usage
  scripts/codeweaver  # prebuilt static binary (CGO_ENABLED=0), or
  scripts/run.sh      # thin wrapper: locate project root, run binary, stream stdout
```

- `SKILL.md` description should trigger on phrasings like "weave this codebase",
  "dump the repo to markdown for context", "package the project for an LLM".
- The wrapper runs:
  `codeweaver -input "$(git rev-parse --show-toplevel 2>/dev/null || pwd)" -output - -gitignore`
  i.e. safe-by-default: gitignore-aware, secrets redacted, streamed to stdout.
- Because the build is `CGO_ENABLED=0` static (already configured in `goreleaser.yaml`),
  the same binary is portable across machines with no runtime deps. **Caveat:** the
  `golang.design/x/clipboard` dependency needs CGO/X11 on Linux for the `-clipboard`
  feature; for the headless skill path, clipboard is irrelevant (we stream to stdout), so
  consider a build tag that compiles clipboard out for the skill binary to keep it truly
  static and dependency-free.

### 5c. Optional: distribute as a one-liner

- `go install github.com/tesserato/CodeWeaver@latest` already makes it global on any box
  with Go. Document this plus the skill wrapper so the consultant can choose binary-drop
  vs. `go install`.

---

## 6. Effort estimate

| Phase | Items | Est. |
|-------|-------|------|
| P1 | S1, S2, tests | ~0.5–1 day |
| P2 | S3, S4, S5, B1 | ~0.5 day |
| P3 | S6 | ~1–2 hrs |
| Skill | 5a code + 5b packaging | ~0.5 day |
