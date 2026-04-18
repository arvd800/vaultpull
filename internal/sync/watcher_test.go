package sync

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestWatchAndSync_CancelImmediately(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/.env"
	_ = os.WriteFile(path, []byte("KEY=val\n"), 0600)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before Watch even polls

	s := &Syncer{} // zero value; Run will error but watcher exits via ctx
	err := WatchAndSync(ctx, s, path, 10*time.Millisecond)
	if err == nil {
		t.Error("expected context cancellation error")
	}
}

func TestWatchAndSync_FileNotExistNoHang(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()

	s := &Syncer{}
	// Missing file: hashFile will fail silently each tick, watcher still exits via ctx.
	err := WatchAndSync(ctx, s, "/tmp/nonexistent_vaultpull.env", 10*time.Millisecond)
	if err == nil {
		t.Error("expected context deadline error")
	}
}
