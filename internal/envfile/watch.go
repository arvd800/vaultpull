package envfile

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"time"
)

// WatchOptions configures the file watcher.
type WatchOptions struct {
	Interval time.Duration
	OnChange func(path string)
}

// hashFile returns a SHA-256 hex digest of the file at path.
func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Watch polls path at opts.Interval and calls opts.OnChange when the file
// content changes. It blocks until ctx is cancelled.
func Watch(ctx context.Context, path string, opts WatchOptions) error {
	if opts.Interval <= 0 {
		opts.Interval = 10 * time.Second
	}
	last, _ := hashFile(path)
	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			current, err := hashFile(path)
			if err != nil {
				continue
			}
			if current != last {
				last = current
				if opts.OnChange != nil {
					opts.OnChange(path)
				}
			}
		}
	}
}
