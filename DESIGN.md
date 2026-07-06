# Design Decisions

Notes on the choices that actually matter for correctness and why we made them. Updated as we go.

---

## Architecture

```
┌──────────────────────────┐
│        server.go         │
│   networking, RPC, I/O   │
└────────────┬─────────────┘
             │
┌────────────▼─────────────┐
│       consensus.go       │
│   pure algorithm, no I/O │
└──────────────────────────┘
```

`ConsensusModule` is the algorithm. `Server` is the network layer. They don't know each other's internals.

The reason for this split is testability. In `consensus_test.go` we can run a full multi-node cluster in a single process by having modules call each other's methods directly, no real sockets needed. Without this separation, every test would require real networking and timing, which makes bugs much harder to reproduce.

---

## State Machine

```mermaid
flowchart LR
    F[Follower] -->|timeout| C[Candidate]
    C -->|majority votes| L[Leader]
    C -->|higher term| F
    L -->|higher term| F
    C -->|split vote| C
```

---

## Election Flow

```mermaid
sequenceDiagram
    participant A as Node A
    participant B as Node B
    participant C as Node C

    Note over A: timeout fires
    A->>B: RequestVote term=2
    A->>C: RequestVote term=2
    B-->>A: vote granted
    C-->>A: vote granted
    Note over A: majority, now Leader
    A->>B: AppendEntries
    A->>C: AppendEntries
    Note over B,C: timers reset
```

---

## Why strict majority (`n/2 + 1`)

A candidate needs more than half the cluster, not just half.

Two majorities in the same cluster always overlap by at least one node, and that node can only vote once per term. So you can never have two candidates both reach majority at the same time. Without this guarantee you could elect two leaders in the same term, which means two different entries committed at the same log index. That's the one thing Raft cannot allow.

A 5-node cluster tolerates 2 failures. Adding nodes improves fault tolerance but also raises the number of votes needed, which slows commits.

---

## What gets persisted

Only three fields are written to disk before replying to any RPC: `currentTerm`, `votedFor`, `log[]`.

Each one protects something specific:

- `currentTerm` — if a node restarts at term 0 it will accept messages from leaders it should reject and vote in terms it already participated in. Terms are how the whole cluster agrees on who is authoritative.
- `votedFor` — without this, a node that crashes after voting can restart and vote again in the same term. In a small cluster that can give two candidates a majority at once.
- `log[]` — committed entries live here. Losing this means losing committed data.

`commitIndex` and `lastApplied` are not persisted because they can be reconstructed by replaying the log on restart.

The tradeoff is that every RPC reply requires a disk write first. This is the main performance bottleneck in Raft.

---

## Timing: heartbeat vs election timeout

Heartbeat: **50ms**. Election timeout: **150–300ms** (randomized per node).

The heartbeat interval needs to be well below the election timeout. If a heartbeat is slow and the timeout fires, you get an unnecessary election even though the leader is still alive. The ratio matters more than the absolute values.

The randomization within 150–300ms is what prevents split votes. If every node had the same timeout they would all call elections at the same time every time the leader died, and you'd get repeated splits with no winner. Staggering the timeouts means one node almost always fires first and wins before others wake up.

---

## Higher term always wins

Whenever any node sees a term higher than its own in any RPC, it immediately steps down to Follower and updates its term. No exceptions.

This is how Raft prevents stale leaders. A leader that got partitioned from the cluster could come back thinking it's still in charge. If the rest of the cluster elected a new leader in a higher term, the old leader needs to step down the moment it sees that term. Without this rule you can have two nodes both thinking they're leader, which breaks everything.

This applies in both directions — request and response, any RPC type.

---

## Networking: `net/rpc`

Using Go's standard library `net/rpc` instead of gRPC.

The handler signature Raft needs is `func (s *T) Method(args T, reply *T) error`, which is exactly what `net/rpc` expects. Zero dependencies, no adapter code.

For a production system this would be gRPC — it's cross-language, has TLS, retry, and observability built in. `net/rpc` is Go-only and effectively deprecated upstream. Fine for this project, not for prod.

---

## References

- [Raft paper — Ongaro & Ousterhout (2014)](https://raft.github.io/raft.pdf)
- [Go net/rpc](https://pkg.go.dev/net/rpc)
