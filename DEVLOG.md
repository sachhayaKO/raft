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
- Next Steps: building the NewConsensusModule constructorr. 