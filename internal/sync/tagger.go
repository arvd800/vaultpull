package sync

import (
	"time"

	"github.com/yourusername/vaultpull/internal/envfile"
)

// TagOptions configures metadata tagging after a sync.
type TagOptions struct {
	TagsFilePath string
	Source       string
	ExtraTags    map[string]string
}

// ApplyTags writes or updates a tags file after a successful sync.
func ApplyTags(opts TagOptions) error {
	if opts.TagsFilePath == "" {
		return nil
	}

	existing, err := envfile.LoadTags(opts.TagsFilePath)
	if err != nil {
		return err
	}

	existing.Source = opts.Source
	existing.FetchedAt = time.Now().UTC()

	if existing.Tags == nil {
		existing.Tags = map[string]string{}
	}
	envfile.MergeTags(&existing, opts.ExtraTags)

	return envfile.SaveTags(opts.TagsFilePath, existing)
}
