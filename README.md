# Raft (Go)

A from-scratch implementation of the Raft consensus algorithm in Go. Building this to understand distributed systems and consensus from the ground up. The longer-term goal is to apply it to a real finance or healthcare application.

## Status

In progress. Leader election and gRPC networking layer implemented. Log replication and persistence coming next.

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
├── proto/
│   ├── raft.proto             - gRPC service and message definitions
│   ├── raft.pb.go             - generated message types
│   └── raft_grpc.pb.go        - generated server interface and client
├── storage/
│   ├── storage.go             - Storage interface, PersistentState/LogEntry types
│   └── file_storage.go        - JSON file-based Storage implementation
└── raft/
    ├── types.go               - internal data structures
    ├── consensus.go           - Raft algorithm (election, voting, heartbeat)
    ├── server.go              - gRPC server, client connections
    └── consensus_test.go      - election and state tests
```

## Running

Coming once the networking layer is complete.

## Design

See [DESIGN.md](./DESIGN.md) for decisions on quorum, persistence, timing, and networking.

## References

- [In Search of an Understandable Consensus Algorithm](https://raft.github.io/raft.pdf)
