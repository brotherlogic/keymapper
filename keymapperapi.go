package main

import (
	"fmt"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/keymapper/proto"
)

const (
	// CONFIG - key storage location
	CONFIG = "/github.com/brotherlogic/keymapper/config"
)

//Get a key
func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	data, _, err := s.KSclient.Read(ctx, CONFIG, &pb.Keys{})
	if err != nil {
		return nil, err
	}
	keys := data.(*pb.Keys)

	for _, key := range keys.GetKeys() {
		if key.GetKey() == req.GetKey() {
			return &pb.GetResponse{Key: key}, nil
		}
	}

	return nil, fmt.Errorf("Cannot find key: %v", req.GetKey())
}

//Set a key
func (s *Server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	data, _, err := s.KSclient.Read(ctx, CONFIG, &pb.Keys{})
	if err != nil {
		return nil, err
	}
	keys := data.(*pb.Keys)

	keys.Keys = append(keys.Keys, &pb.Key{Key: req.GetKey(), Value: req.GetValue()})

	return nil, s.KSclient.Save(ctx, CONFIG, keys)
}
