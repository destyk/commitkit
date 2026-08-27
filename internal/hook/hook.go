package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Options controls hook installation behaviour.
type Options struct {
	// RepoRoot is the path to the git repository root.
	// If empty, the current working directory is used and walked upwards
	// until a .git directory is found.
	RepoRoot string

	// Force overwrites an existing hook that was not installed by commitkit.
	Force bool

	// Script overrides the default hook script body.
	Script string
}

// Result describes what Install did.
type Result struct {
	HookPath string
	Created  bool
	Updated  bool
	Skipped  bool
	Message  string
}

// Install installs the commit-msg hook into the repository.
func Install(opts Options) (Result, error) {
	root, err := resolveRepoRoot(opts.RepoRoot)
	if err != nil {
		return Result{}, err
	}

	hooksDir := filepath.Join(root, ".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return Result{}, fmt.Errorf("create hooks directory: %w", err)
	}

	hookPath := filepath.Join(hooksDir, commitMsgHookName)
	script := opts.Script
	if script == "" {
		script = string(defaultCommitMsgHookScript)
	}
	// Ensure trailing newline.
	if !strings.HasSuffix(script, "\n") {
		script += "\n"
	}

	existing, err := os.ReadFile(hookPath)
	if err != nil && !os.IsNotExist(err) {
		return Result{}, fmt.Errorf("read existing hook: %w", err)
	}

	if err == nil {
		// File exists.
		if string(existing) == script {
			return Result{
				HookPath: hookPath,
				Skipped:  true,
				Message:  "commit-msg hook already installed",
			}, nil
		}

		if !opts.Force && !isManagedByCommitkit(string(existing)) {
			return Result{}, fmt.Errorf(
				"hook already exists at %s and does not appear to be managed by commitkit; re-run with --force to overwrite",
				hookPath,
			)
		}

		if err := os.WriteFile(hookPath, []byte(script), 0o755); err != nil {
			return Result{}, fmt.Errorf("write hook: %w", err)
		}

		// Ensure executable bit (WriteFile may not set it on all platforms the same way).
		if err := os.Chmod(hookPath, 0o755); err != nil {
			return Result{}, fmt.Errorf("chmod hook: %w", err)
		}

		return Result{
			HookPath: hookPath,
			Updated:  true,
			Message:  "commit-msg hook updated",
		}, nil
	}

	// Does not exist – create.
	if err := os.WriteFile(hookPath, []byte(script), 0o755); err != nil {
		return Result{}, fmt.Errorf("write hook: %w", err)
	}

	if err := os.Chmod(hookPath, 0o755); err != nil {
		return Result{}, fmt.Errorf("chmod hook: %w", err)
	}

	return Result{
		HookPath: hookPath,
		Created:  true,
		Message:  "commit-msg hook installed",
	}, nil
}

// Uninstall removes the commit-msg hook if it is managed by commitkit
// (or if Force is true).
func Uninstall(opts Options) (Result, error) {
	root, err := resolveRepoRoot(opts.RepoRoot)
	if err != nil {
		return Result{}, err
	}

	hookPath := filepath.Join(root, ".git", "hooks", commitMsgHookName)

	existing, err := os.ReadFile(hookPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Result{
				HookPath: hookPath,
				Skipped:  true,
				Message:  "commit-msg hook not present",
			}, nil
		}

		return Result{}, fmt.Errorf("read hook: %w", err)
	}

	if !opts.Force && !isManagedByCommitkit(string(existing)) {
		return Result{}, fmt.Errorf(
			"hook at %s does not appear to be managed by commitkit; re-run with --force to remove anyway",
			hookPath,
		)
	}

	if err := os.Remove(hookPath); err != nil {
		return Result{}, fmt.Errorf("remove hook: %w", err)
	}

	return Result{
		HookPath: hookPath,
		Message:  "commit-msg hook removed",
	}, nil
}

func isManagedByCommitkit(content string) bool {
	return strings.Contains(content, "Managed by commitkit")
}

func resolveRepoRoot(explicit string) (string, error) {
	if explicit != "" {
		gitDir := filepath.Join(explicit, ".git")
		info, err := os.Stat(gitDir)
		if err != nil || !info.IsDir() {
			// .git may be a file (worktree / submodule)
			if err == nil {
				return explicit, nil
			}
			return "", fmt.Errorf("%s is not a git repository", explicit)
		}

		return explicit, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := cwd
	for {
		gitDir := filepath.Join(dir, ".git")
		if info, err := os.Stat(gitDir); err == nil {
			_ = info
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf(
				"not inside a git repository (no .git found from %s)",
				cwd,
			)
		}

		dir = parent
	}
}
