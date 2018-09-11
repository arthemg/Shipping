package main

import (
	pb "github.com/arthemg/Shipping/consignment-service/proto/consignment"
	"context"
	"fmt"
	micro "github.com/micro/go-micro"
)
const(
	port = ":50051"
)

type IRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error){
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return  consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment{
	return repo.consignments
}

type service struct {
	repo IRepository
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response)error{
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
	repo := &Repository{}

	//Create a new service. Optionally include some options here
	srv := micro.NewService(
		//This name must mach the package name given in your protobuff definition
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
		)

	// Init will parse the command line flags
	srv.Init()

	//Register Handler
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo})

	//Run server
	if err := srv.Run(); err != nil{
		fmt.Println(err)
	}
}