# test-release-please

Throwaway repo for testing CI tooling. Currently configured as a PoC of the
**push-based release flow** (no Release PR; conventional commits on `main`
auto-tag, the tag triggers publishing).

## How it works

```
push to main
   │
   ▼
.github/workflows/auto-tag.yaml  ── git-cliff inspects commits since last tag
                                  ── creates and pushes vX.Y.Z if bump warranted
   │
   ▼ (tag push)
.github/workflows/goreleaser.yaml ── goreleaser builds binaries, opens GH Release
```

For chart/service repos the `goreleaser.yaml` would be replaced by the existing
CircleCI architect pipeline (same tag trigger, different publisher).

## Files

| Path | Purpose |
|---|---|
| `.github/workflows/auto-tag.yaml` | Tagger — runs on push to `main` or any `release-*` branch |
| `.github/workflows/goreleaser.yaml` | Publisher — runs on tag push, invokes goreleaser |
| `cliff.toml`                      | git-cliff config: bump rules + (unused here) changelog template |
| `.goreleaser.yaml`                | goreleaser config: builds, archives, release notes |
| `main.go` / `go.mod`              | Minimal Go CLI so goreleaser has something to build |
| `CHANGELOG.md`                    | **Frozen** at v1.0.0 (the last release-please-produced version). New releases publish notes only to GitHub Releases. |

## How to test

Land any conventional commit on `main`:

| Commit message | Effect |
|---|---|
| `fix: something` | Patch bump (e.g. v1.0.0 → v1.0.1) |
| `feat: something` | Minor bump (e.g. v1.0.0 → v1.1.0) |
| `feat!: something` or `BREAKING CHANGE:` in body | Major bump (e.g. v1.0.0 → v2.0.0) |
| `chore:` / `docs:` / `ci:` / `test:` / `build:` / `style:` | No release |

After the tag is pushed, watch the **GoReleaser** workflow run; it'll produce
a GitHub Release with binaries under "Assets".

## Backports

If you've shipped `v3.0.0` from `main` and now need to ship a fix on the 2.x
line:

```bash
# 1. Once, when starting the 2.x maintenance line:
git switch -c release-2.x v2.3.5
git push -u origin release-2.x

# 2. Open a PR with `release-2.x` as the base (not main) carrying the fix.
#    PR title follows conventional commits as usual: "fix: backport thing".

# 3. Merge the PR.
#    Auto-tag fires on release-2.x → git-cliff sees commits since v2.3.5,
#    not v3.0.0, because v3.0.0 isn't reachable from release-2.x's history.
#    Tags v2.3.6 and pushes.
#    GoReleaser fires on the v2.3.6 tag → publishes binaries.
```

Two gotchas to keep in mind:

1. **The release branch must contain the workflow files.** GitHub Actions
   loads workflows from the branch being pushed to. If you create a release
   branch from a commit that pre-dates these workflows, cherry-pick
   `.github/workflows/` forward into the new branch before opening backport
   PRs.

2. **`feat:` on a release branch still bumps minor.** Convention is to keep
   release branches `fix:`-only; nothing in the tooling enforces this.
   Reviewer discipline (or a tighter `cliff.toml` on release branches if
   you want to be strict).

## Migrating an existing repo to this flow

What's safe to remove and what's not, for a repo currently on the legacy
manual-Release-PR flow (or release-please):

- **Git tags and GitHub Releases:** keep. They're not tied to which workflow
  created them. All historical release pages stay accessible at their URLs.
- **`CHANGELOG.md`:** **do not delete.** Freeze it instead — add a header
  noting the version at which the cutoff happened and pointing readers at
  the GitHub Releases page going forward. See this repo's own
  `CHANGELOG.md` as the example. New releases publish their notes to GitHub
  Releases only; nothing writes back to this file.
- **Legacy workflow files** (`zz_generated.create_release_pr.yaml`,
  `zz_generated.create_release.yaml`, `zz_generated.validate_changelog.yaml`):
  remove. devctl will handle this when the team's `releaseWorkflow` config
  is switched to the push-based mode.
- **`release-please-config.json` / `.release-please-manifest.json`** (if
  migrating from release-please): remove. devctl handles this too.
