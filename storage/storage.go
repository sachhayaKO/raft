package storage

// LogEntry mirrors raft.LogEntry, but keeps Command as opaque bytes instead
// of interface{}. Round-tripping interface{} through JSON loses the original
// concrete type (structs come back as map[string]interface{}, numbers as
// float64), so callers are responsible for encoding/decoding their own
// commands before handing them to Storage.
type LogEntry struct {
	Index   int
	Term    int
	Command []byte
}

// PersistentState is the subset of Raft state that must survive a crash:
// CurrentTerm and VotedFor prevent double-voting on restart, Log holds
// committed entries.
type PersistentState struct {
	CurrentTerm int
	VotedFor    int
	Log         []LogEntry
}

// Storage persists and loads Raft's durable state.
type Storage interface {
	// Save must not return until state is durable on disk.
	Save(state PersistentState) error
	// Load returns the last saved state. found is false if nothing has
	// been saved yet (first run on a fresh node).
	Load() (state PersistentState, found bool, err error)
}
