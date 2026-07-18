package raft

import (
	"testing"
)

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

func TestElection(t *testing.T) {
	modules := newTestCluster(t)

}
