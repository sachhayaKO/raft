# DEVLOG

## 2026-06-28 — Sachin
- Completed: `requestVote` RPC handler — acquires lock, steps down to Follower on higher term, sets `reply.Term`, grants vote if `votedFor == -1` or `votedFor == candidateId`, denies otherwise
- Fixed: `startElection` bugs — moved `mu.Lock()` before `currentTerm` increment, replaced `cm.currentTerm` with `savedTerm` in goroutine comparison
- Learning: `defer` for automatic unlock, Go `else` must be on same line as closing `}`, pointer receivers vs value receivers for RPC reply
- Key concepts: why `savedTerm` is needed (term may change before goroutine runs), why reply is a pointer (caller needs to see your changes), why higher term always triggers step-down
- Next: implement `AppendEntries` RPC handler (heartbeat receiver)

## 2026-06-26 — Sachin
- Completed: `startElection` method — increments term, transitions to Candidate, votes for self, launches parallel goroutines per peer
- Completed: vote counting logic — goroutines lock, check reply term, increment tally, transition to Leader on majority
- Learning: anonymous goroutines with `go func(arg) { }(arg)`, closure variable capture gotcha (why you pass peerId as argument), majority calculation (`len(peers)/2 + 1`)
- Key concepts: why you save `currentTerm` to a local var before launching goroutines (term may change by the time goroutines run), why votes starts at 1 (self-vote), why higher term in any reply means immediate revert to Follower
- Next: implement `RequestVote` RPC handler — the logic that receives and replies to vote requests

## 2026-06-15 — Sachin
- Completed: `ConsensusModule` struct with all persistent, volatile, and leader-only state fields
- Completed: `NewConsensusModule` constructor with correct sentinel value for `votedFor` (-1)
- Completed: `runElectionTimer` goroutine — sleeps 10ms per tick, checks state, randomizes timeout per iteration
- Learning: Go structs and methods, pointer receivers, `sync.Mutex` patterns, `time.Timer` vs `time.Ticker`, infinite loops with `for`, goroutines with `go`
- Key concepts: why `votedFor` must be persisted (double voting across crashes), why mutex must be unlocked before RPC calls, why election timeout must re-randomize each check
- Next: implement `startElection` — increment term, vote for self, send `RequestVote` to all peers in parallel goroutines

## 2026-06-04 — Sachin
- Completed: `ConsensusModule` struct definition
- Learning: typing and struct creation in Go
- Key concepts: each module has persistent state, volatile state, and leaders have an additional volatile state
- Next: build the `NewConsensusModule` constructor

## 2026-06-01 — Chris
- Completed: basic file layout (`raft/types.go`, `raft/consensus.go`, `raft/server.go`), initialized Go module, defined project data types
- Learning: goroutines, channels, mutexes, select statements
- Next: finish project timeline and begin consensus logic

## 2026-05-29 — Sachin
- Learning: goroutines and channels
- Next: planning project timeline
