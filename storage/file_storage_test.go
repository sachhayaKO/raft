package storage

import (
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")
	fs := NewFileStorage(path)

	var entries []LogEntry
	for i := 1; i <= 2; i++ {
		entries = append(entries, LogEntry{
			Index:   i,
			Term:    i,
			Command: []byte("cmd" + strconv.Itoa(i)),
		})
	}

	state := PersistentState{
		CurrentTerm: 3,
		VotedFor:    1,
		Log:         entries,
	}

	if err := fs.Save(state); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, found, err := fs.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if found == false {
		t.Fatalf("expected found to be true")
	}
	if !reflect.DeepEqual(state, loaded) {
		t.Errorf("loaded state doesn't match saved state\nwant: %+v\ngot:  %+v", state, loaded)
	}
}

func TestLoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	fs := NewFileStorage(dir + "/does-not-exist.json")

	state, found, err := fs.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if found {
		t.Errorf("expected found to be false")
	}
	if state.CurrentTerm != 0 || state.VotedFor != 0 {
		t.Errorf("expected zero-value state, got %+v", state)
	}
}

func TestSaveOverwrites(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")
	fs := NewFileStorage(path)

	first := PersistentState{CurrentTerm: 1, VotedFor: 0}
	second := PersistentState{CurrentTerm: 5, VotedFor: 2}

	fs.Save(first)
	fs.Save(second)

	loaded, found, err := fs.Load()
	if err != nil || !found {
		t.Fatalf("Load failed: found=%v err=%v", found, err)
	}
	if loaded.CurrentTerm != 5 || loaded.VotedFor != 2 {
		t.Errorf("expected latest save to win, got %+v", loaded)
	}
}

// A handful of saves in a row shouldn't leave any of the .tmp-state-* files
// behind in the directory - only the final renamed file should exist.
func TestSaveDoesNotLeaveTempFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")
	fs := NewFileStorage(path)

	for i := 0; i < 5; i++ {
		err := fs.Save(PersistentState{CurrentTerm: i})
		if err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected exactly 1 file in dir after saves, found %d", len(entries))
	}
}
