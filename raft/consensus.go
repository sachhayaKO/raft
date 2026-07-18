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

	//Server functions we need to call to access the peerAddrs without exposing
	requestVoteFn   func(peerId int, args RequestVoteArgs) (*RequestVoteReply, error)
	appendEntriesFn func(peerId int, args AppendEntriesArgs) (*AppendEntriesReply, error)
}

// NewConsensusModule creates a new ConsensusModule and immediately starts its election timer.
func NewConsensusModule(id int, peers []int) *ConsensusModule {
	cm := &ConsensusModule{
		peers:    peers,
		id:       id,
		votedFor: -1,
	}
	go cm.runElectionTimer()
	return cm
}

// runElectionTimer runs in a goroutine for the lifetime of a non-leader node.
// It fires an election if no heartbeat arrives within a randomized timeout.
func (cm *ConsensusModule) runElectionTimer() {
	for {
		time.Sleep(10 * time.Millisecond)
		cm.mu.Lock()
		if cm.state == Leader || cm.state == Dead || cm.state == Candidate {
			cm.mu.Unlock()
			return
		}
		timeout := time.Duration(150+rand.Intn(151)) * time.Millisecond
		if time.Since(cm.lastHeartbeat) >= timeout {
			go cm.startElection()
		}
		cm.mu.Unlock()
	}
}

// startElection transitions the node to Candidate, votes for itself,
// and sends RequestVote RPCs to all peers in parallel.
func (cm *ConsensusModule) startElection() {
	cm.mu.Lock()
	cm.currentTerm += 1
	cm.state = Candidate
	cm.votedFor = cm.id
	savedTerm := cm.currentTerm
	cm.mu.Unlock()

	votes := 1 // counts self-vote

	for _, peerId := range cm.peers {
		go func(peerId int) {

			args := RequestVoteArgs{
				Term:        savedTerm,
				CandidateId: cm.id,
			}
			reply, err := cm.requestVoteFn(peerId, args)

			if err != nil {
				return
			}

			cm.mu.Lock()
			// higher term in reply means we're stale — revert to follower
			if reply.Term > cm.currentTerm {
				cm.state = Follower
				cm.mu.Unlock()
				return
			}
			votes++
			if votes >= len(cm.peers)/2+1 && cm.state == Candidate {
				cm.state = Leader
				go cm.leaderHeartbeat()
			}
			cm.mu.Unlock()
		}(peerId)
	}
}

// RequestVote handles an incoming RequestVote RPC from a candidate.
func (cm *ConsensusModule) RequestVote(args RequestVoteArgs) (*RequestVoteReply, error) {
	cm.mu.Lock()
	reply := &RequestVoteReply{}
	// higher term forces us to update and clear any prior vote
	if args.Term > cm.currentTerm {
		cm.currentTerm = args.Term
		cm.state = Follower
		cm.votedFor = -1
	}
	// grant vote if we haven't voted yet this term (or already voted for this candidate)
	if args.Term == cm.currentTerm && (cm.votedFor == -1 || cm.votedFor == args.CandidateId) {
		cm.votedFor = args.CandidateId
		reply.VoteGranted = true
	}
	reply.Term = cm.currentTerm
	cm.mu.Unlock()
	return reply, nil
}

func (cm *ConsensusModule) AppendEntries(args AppendEntriesArgs) (*AppendEntriesReply, error) {
	cm.mu.Lock()
	reply := &AppendEntriesReply{}
	if args.Term < cm.currentTerm {
		reply.Success = false
		reply.Term = cm.currentTerm
	} else {
		cm.state = Follower
		reply.Success = true
		cm.lastHeartbeat = time.Now()
	}
	cm.mu.Unlock()
	return reply, nil
}

// leaderHeartbeat runs in a goroutine while this node is Leader.
// It sends AppendEntries RPCs to all peers every 50ms to suppress new elections.
func (cm *ConsensusModule) leaderHeartbeat() {
	for {
		time.Sleep(50 * time.Millisecond)
		cm.mu.Lock()
		if cm.state != Leader || cm.state == Dead {
			cm.mu.Unlock()
			return
		}
		term := cm.currentTerm
		id := cm.id
		cm.mu.Unlock()
		for _, peerId := range cm.peers {
			args := AppendEntriesArgs{
				Term:     term,
				LeaderId: id,
			}
			reply, err := cm.appendEntriesFn(peerId, args)
			if err != nil {
				continue
			}
			if reply.Term > term {
				cm.mu.Lock()
				cm.state = Follower
				cm.currentTerm = reply.Term
				cm.mu.Unlock()
				return
			}
		}
	}
}
