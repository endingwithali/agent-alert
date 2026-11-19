package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadUsesSampleDefaults(t *testing.T) {
	clearConfigEnvironment(t)

	got, err := Load(filepath.Join(t.TempDir(), ".env"))
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if got.SessionPath != "test/data/sample-session.jsonl" {
		t.Fatalf("SessionPath = %q", got.SessionPath)
	}
	if got.PollInterval != 5*time.Second {
		t.Fatalf("PollInterval = %s", got.PollInterval)
	}
}

func TestLoadReadsDotenv(t *testing.T) {
	clearConfigEnvironment(t)

	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("CODEX_SESSION_PATH=/tmp/sessions\nPOLL_INTERVAL=750ms\n"), 0o600); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if got.SessionPath != "/tmp/sessions" {
		t.Fatalf("SessionPath = %q", got.SessionPath)
	}
	if got.PollInterval != 750*time.Millisecond {
		t.Fatalf("PollInterval = %s", got.PollInterval)
	}
}

func TestProcessEnvironmentOverridesDotenv(t *testing.T) {
	clearConfigEnvironment(t)
	t.Setenv("CODEX_SESSION_PATH", "/environment/sessions")

	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("CODEX_SESSION_PATH=/dotenv/sessions\n"), 0o600); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if got.SessionPath != "/environment/sessions" {
		t.Fatalf("SessionPath = %q", got.SessionPath)
	}
}

func clearConfigEnvironment(t *testing.T) {
	t.Helper()
	for _, key := range []string{"CODEX_SESSION_PATH", "CODEX_SESSION_DIR", "POLL_INTERVAL"} {
		oldValue, existed := os.LookupEnv(key)
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("unset %s: %v", key, err)
		}
		t.Cleanup(func() {
			if existed {
				_ = os.Setenv(key, oldValue)
			} else {
				_ = os.Unsetenv(key)
			}
		})
	}
}
