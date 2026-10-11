package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindIssuesDir(t *testing.T) {
	// Create temp directory structure
	tmp := t.TempDir()
	gitRoot := filepath.Join(tmp, "project")
	subDir := filepath.Join(gitRoot, "src", "pkg")
	issuesDir := filepath.Join(gitRoot, IssuesDirName)

	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(gitRoot, ".git"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(issuesDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Should find .issues from subdirectory
	found := FindIssuesDir(subDir)
	if found != gitRoot {
		t.Errorf("FindIssuesDir(%s) = %q, want %q", subDir, found, gitRoot)
	}

	// Should find .issues from root
	found = FindIssuesDir(gitRoot)
	if found != gitRoot {
		t.Errorf("FindIssuesDir(%s) = %q, want %q", gitRoot, found, gitRoot)
	}
}

func TestFindIssuesDirStopsAtGitRoot(t *testing.T) {
	tmp := t.TempDir()
	gitRoot := filepath.Join(tmp, "project")
	subDir := filepath.Join(gitRoot, "src")

	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(gitRoot, ".git"), 0755); err != nil {
		t.Fatal(err)
	}
	// No .issues directory

	found := FindIssuesDir(subDir)
	if found != "" {
		t.Errorf("FindIssuesDir(%s) = %q, want empty (should stop at git root)", subDir, found)
	}
}

func TestFindIssuesDirEnvOverride(t *testing.T) {
	tmp := t.TempDir()
	customIssues := filepath.Join(tmp, "custom-issues")
	if err := os.MkdirAll(customIssues, 0755); err != nil {
		t.Fatal(err)
	}

	t.Setenv(EnvIssuesDir, customIssues)

	found := FindIssuesDir("/some/random/dir")
	if found != tmp {
		t.Errorf("FindIssuesDir with env = %q, want %q", found, tmp)
	}
}

func TestFindIssuesDirEnvOverrideMissingDir(t *testing.T) {
	tmp := t.TempDir()
	customIssues := filepath.Join(tmp, "does-not-exist")
	t.Setenv(EnvIssuesDir, customIssues)

	found := FindIssuesDir("/some/random/dir")
	if found != tmp {
		t.Errorf("FindIssuesDir with env (missing dir) = %q, want %q", found, tmp)
	}
}

func TestIssuesDirPathEnvOverride(t *testing.T) {
	tmp := t.TempDir()
	customIssues := filepath.Join(tmp, "custom-issues")
	t.Setenv(EnvIssuesDir, customIssues)

	p := New(tmp)
	if p.IssuesDir != customIssues {
		t.Errorf("IssuesDir = %q, want %q", p.IssuesDir, customIssues)
	}
	if p.OpenDir != filepath.Join(customIssues, OpenDirName) {
		t.Errorf("OpenDir = %q", p.OpenDir)
	}
	if p.ConfigPath != filepath.Join(customIssues, SyncDirName, ConfigFileName) {
		t.Errorf("ConfigPath = %q", p.ConfigPath)
	}
}

func TestIssuesDirPathEnvOverrideRelative(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv(EnvIssuesDir, "work/.issues")

	p := New(tmp)
	want := filepath.Join(tmp, "work", ".issues")
	if p.IssuesDir != want {
		t.Errorf("IssuesDir = %q, want %q", p.IssuesDir, want)
	}
}

func TestIssuesDirPathNoEnv(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv(EnvIssuesDir, "")

	p := New(tmp)
	if p.IssuesDir != filepath.Join(tmp, IssuesDirName) {
		t.Errorf("IssuesDir = %q", p.IssuesDir)
	}
	if p.Root != tmp {
		t.Errorf("Root = %q, want %q", p.Root, tmp)
	}
}

func TestFindIssuesDirAndNewRoundTripRelativeEnv(t *testing.T) {
	// Mirrors the production pipeline: cmd/gh-issue-sync/main.go sets
	// root := FindIssuesDir(cwd), then App.Init calls paths.New(root). A
	// relative GH_ISSUE_SYNC_DIR must resolve to the same location in both
	// steps.
	base := t.TempDir()
	t.Setenv(EnvIssuesDir, filepath.Join("work", ".issues"))

	root := FindIssuesDir(base)
	p := New(root)
	want := filepath.Join(base, "work", ".issues")
	if p.IssuesDir != want {
		t.Fatalf("round-trip IssuesDir = %q, want %q", p.IssuesDir, want)
	}
}

func TestFindGitRoot(t *testing.T) {
	tmp := t.TempDir()
	gitRoot := filepath.Join(tmp, "project")
	subDir := filepath.Join(gitRoot, "src", "pkg")

	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(gitRoot, ".git"), 0755); err != nil {
		t.Fatal(err)
	}

	found := FindGitRoot(subDir)
	if found != gitRoot {
		t.Errorf("FindGitRoot(%s) = %q, want %q", subDir, found, gitRoot)
	}

	// No git root
	found = FindGitRoot(tmp)
	if found != "" {
		t.Errorf("FindGitRoot(%s) = %q, want empty", tmp, found)
	}
}
