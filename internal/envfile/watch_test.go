package envfile

import (
	"context"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

func TestWatch_DetectsChange(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "watch*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString("KEY=old\n")
	f.Close()

	var calls atomic.Int32
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- Watch(ctx, f.Name(), WatchOptions{
			Interval: 20 * time.Millisecond,
			OnChange: func(_ string) { calls.Add(1) },
		})
	}()

	time.Sleep(40 * time.Millisecond)
	_ = os.WriteFile(f.Name(), []byte("KEY=new\n"), 0600)
	time.Sleep(60 * time.Millisecond)
	cancel()
	<-done

	if calls.Load() == 0 {
		t.Error("expected OnChange to be called at least once")
	}
}

func TestWatch_NoChangeNoCallback(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "watch*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString("KEY=stable\n")
	f.Close()

	var calls atomic.Int32
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	_ = Watch(ctx, f.Name(), WatchOptions{
		Interval: 20 * time.Millisecond,
		OnChange: func(_ string) { calls.Add(1) },
	})

	if calls.Load() != 0 {
		t.Errorf("expected no callbacks, got %d", calls.Load())
	}
}

func TestHashFile_NonExistent(t *testing.T) {
	_, err := hashFile("/no/such/file.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
