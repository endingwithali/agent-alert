// Package sessionpoller finds and reads the most recently modified Codex session file.
package sessionpoller

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// ErrNoSessionFiles indicates that the configured session directory contains no files.
var ErrNoSessionFiles = errors.New("no session files found")

// Snapshot contains the data read from the most recently modified session file.
type Snapshot struct {
	Path       string
	Data       []byte
	ModifiedAt time.Time
}

// Poller reads a specific session file or the newest file in a directory tree.
type Poller struct {
	sessionPath string
}

// New creates a Poller for a session file or directory.
func New(sessionPath string) (*Poller, error) {
	if sessionPath == "" {
		return nil, errors.New("session path must not be empty")
	}

	return &Poller{sessionPath: filepath.Clean(sessionPath)}, nil
}

// Poll reads the configured file, or finds and reads the newest regular file
// when the configured path is a directory.
func (p *Poller) Poll() (Snapshot, error) {
	path, info, err := selectedFile(p.sessionPath)
	if err != nil {
		return Snapshot{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("read most recent session file %q: %w", path, err)
	}

	return Snapshot{
		Path:       path,
		Data:       data,
		ModifiedAt: info.ModTime(),
	}, nil
}

func selectedFile(path string) (string, fs.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", nil, fmt.Errorf("inspect session path %q: %w", path, err)
	}
	if info.Mode().IsRegular() {
		return path, info, nil
	}
	if !info.IsDir() {
		return "", nil, fmt.Errorf("session path %q is not a regular file or directory", path)
	}

	return newestFile(path)
}

func newestFile(root string) (string, fs.FileInfo, error) {
	var newestPath string
	var newestInfo fs.FileInfo

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("get file information for %q: %w", path, err)
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		if newestInfo == nil || info.ModTime().After(newestInfo.ModTime()) ||
			(info.ModTime().Equal(newestInfo.ModTime()) && path > newestPath) {
			newestPath = path
			newestInfo = info
		}

		return nil
	})
	if err != nil {
		return "", nil, fmt.Errorf("scan session directory %q: %w", root, err)
	}
	if newestInfo == nil {
		return "", nil, fmt.Errorf("%w in %q", ErrNoSessionFiles, root)
	}

	return newestPath, newestInfo, nil
}
