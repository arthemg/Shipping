package main

import (
	"context"
	"fmt"
	pb "github.com/arthemg/Shipping/vessel-service/proto/vessel"
	"github.com/go-openapi/errors"
	"github.com/micro/go-micro"
)
type Repository interface {
	FindAvailable(*pb.Specification)(*pb.Vessel, error)
}

type VesselRepository struct {
	vessels []*pb.Vessel
}

type service struct {
	repo Repository
}

func (repo *VesselRepository) FindAvailable(spec *pb.Specification)(*pb.Vessel, error){
	for _, vessel := range repo.vessels {
		if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.MaxWeight {
			return vessel, nil
		}
	}
	return nil, errors.NotFound("No vessel found for this spec")
}

func (s *service) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error{
	vessel, err := s.repo.FindAvailable(req)
	if err != nil{
		return err
	}

	res.Vessel = vessel
	return nil
}

func main()  {
	vessels := []*pb.Vessel{
		&pb.Vessel{Id:"vessel001", Name:"Boaty McBoatface", MaxWeight:200000, Capacity:500},
	}

	repo :=&VesselRepository{vessels}

	srv:= micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
		)

	srv.Init()

	pb.RegisterVesselServiceHandler(srv.Server(), &service{repo})

	if err := srv.Run(); err != nil{
		fmt.Println(err)
	}
}



