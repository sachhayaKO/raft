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

## 2026-06-26 — Sachin
- Completed: `startElection` method — increments term, transitions to Candidate, votes for self, launches parallel goroutines per peer
- Completed: vote counting logic — goroutines lock, check reply term, increment tally, transition to Leader on majority
- Learning: anonymous goroutines with `go func(arg) { }(arg)`, closure variable capture gotcha (why you pass peerId as argument), majority calculation (`len(peers)/2 + 1`)
- Key concepts: why you save `currentTerm` to a local var before launching goroutines (term may change by the time goroutines run), why votes starts at 1 (self-vote), why higher term in any reply means immediate revert to Follower
- Next: implement `RequestVote` RPC handler — the logic that receives and replies to vote requests

## 2026-07-04 — Sachin
- Completed: `RequestVote` RPC handler — term check, votedFor check, vote grant logic
- Completed: `leaderHeartbeat` goroutine — sends AppendEntries to all peers every 50ms, stops when no longer Leader
- Completed: wired `go cm.leaderHeartbeat()` inside `startElection` when majority votes received
- Learning: why `go` is required before function calls in goroutines (blocking vs non-blocking), mutex must be released before any RPC call, `reply.VoteGranted` defaults to false via Go zero values
- Key concepts: deadlock caused by holding mutex inside infinite loop, why leaderHeartbeat launches only at the moment of becoming Leader, `go` means fire-and-forget
- Next: implement `server.go` — wire up `net/rpc` so nodes can actually communicate over the network