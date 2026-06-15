package raft

import (
	"math/rand"
	"sync"
	"time"
)

type ConsensusModule struct {
	id    int
	peers []int

	// persistent state
	currentTerm int // logical clock — losing this lets stale candidates win elections
	votedFor    int // -1 means no vote; persisted so a crash can't cause double-voting
	log         []LogEntry

	// volatile state
	lastHeartbeat time.Time // reset on every heartbeat; drives election timeout
	commitIndex   int       // highest index known to be replicated on a majority
	lastApplied   int       // highest index applied to the state machine (lags commitIndex)

	// leader-only volatile state (reinitialized on election)
	nextIndex  []int // next log index to send to each peer (optimistic)
	matchIndex []int // highest log index confirmed replicated on each peer

	state CMState    // current role: Follower, Candidate, Leader, or Dead
	mu    sync.Mutex // protects all fields above
}

func NewConsensusModule(id int, peers []int) *ConsensusModule {

	cm := &ConsensusModule{
		peers:    peers,
		id:       id,
		votedFor: -1,
	}
	go cm.runElectionTimer()
	return cm
}
func (cm *ConsensusModule) runElectionTimer() {
	for {
		time.Sleep(10 * time.Millisecond)
		cm.mu.Lock()
		if cm.state == Leader || cm.state == Dead {
			cm.mu.Unlock()
			return
		}
		timeout := time.Duration(150+rand.Intn(151)) * time.Millisecond
		if time.Since(cm.lastHeartbeat) >= timeout {
			//Start a new election
		}
		cm.mu.Unlock()
	}
}
