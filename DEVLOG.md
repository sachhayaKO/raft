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

## 2026-07-06 — Sachin
- Completed: `Server` struct — `cm`, `listener`, `peerAddrs`, `address`
- Completed: `NewServer` constructor — takes cm, peer addresses, and own address; starts TCP listener with `net.Listen`
- Learning: Go imports syntax, `net.Listener` type, error handling with `log.Fatal`, difference between receivers and regular functions
- Key concepts: server owns its own cm only (not peers'), outbound needs peer address map, `NewServer` receives cm from caller rather than creating it
- Next: implement `Start` method — accept connections and register RPC handlers

## 2026-07-07 — Sachin
- Decision: switched from net/rpc to gRPC for networking layer. Better for resume, cross-language, has TLS and observability built in. Plan is to finish with net/rpc understanding first then migrate.
- Completed: proto/raft.proto — defined Raft service with RequestVote and AppendEntries RPCs, all four message types matching types.go
- Completed: ran protoc to generate proto/raft.pb.go and proto/raft_grpc.pb.go
- Learning: proto3 syntax, message field numbers, go_package option, protoc flags (source_relative to control output path)
- Key concepts: proto file defines the RPC contract, protoc generates the server interface and client stub, gRPC signatures use context.Context and return values instead of out parameters
- Next: decide whether ConsensusModule or Server implements the gRPC interface, then update method signatures to match generated interface

## 2026-07-08 — Sachin
- Decision: Server implements the gRPC interface, ConsensusModule stays pure algorithm with no networking imports
- Completed: updated RequestVote signature to return (*RequestVoteReply, error) instead of using out parameter
- Completed: AppendEntries handler on ConsensusModule — rejects stale leaders, resets lastHeartbeat on valid heartbeat
- Completed: wired go cm.startElection() into runElectionTimer, fixed deadlock (was calling directly while holding mutex)
- Completed: leaderHeartbeat reads currentTerm and id into locals before unlocking, RPC call still placeholder
- Completed: startElection now builds RequestVoteArgs with savedTerm and calls cm.RequestVote directly (will be replaced with gRPC client call)
- Next: implement Server methods for gRPC interface, wire up gRPC server in Start method

## 2026-07-12 — Sachin
- Completed: Server.AppendEntries and Server.RequestVote — convert proto types to internal types, call cm, convert reply back
- Completed: Start method — creates gRPC server, registers RaftServer, serves on listener in goroutine
- Completed: peerClients map on Server — dials all peers on construction using grpc.NewClient, stores proto.RaftClient per peer ID
- Learning: gRPC client connections are persistent and stored, not created per call. grpc.Dial is deprecated, use grpc.NewClient. Must initialize map with make() before writing to it
- Key concepts: nil pointer panic if you return reply fields when err != nil, insecure.NewCredentials() needed for local connections without TLS
- Next: replace placeholder comments in startElection and leaderHeartbeat with real gRPC client calls using peerClients

## 2026-07-14 — Sachin
- Completed: wired requestVoteFn and appendEntriesFn callbacks into startElection and leaderHeartbeat
- Completed: leaderHeartbeat now handles higher term in reply — steps down to Follower immediately
- Completed: server.go comments added
- Key concepts: callback pattern keeps ConsensusModule free of gRPC imports, higher term in any reply triggers immediate stepdown
- Next: write tests in consensus_test.go, then move to log replication

## 2026-07-17 — Sachin
- Completed: wrote consensus_test.go — newTestCluster helper, TestElection asserts exactly 1 leader after 500ms
- Completed: fixed two algorithm bugs caught by the test (see below)
- Bug 1: runElectionTimer did not check Candidate state — kept firing elections while a node was already running one. Fixed by adding Candidate to the exit condition
- Bug 2 (root cause of 3 leaders): RequestVote and AppendEntries updated currentTerm when seeing a higher term but never set state = Follower. A Leader that received a RequestVote from a higher-term candidate would update its term, vote for the new candidate, but stay as Leader. Its heartbeats now used the new term and succeeded — so it never stepped down. All 3 nodes ended up as permanent leaders at the same term. Fixed by adding cm.state = Follower in both handlers when a higher term is seen
- Key concept: "higher term always wins" is not just about updating currentTerm — it must also trigger a state transition. Every RPC handler must enforce this
- Next: log replication — update types.go with missing fields, implement real AppendEntries logic

## 2026-07-19 — Chris
- Completed: `storage` package — `Storage` interface (`Save`/`Load`), `PersistentState`/`LogEntry` types
- Completed: `FileStorage` — JSON to a temp file, then `os.Rename` into place so a crash mid-write can't corrupt the state file
- Learning: `LogEntry.Command` is `[]byte` not `interface{}` — JSON round-tripping `interface{}` loses the concrete type
- Key concept: persistence lives in its own package, not `ConsensusModule` — testable standalone, no conflicts with in-progress algorithm work
- Next: write `file_storage_test.go`, then wire `Storage` into `ConsensusModule`