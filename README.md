# Raft (Go)

A from-scratch implementation of the Raft consensus algorithm in Go. Building this to understand distributed systems and consensus from the ground up. The longer-term goal is to apply it to a real finance or healthcare application.

## Status

In progress. Leader election implemented. Log replication and persistence coming next.

## What is this

Raft is a consensus algorithm that lets a cluster of servers agree on a shared log of commands even when some servers fail. This implementation follows the original paper by Diego Ongaro and John Ousterhout.

Covers:
- Leader election
- Log replication (in progress)
- Safety guarantees
- Persistence (planned)

## Layout

```
raft/
├── go.mod
└── raft/
    ├── types.go           - data structures, RPC message types
    ├── consensus.go       - Raft algorithm (election, voting, heartbeat)
    ├── server.go          - networking and RPC
    └── consensus_test.go  - election and state tests
```

## Running

Coming once the networking layer is complete.

## Design

See [DESIGN.md](./DESIGN.md) for decisions on quorum, persistence, timing, and networking.

## References

- [In Search of an Understandable Consensus Algorithm](https://raft.github.io/raft.pdf)
