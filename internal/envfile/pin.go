package envfile

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// PinRecord holds a pinned value for a key with metadata.
type PinRecord struct {
	Value     string    `json:"value"`
	PinnedAt  time.Time `json:"pinned_at"`
	PinnedBy  string    `json:"pinned_by,omitempty"`
}

// PinMap maps env key names to their pinned records.
type PinMap map[string]PinRecord

// SavePins writes the pin map to a JSON file at path.
func SavePins(path string, pins PinMap) error {
	if path == "" {
		return errors.New("pin: empty path")
	}
	data, err := json.MarshalIndent(pins, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadPins reads a pin map from a JSON file at path.
// Returns an empty PinMap if the file does not exist.
func LoadPins(path string) (PinMap, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return PinMap{}, nil
	}
	if err != nil {
		return nil, err
	}
	var pins PinMap
	if err := json.Unmarshal(data, &pins); err != nil {
		return nil, err
	}
	return pins, nil
}

// PinKey adds or updates a pin for key with the given value and optional author.
func PinKey(pins PinMap, key, value, author string) PinMap {
	out := make(PinMap, len(pins)+1)
	for k, v := range pins {
		out[k] = v
	}
	out[key] = PinRecord{Value: value, PinnedAt: time.Now().UTC(), PinnedBy: author}
	return out
}

// UnpinKey removes a pin for key. Returns the updated map.
func UnpinKey(pins PinMap, key string) PinMap {
	out := make(PinMap, len(pins))
	for k, v := range pins {
		if k != key {
			out[k] = v
		}
	}
	return out
}

// ApplyPins overrides values in secrets with any pinned values.
func ApplyPins(secrets map[string]string, pins PinMap) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for key, rec := range pins {
		out[key] = rec.Value
	}
	return out
}
