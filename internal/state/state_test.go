package state

import (
	"os"
	"testing"
)

func TestSaveLoad(t *testing.T) {
	filename := "test_state.bin"
	defer os.Remove(filename)

	original := &State{
		Version:        1,
		PublicKeyHex:   "02abc",
		StartRangeHex:  "1000",
		EndRangeHex:    "2000",
		TamePosX:       "dead",
		TamePosY:       "beef",
		TameDist:       "123",
		TotalTameSteps: 42,
		TotalWildSteps: 100500,
		Map: map[string]string{
			"a,b": "10",
			"c,d": "20",
		},
		FoundKey: "",
	}

	if err := original.Save(filename); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(filename)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Version != original.Version {
		t.Errorf("Version: got %d, want %d", loaded.Version, original.Version)
	}
	if loaded.PublicKeyHex != original.PublicKeyHex {
		t.Errorf("PublicKeyHex: got %s, want %s", loaded.PublicKeyHex, original.PublicKeyHex)
	}
	if len(loaded.Map) != len(original.Map) {
		t.Errorf("Map len: got %d, want %d", len(loaded.Map), len(original.Map))
	}
	for k, v := range original.Map {
		if loaded.Map[k] != v {
			t.Errorf("Map[%s]: got %s, want %s", k, loaded.Map[k], v)
		}
	}
}

func TestIsValid(t *testing.T) {
	s := &State{
		PublicKeyHex:  "abc",
		StartRangeHex: "123",
		EndRangeHex:   "456",
	}
	if !s.IsValid("abc", "123", "456") {
		t.Error("IsValid returned false for matching params")
	}
	if s.IsValid("abc", "123", "457") {
		t.Error("IsValid returned true for non-matching end range")
	}
}

func TestLoadNonexistent(t *testing.T) {
	_, err := Load("nonexistent.bin")
	if err == nil {
		t.Error("Load nonexistent.bin should return error")
	}
}
