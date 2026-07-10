package raft

import (
	"context"
	"log"
	"net"
	"raftproject/proto"

	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedRaftServer
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
func (sv *Server) AppendEntries(ctx context.Context, req *proto.AppendEntriesArgs) (*proto.AppendEntriesReply, error) {
	internalArgs := AppendEntriesArgs{
		Term:     int(req.Term),
		LeaderId: int(req.LeaderId),
	}
	reply, err := sv.cm.AppendEntries(internalArgs)
	return &proto.AppendEntriesReply{
		Term:    int32(reply.Term),
		Success: bool(reply.Success),
	}, err
}

func (sv *Server) RequestVote(ctx context.Context, req *proto.RequestVoteArgs) (*proto.RequestVoteReply, error) {
	internalArgs := RequestVoteArgs{
		Term:        int(req.Term),
		CandidateId: int(req.CandidateId),
	}
	reply, err := sv.cm.RequestVote(internalArgs)

	return &proto.RequestVoteReply{
		Term:        int32(reply.Term),
		VoteGranted: reply.VoteGranted,
	}, err
}
func (sv *Server) Start() {
	grpcServer := grpc.NewServer()
	proto.RegisterRaftServer(grpcServer, sv)
	grpcServer.Serve(sv.listener)
}
