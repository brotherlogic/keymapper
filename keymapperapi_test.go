package main

import (
	"testing"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/keymapper/proto"
	keystoreclient "github.com/brotherlogic/keystore/client"
)

func InitTest() *Server {
	s := Init()
	s.SkipLog = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient("./testing")
	s.GoServer.KSclient.Save(context.Background(), CONFIG, &pb.Keys{})
	return s
}

func TestBasic(t *testing.T) {
	s := InitTest()

	_, err := s.Set(context.Background(), &pb.SetRequest{Key: "testkey", Value: "donkey"})
	if err != nil {
		t.Errorf("Bad request: %v", err)
	}

	resp, err := s.Get(context.Background(), &pb.GetRequest{Key: "testkey"})
	if err != nil {
		t.Errorf("Bad get request: %v", err)
	}

	if resp.GetKey().GetValue() != "donkey" {
		t.Errorf("Bad return: %v", resp)
	}
}

func TestBadSet(t *testing.T) {
	s := InitTest()
	s.GoServer.KSclient.Fail = true

	r, err := s.Set(context.Background(), &pb.SetRequest{Key: "testkey", Value: "donkey"})
	if err == nil {
		t.Errorf("Bad request: %v", r)
	}
}

func TestEmptySet(t *testing.T) {
	s := InitTest()
	s.GoServer.KSclient = *keystoreclient.GetTestClient("./sptest")

	r, err := s.Set(context.Background(), &pb.SetRequest{Key: "testkey", Value: "donkey"})
	if err != nil {
		t.Errorf("Bad request: %v", r)
	}
}

func TestBadGet(t *testing.T) {
	s := InitTest()
	s.GoServer.KSclient.Fail = true
	resp, err := s.Get(context.Background(), &pb.GetRequest{Key: "testkey"})
	if err == nil {
		t.Errorf("Bad get request: %v", resp)
	}
}

func TestMissingKey(t *testing.T) {
	s := InitTest()

	_, err := s.Set(context.Background(), &pb.SetRequest{Key: "testkey", Value: "donkey"})
	if err != nil {
		t.Errorf("Bad request: %v", err)
	}

	resp, err := s.Get(context.Background(), &pb.GetRequest{Key: "differentkey"})
	if err == nil {
		t.Errorf("Bad get request: %v", resp)
	}
}
