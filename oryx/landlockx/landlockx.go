// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package landlockx

import (
	"os"
	"path/filepath"

	ll "github.com/landlock-lsm/go-landlock/landlock"
	llsyscall "github.com/landlock-lsm/go-landlock/landlock/syscall"

	"github.com/ory/x/logrusx"
)

// IsSupported returns true if the running kernel supports Landlock LSM.
// On non-Linux platforms or kernels older than 5.13, this returns false.
func IsSupported() bool {
	v, err := llsyscall.LandlockGetABIVersion()
	return err == nil && v > 0
}

// runningExe returns the kernel-canonical path of the running binary.
// jsonnetsecure re-execs the running binary as a sandboxed worker, so
// the rule we attach must use the inode the kernel resolves at execve
// time. os.Executable on Linux returns the readlink of /proc/self/exe;
// EvalSymlinks resolves any further symlinks so the rule attaches to
// the actual file inode.
func runningExe() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	if real, err := filepath.EvalSymlinks(exe); err == nil {
		return real
	}
	return exe
}

// Apply activates the Landlock LSM filesystem sandbox for the calling process.
//
// roPaths is a list of paths that the process may read. Each path is stat-ed
// to determine whether to apply directory or file access rights; non-existent
// paths are silently skipped.
//
// rwDirs is a list of directories that the process may read and write
// recursively. Missing paths are silently ignored.
//
// Apply is a no-op on non-Linux platforms and on Linux kernels older than
// 5.13 (BestEffort degrades gracefully).
func Apply(l *logrusx.Logger, roPaths, rwDirs []string) error {
	rules := defaultRules()

	for _, p := range roPaths {
		info, err := os.Stat(p)
		if err != nil {
			continue
		}
		if info.IsDir() {
			rules = append(rules, ll.RODirs(p))
		} else {
			rules = append(rules, ll.ROFiles(p))
		}
	}

	for _, d := range rwDirs {
		rules = append(rules, ll.RWDirs(d).IgnoreIfMissing())
	}

	if !IsSupported() {
		l.Warn("Landlock LSM is not available on this kernel; filesystem sandbox not applied.")
		return nil
	}

	if err := ll.V3.BestEffort().RestrictPaths(rules...); err != nil {
		return err
	}

	l.Info("Landlock filesystem sandbox is active.")
	return nil
}

// ApplyEmpty activates a Landlock filesystem sandbox with no allowed
// paths, denying every filesystem access the V3 ABI handles. Already-open
// file descriptors (stdin, stdout, stderr inherited from the parent)
// keep working, but no path-based open/read/write/exec is permitted.
//
// Use it inside helper subprocesses that only need stdio — for example,
// the hidden `jsonnet` worker that jsonnetsecure re-execs to evaluate
// mappers in isolation.
//
// Landlock layers stack: this call adds a strictly-more-restrictive
// layer on top of whatever the parent already enforced.
//
// ApplyEmpty is a no-op on non-Linux platforms and on Linux kernels
// older than 5.13.
func ApplyEmpty(l *logrusx.Logger) error {
	if !IsSupported() {
		if l != nil {
			l.Warn("Landlock LSM is not available on this kernel; empty sandbox not applied.")
		}
		return nil
	}
	if err := ll.V3.BestEffort().RestrictPaths(); err != nil {
		return err
	}
	if l != nil {
		l.Info("Landlock empty filesystem sandbox is active.")
	}
	return nil
}

// defaultRules returns the rules every kratos serve process needs regardless
// of operator config:
//   - /dev/null is a kernel-managed sink used by subprocess plumbing
//     (exec.Command stdin redirection, log discarding, etc.); it cannot
//     leak data.
//   - /etc/resolv.conf carries the DNS nameserver list and /etc/hosts
//     the static host table; Go's pure-Go resolver reads both for every
//     hostname lookup. /etc/nsswitch.conf is intentionally not granted:
//     the binary is built with CGO_ENABLED=0, so when nsswitch.conf is
//     unreadable the resolver falls back to "files then DNS", which is
//     what we want anyway.
//   - the running binary is re-exec'd by jsonnetsecure for sandboxed mapper
//     evaluation, so it needs read+execute. The binary must be statically
//     linked (CGO_ENABLED=0) — execve of a dynamically-linked binary also
//     requires read+execute on the system dynamic linker (PT_INTERP) and
//     its NEEDED libraries, which we deliberately do not grant.
//
// /etc/ssl (system trust roots for outbound HTTPS) is intentionally not
// included. Cloud binaries blank-import
// golang.org/x/crypto/x509roots/fallback (see cloudlib/x509_roots.go)
// and set `godebug x509usefallbackroots=1` in their go.mod, so
// crypto/x509 ships its own embedded Mozilla bundle and never reads
// /etc/ssl at runtime. Operators who need to trust an additional CA
// continue to set SSL_CERT_FILE / SSL_CERT_DIR and list the file under
// security.landlock.allowed_paths.
func defaultRules() []ll.Rule {
	rules := []ll.Rule{
		ll.RWFiles("/dev/null"),
		ll.ROFiles("/etc/resolv.conf").IgnoreIfMissing(),
		ll.ROFiles("/etc/hosts").IgnoreIfMissing(),
	}
	if exe := runningExe(); exe != "" {
		access := ll.AccessFSSet(llsyscall.AccessFSReadFile | llsyscall.AccessFSExecute)
		rules = append(rules, ll.PathAccess(access, exe).IgnoreIfMissing())
	}
	return rules
}
