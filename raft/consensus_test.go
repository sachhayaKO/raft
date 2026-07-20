package raft

import (
	"testing"
	"time"
)

// newTestCluster spins up a 3-node in-process cluster with no real networking.
// Callbacks are wired to call peer methods directly so tests run without gRPC.
func newTestCluster(t *testing.T) map[int]*ConsensusModule {

	cm0 := NewConsensusModule(0, []int{1, 2})
	cm1 := NewConsensusModule(1, []int{0, 2})
	cm2 := NewConsensusModule(2, []int{0, 1})

	modules := map[int]*ConsensusModule{
		0: cm0,
		1: cm1,
		2: cm2,
	}
	for _, cm := range modules {
		cm.requestVoteFn = func(peerId int, args RequestVoteArgs) (*RequestVoteReply, error) {
			return modules[peerId].RequestVote(args)
		}
		cm.appendEntriesFn = func(peerId int, args AppendEntriesArgs) (*AppendEntriesReply, error) {
			return modules[peerId].AppendEntries(args)
		}
	}
	return modules
}

// TestElection verifies that exactly one leader is elected after a full election timeout.
func TestElection(t *testing.T) {
	modules := newTestCluster(t)
	// wait longer than the max election timeout (300ms) to guarantee an election has completed
	time.Sleep(500 * time.Millisecond)
	count := 0
	for _, cm := range modules {
		cm.mu.Lock()
		if cm.state == Leader {
			count += 1
		}
		cm.mu.Unlock()
	}
	if count != 1 {
		t.Errorf("expected 1 leader, got %d", count)
	}
}
