package envfile

import (
	"errors"
	"testing"
)

func TestValidate_ValidKeys(t *testing.T) {
	secrets := map[string]string{
		"FOO":       "bar",
		"_PRIVATE":  "val",
		"MY_VAR_99": "val",
	}
	if err := Validate(secrets); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidate_InvalidKeyStartsWithDigit(t *testing.T) {
	secrets := map[string]string{
		"1INVALID": "val",
	}
	err := Validate(secrets)
	if err == nil {
		t.Fatal("expected error for key starting with digit")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.InvalidKeys) != 1 || ve.InvalidKeys[0] != "1INVALID" {
		t.Errorf("unexpected invalid keys: %v", ve.InvalidKeys)
	}
}

func TestValidate_InvalidKeyWithHyphen(t *testing.T) {
	secrets := map[string]string{
		"MY-VAR": "val",
	}
	err := Validate(secrets)
	if err == nil {
		t.Fatal("expected error for key with hyphen")
	}
}

func TestValidate_EmptyMap(t *testing.T) {
	if err := Validate(map[string]string{}); err != nil {
		t.Errorf("expected no error for empty map, got: %v", err)
	}
}

func TestValidate_MixedKeys(t *testing.T) {
	secrets := map[string]string{
		"GOOD_KEY": "val",
		"bad key":  "val",
	}
	err := Validate(secrets)
	if err == nil {
		t.Fatal("expected error for key with space")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.InvalidKeys) != 1 {
		t.Errorf("expected 1 invalid key, got %d: %v", len(ve.InvalidKeys), ve.InvalidKeys)
	}
}
