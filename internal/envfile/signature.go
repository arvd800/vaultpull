package envfile

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// SignatureRecord holds an HMAC signature for a secrets map.
type SignatureRecord struct {
	Signature string    `json:"signature"`
	SignedAt  time.Time `json:"signed_at"`
	KeyCount  int       `json:"key_count"`
}

// SignMap computes an HMAC-SHA256 signature over the sorted key=value pairs.
func SignMap(secrets map[string]string, passphrase string) (string, error) {
	if passphrase == "" {
		return "", fmt.Errorf("passphrase must not be empty")
	}
	payload := marshalSorted(secrets)
	mac := hmac.New(sha256.New, []byte(passphrase))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// VerifySignature checks that a secrets map matches a stored SignatureRecord.
func VerifySignature(secrets map[string]string, passphrase string, record SignatureRecord) error {
	got, err := SignMap(secrets, passphrase)
	if err != nil {
		return err
	}
	if !hmac.Equal([]byte(got), []byte(record.Signature)) {
		return fmt.Errorf("signature mismatch: secrets may have been tampered with")
	}
	return nil
}

// SaveSignature writes a SignatureRecord to path.
func SaveSignature(path string, record SignatureRecord) error {
	if path == "" {
		return nil
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadSignature reads a SignatureRecord from path.
func LoadSignature(path string) (SignatureRecord, error) {
	var record SignatureRecord
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return record, nil
		}
		return record, err
	}
	err = json.Unmarshal(data, &record)
	return record, err
}

func marshalSorted(secrets map[string]string) string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var out string
	for _, k := range keys {
		out += k + "=" + secrets[k] + "\n"
	}
	return out
}
