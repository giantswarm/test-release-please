# test-release-please

Throwaway repo for testing CI tooling. Currently configured as a PoC of the
**push-based release flow** (no Release PR; conventional commits on `main`
auto-create a tag, the tag triggers publishing).

## How it works

```
push to main
   │
   ▼
.github/workflows/auto-release.yaml  ── git-cliff inspects commits since last tag
                                  ── creates and pushes vX.Y.Z if bump warranted
   │
   ▼ (tag push)
.github/workflows/release.yaml   ── goreleaser builds binaries, opens GH Release
```

For chart/service repos the `release.yaml` would be replaced by the existing
CircleCI architect pipeline (same tag trigger, different publisher).

## Files

| Path | Purpose |
|---|---|
| `.github/workflows/auto-release.yaml` | Tagger — runs on push to main |
| `.github/workflows/release.yaml`  | Publisher — runs on tag push, invokes goreleaser |
| `cliff.toml`                      | git-cliff config: bump rules + (unused here) changelog template |
| `.goreleaser.yaml`                | goreleaser config: builds, archives, release notes |
| `main.go` / `go.mod`              | Minimal Go CLI so goreleaser has something to build |

## How to test

Land any conventional commit on `main`:

| Commit message | Effect |
|---|---|
| `fix: something` | Patch bump (e.g. v1.0.0 → v1.0.1) |
| `feat: something` | Minor bump (e.g. v1.0.0 → v1.1.0) |
| `feat!: something` or `BREAKING CHANGE:` in body | Major bump (e.g. v1.0.0 → v2.0.0) |
| `chore:` / `docs:` / `ci:` / `test:` / `build:` / `style:` | No release |

After the tag is pushed, watch the **Release** workflow run; it'll produce a
GitHub Release with binaries under "Assets".
