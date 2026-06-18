package project

import "runtime/debug"

// dev is the default for unset build identifiers.
const dev = "dev"

// devel is the module version runtime/debug reports for a build that carries
// no resolvable VCS tag (no .git, or built outside a module checkout).
const devel = "(devel)"

// Build identifiers. Mirrors muster's pkg/project: a CI go-build job overrides
// gitSHA and buildTimestamp via -X ldflags but does NOT inject version. The
// version is derived at runtime from the Go build info (see Version).
var (
	version        = dev
	gitSHA         = dev
	buildTimestamp = "unknown"
)

// Version returns the best human-readable build identifier available, in
// order: an explicitly injected version ldflag, the VCS version stamped into
// the Go build info, the injected commit SHA, then "dev".
func Version() string {
	if version != dev && version != "" {
		return version
	}
	if v := buildInfoVersion(); v != "" {
		return v
	}
	if gitSHA != dev {
		return gitSHA
	}
	return dev
}

// buildInfoVersion reads the main module version the Go toolchain embedded from
// version control, or "" when none is usable.
var buildInfoVersion = func() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	if v := info.Main.Version; v != "" && v != devel {
		return v
	}
	return ""
}

// GitSHA returns the commit SHA the binary was built from.
func GitSHA() string { return gitSHA }

// BuildTimestamp returns the build time, or "unknown" when not injected.
func BuildTimestamp() string { return buildTimestamp }
