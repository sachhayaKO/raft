package raft

// CMState represents the state of a Consensus Module.
type CMState int

const (
	Follower CMState = iota
	Candidate
	Leader
	Dead
)

// LogEntry is a single entry in the replicated log.
type LogEntry struct {
	Index   int
	Term    int
	Command interface{}
}

// RequestVoteArgs is the arguments for a RequestVote RPC.
type RequestVoteArgs struct {
	Term        int // candidate's term
	CandidateId int // candidate requesting vote
}

// RequestVoteReply is the reply to a RequestVote RPC.
type RequestVoteReply struct {
	Term        int  // currentTerm, for candidate to update itself
	VoteGranted bool // true means candidate received vote
}

// AppendEntriesArgs is the arguments for an AppendEntries RPC.
type AppendEntriesArgs struct {
	Term     int // leader's term
	LeaderId int // so follower can redirect clients
}

// AppendEntriesReply is the reply to an AppendEntries RPC.
type AppendEntriesReply struct {
	Term    int  // currentTerm, for leader to update itself
	Success bool // true if follower accepted heartbeat
}
