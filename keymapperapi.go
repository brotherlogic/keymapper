package main

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/brotherlogic/keymapper/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// CONFIG - key storage location
	CONFIG = "github.com/brotherlogic/keymapper/config"
)

var (
	//KeySize - the print queue
	KeySize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "keymapper_keys",
		Help: "The size of the print queue",
	})
)

//Get a key
func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	keys := &pb.Keys{}
	err := s.Store.Load(ctx, CONFIG, keys)
	if err != nil {
		return nil, err
	}

	KeySize.Set(float64(len(keys.GetKeys())))

	for _, key := range keys.GetKeys() {
		if key.GetKey() == req.GetKey() {
			return &pb.GetResponse{Key: key}, nil
		}
	}

	return &pb.GetResponse{}, fmt.Errorf("Cannot find key: %v", req.GetKey())
}

//Set a key
func (s *Server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	keys := &pb.Keys{}
	err := s.Store.Load(ctx, CONFIG, keys)
	if err != nil {
		if status.Convert(err).Code() != codes.InvalidArgument {
			return nil, err
		}
	}

	KeySize.Set(float64(len(keys.GetKeys())))

	found := false
	for _, k := range keys.Keys {
		if k.GetKey() == req.GetKey() {
			k.Value = req.GetValue()
			found = true
		}
	}

	if !found {
		keys.Keys = append(keys.Keys, &pb.Key{Key: req.GetKey(), Value: req.GetValue()})
	}

	return &pb.SetResponse{}, s.Store.Save(ctx, CONFIG, keys)
}
