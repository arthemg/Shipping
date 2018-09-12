package main

import (
	pb "github.com/arthemg/Shipping/consignment-service/proto/consignment"
	"context"
	"fmt"
	vesselProto "github.com/arthemg/Shipping/vessel-service/proto/vessel"
	micro "github.com/micro/go-micro"
	"log"
)

type Repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

type ConsignmentRepository struct {
	consignments []*pb.Consignment
}

func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) (*pb.Consignment, error){
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return  consignment, nil
}

func (repo *ConsignmentRepository) GetAll() []*pb.Consignment{
	return repo.consignments
}

type service struct {
	repo Repository
	VesselClient vesselProto.VesselServiceClient
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response)error{
	vesselResponse, err := s.VesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity: int32(len(req.Containers)),
	})
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
	if err != nil{
		return err
	}

	req.VesselId = vesselResponse.Vessel.Id


	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response)error  {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main(){
	repo := &ConsignmentRepository{}

	//Create a new service. Optionally include some options here
	srv := micro.NewService(
		//This name must mach the package name given in your protobuff definition
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
		)

	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())


	// Init will parse the command line flags
	srv.Init()

	//Register Handler
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})

	//Run server
	if err := srv.Run(); err != nil{
		fmt.Println(err)
	}
}