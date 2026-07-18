package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// FileStorage persists state to a single JSON file. Save writes to a temp
// file in the same directory and renames it into place, so a crash mid-write
// can never leave a corrupt or partially-written state file behind — the
// rename is atomic, the old file stays intact until it succeeds.
type FileStorage struct {
	path string
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{path: path}
}

func (fs *FileStorage) Save(state PersistentState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	dir := filepath.Dir(fs.path)
	tmp, err := os.CreateTemp(dir, ".tmp-state-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return os.Rename(tmpPath, fs.path)
}

func (fs *FileStorage) Load() (PersistentState, bool, error) {
	data, err := os.ReadFile(fs.path)
	if os.IsNotExist(err) {
		return PersistentState{}, false, nil
	}
	if err != nil {
		return PersistentState{}, false, err
	}

	var state PersistentState
	if err := json.Unmarshal(data, &state); err != nil {
		return PersistentState{}, false, err
	}
	return state, true, nil
}
