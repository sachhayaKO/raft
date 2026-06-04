package raft

import "sync"

type ConsensusModule struct {
	id    int
	peers []int

	//persistent state

	/*
		This variable prevents past elections conflicting with present
		For example, if you restart with term 0, you'll grant votes to candidates
		you already voted against, and you might accept a leader from a term you already
		participated in.
	*/
	currentTerm int

	// Makes sure that if a node crashes, it remembers who they votes for
	// Nodes cannot vote for more than 1 candidate
	votedFor int

	//Log of all the LogEntries
	log []LogEntry

	//volatile state

	//Highest index where majority of nodes have the entry in their log
	commitIndex int

	//Highest term entry has been executed against the state machine
	lastApplied int

	//leader volatile state

	//Index of the next log entry the leader will send to each follower
	nextIndex []int

	//Highest log index the leader knows has been replicated on a follower
	matchIndex []int

	//Tracks whether this node is a Follower, Candidate or Leader
	state CMState

	//Lock for writing to the state
	mu sync.Mutex
}
