package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Annotation holds metadata attached to a secret key.
type Annotation struct {
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Annotations maps secret keys to their annotations.
type Annotations map[string]Annotation

// SaveAnnotations persists annotations to a JSON file.
func SaveAnnotations(path string, annotations Annotations) error {
	if path == "" {
		return nil
	}
	data, err := json.MarshalIndent(annotations, "", "  ")
	if err != nil {
		return fmt.Errorf("annotate: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadAnnotations reads annotations from a JSON file.
// Returns an empty map if the file does not exist.
func LoadAnnotations(path string) (Annotations, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Annotations{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("annotate: read: %w", err)
	}
	var annotations Annotations
	if err := json.Unmarshal(data, &annotations); err != nil {
		return nil, fmt.Errorf("annotate: unmarshal: %w", err)
	}
	return annotations, nil
}

// Annotate sets or updates the note for a key, returning a new Annotations map.
func Annotate(existing Annotations, key, note string) Annotations {
	out := make(Annotations, len(existing))
	for k, v := range existing {
		out[k] = v
	}
	now := time.Now().UTC()
	if prev, ok := out[key]; ok {
		prev.Note = note
		prev.UpdatedAt = now
		out[key] = prev
	} else {
		out[key] = Annotation{
			Note:      note,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return out
}

// RemoveAnnotation removes the annotation for a key, returning a new Annotations map.
func RemoveAnnotation(existing Annotations, key string) Annotations {
	out := make(Annotations, len(existing))
	for k, v := range existing {
		if k != key {
			out[k] = v
		}
	}
	return out
}
