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
	peerClients map[int]proto.RaftClient // gRPC client per peer, keyed by node ID
	cm          *ConsensusModule
	listener    net.Listener
	peerAddrs   map[int]string
	address     string
}

// NewServer creates a Server, dials all peers, and wires gRPC calls into the ConsensusModule.
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
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(err)
		}
		sv.peerClients[id] = proto.NewRaftClient(conn)
	}

	// wire outbound RPCs into the consensus module as plain function calls
	// keeps ConsensusModule free of gRPC imports
	cm.requestVoteFn = func(peerId int, args RequestVoteArgs) (*RequestVoteReply, error) {
		reply, err := sv.peerClients[peerId].RequestVote(context.Background(), &proto.RequestVoteArgs{
			Term:        int32(args.Term),
			CandidateId: int32(args.CandidateId),
		})
		if err != nil {
			return nil, err
		}
		return &RequestVoteReply{
			Term:        int(reply.Term),
			VoteGranted: reply.VoteGranted,
		}, nil
	}

	cm.appendEntriesFn = func(peerId int, args AppendEntriesArgs) (*AppendEntriesReply, error) {
		reply, err := sv.peerClients[peerId].AppendEntries(context.Background(), &proto.AppendEntriesArgs{
			Term:     int32(args.Term),
			LeaderId: int32(args.LeaderId),
		})
		if err != nil {
			return nil, err
		}
		return &AppendEntriesReply{
			Term:    int(reply.Term),
			Success: reply.Success,
		}, nil
	}

	return sv
}

// AppendEntries handles an incoming AppendEntries RPC, converts proto types, and calls into the algorithm.
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

// RequestVote handles an incoming RequestVote RPC, converts proto types, and calls into the algorithm.
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

// Start registers the gRPC server and begins serving on the listener.
func (sv *Server) Start() {
	grpcServer := grpc.NewServer()
	proto.RegisterRaftServer(grpcServer, sv)
	go grpcServer.Serve(sv.listener)
}
