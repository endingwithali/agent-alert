package sessionpoller

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPollReadsMostRecentlyModifiedFile(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	olderPath := writeTestFile(t, root, "2025/01/01/older.jsonl", "older")
	newerPath := writeTestFile(t, root, "2025/01/02/newer.jsonl", "newer")

	now := time.Now()
	if err := os.Chtimes(olderPath, now.Add(-time.Minute), now.Add(-time.Minute)); err != nil {
		t.Fatalf("set older file time: %v", err)
	}
	if err := os.Chtimes(newerPath, now, now); err != nil {
		t.Fatalf("set newer file time: %v", err)
	}

	poller, err := New(root)
	if err != nil {
		t.Fatalf("create poller: %v", err)
	}

	snapshot, err := poller.Poll()
	if err != nil {
		t.Fatalf("poll: %v", err)
	}

	if snapshot.Path != newerPath {
		t.Fatalf("path = %q, want %q", snapshot.Path, newerPath)
	}
	if got, want := string(snapshot.Data), "newer"; got != want {
		t.Fatalf("data = %q, want %q", got, want)
	}
}

func TestPollReadsConfiguredFile(t *testing.T) {
	t.Parallel()

	path := writeTestFile(t, t.TempDir(), "session.jsonl", "session data")
	poller, err := New(path)
	if err != nil {
		t.Fatalf("create poller: %v", err)
	}

	snapshot, err := poller.Poll()
	if err != nil {
		t.Fatalf("poll: %v", err)
	}
	if snapshot.Path != path {
		t.Fatalf("path = %q, want %q", snapshot.Path, path)
	}
	if got, want := string(snapshot.Data), "session data"; got != want {
		t.Fatalf("data = %q, want %q", got, want)
	}
}

func TestPollReturnsErrorWhenNoSessionFilesExist(t *testing.T) {
	t.Parallel()

	poller, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("create poller: %v", err)
	}

	_, err = poller.Poll()
	if !errors.Is(err, ErrNoSessionFiles) {
		t.Fatalf("error = %v, want ErrNoSessionFiles", err)
	}
}

func TestNewRejectsEmptyDirectory(t *testing.T) {
	t.Parallel()

	if _, err := New(""); err == nil {
		t.Fatal("New returned nil error for an empty directory")
	}
}

func writeTestFile(t *testing.T, root, relativePath, contents string) string {
	t.Helper()

	path := filepath.Join(root, relativePath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create test directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	return path
}
