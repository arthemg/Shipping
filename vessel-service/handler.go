package main

import (
	"context"
	"gopkg.in/mgo.v2"
	pb "github.com/arthemg/Shipping/vessel-service/proto/vessel"
)

type service struct {
	session *mgo.Session
}

func(s *service)GetRepo() Repository{
	return &VesselRepository{s.session.Clone()}
}

func (s *service) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error{
	defer s.GetRepo().Close()
	vessel, err := s.GetRepo().FindAvailable(req)
	if err != nil{
		return err
	}
	// Set the vessel as part of the response message type
	res.Vessel = vessel
	return nil
}

func(s *service)Create(ctx context.Context, req *pb.Vessel, res *pb.Response) error{
	defer s.GetRepo().Close()
	if err := s.GetRepo().Create(req); err != nil {
		return  err
	}
	res.Vessel = req
	res.Created = true
	return nil
}