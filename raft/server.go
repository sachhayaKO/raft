package raft

import (
	"log"
	"net"
)

type Server struct {
	cm        *ConsensusModule
	listener  net.Listener
	peerAddrs map[int]string
	address   string
}

func NewServer(cm *ConsensusModule, peerAddrs map[int]string, address string) *Server {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	sv := &Server{
		cm:        cm,
		address:   address,
		peerAddrs: peerAddrs,
	}
	sv.listener = ln
	return sv
}
