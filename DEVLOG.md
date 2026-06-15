# DEVLOG

## 2026-05-29 — Sachin
- Reading https://raft.github.io/raft.pdf
- Working on: learning goroutines and channels
- Blocked on: nothing
- Next: planning project timeline

## 2026-06-01 — Chris
- Working on: reading https://gobyexample.com and https://raft.github.io/raft.pdf
- Learning: goroutines, channels, mutexes, select statements
- Done: created basic file layout (raft/types.go, raft/consensus.go, raft/server.go), initialized go module, defined project data types
- Blocked on: none
- Next: finish project timeline and begin consensus logic

## 2026-06-04 - Sachin 
- Working on: Consensus struct creation 
- Learning typing and struct creation
- Next: begin coding consensus logic
- Important concepts learned: Each module has a perisistant state, volatile state, and leaders only have another volatile state. 
- Next Steps: building the NewConsensusModule constructor. 

## 2026-06-15 — Sachin
  - Completed: `ConsensusModule` struct with all persistent, volatile, and leader-only state fields
  - Completed: `NewConsensusModule` constructor with correct sentinel value for `votedFor` (-1)
  - Completed: `runElectionTimer` goroutine — sleeps 10ms per tick, checks state, randomizes timeout per
  iteration
  - Learning: Go structs and methods, pointer receivers, `sync.Mutex` patterns, `time.Timer` vs
  `time.Ticker`, infinite loops with `for`, goroutines with `go`
  - Key concepts: why `votedFor` must be persisted (double voting across crashes), why mutex must be
  unlocked before RPC calls, why election timeout must re-randomize each check
  - Next: implement `startElection` — increment term, vote for self, send `RequestVote` to all peers in
  parallel goroutines