package raft

import (
	"context"
	"log"
	"net"
	"raftproject/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	proto.UnimplementedRaftServer
	peerClients map[int]proto.RaftClient
	cm          *ConsensusModule
	listener    net.Listener
	peerAddrs   map[int]string
	address     string
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
	sv.peerClients = make(map[int]proto.RaftClient)
	for id, addr := range peerAddrs {
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))		if err != nil {
			log.Fatal(err)
		}
		sv.peerClients[id] = proto.NewRaftClient(conn)
	}
	return sv
}
func (sv *Server) AppendEntries(ctx context.Context, req *proto.AppendEntriesArgs) (*proto.AppendEntriesReply, error) {
	internalArgs := AppendEntriesArgs{
		Term:     int(req.Term),
		LeaderId: int(req.LeaderId),
	}
	reply, err := sv.cm.AppendEntries(internalArgs)
	if err != nil {
		return nil, err
	}
	return &proto.AppendEntriesReply{
		Term:    int32(reply.Term),
		Success: reply.Success,
	}, nil

}

func (sv *Server) RequestVote(ctx context.Context, req *proto.RequestVoteArgs) (*proto.RequestVoteReply, error) {
	internalArgs := RequestVoteArgs{
		Term:        int(req.Term),
		CandidateId: int(req.CandidateId),
	}
	reply, err := sv.cm.RequestVote(internalArgs)
	if err != nil {
		return nil, err
	}
	return &proto.RequestVoteReply{
		Term:        int32(reply.Term),
		VoteGranted: reply.VoteGranted,
	}, nil
}

func (sv *Server) Start() {
	grpcServer := grpc.NewServer()
	proto.RegisterRaftServer(grpcServer, sv)
	go grpcServer.Serve(sv.listener)
}
