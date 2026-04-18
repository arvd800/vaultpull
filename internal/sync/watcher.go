package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/your-org/vaultpull/internal/envfile"
)

// WatchAndSync watches the given env file and re-runs the syncer whenever
// the file changes on disk (e.g. manual edits). It blocks until ctx is done.
func WatchAndSync(ctx context.Context, s *Syncer, path string, interval time.Duration) error {
	fmt.Printf("[watch] monitoring %s every %s\n", path, interval)
	return envfile.Watch(ctx, path, envfile.WatchOptions{
		Interval: interval,
		OnChange: func(p string) {
			fmt.Printf("[watch] change detected in %s — re-syncing\n", p)
			if err := s.Run(ctx); err != nil {
				fmt.Printf("[watch] sync error: %v\n", err)
			}
		},
	})
}
