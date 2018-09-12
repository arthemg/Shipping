//consignment-service/main.go
package main

import (
	"fmt"
	pb "github.com/arthemg/Shipping/consignment-service/proto/consignment"
	vesselProto "github.com/arthemg/Shipping/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
	"log"
	"os"
)

const (
	defaultHost = "localhost:27017"
)

func main() {
	// Database host from the environment variables
	host := os.Getenv("DB_HOST")

	if host == "" {
		host = defaultHost
	}
	session, err := CreateSession(host)

	defer session.Close()

	if err != nil {
		log.Panicf("Couldn't connect to datastore with host %s - %v", host, err)
	}

	srv := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())
	srv.Init()

	pb.RegisterShippingServiceHandler(srv.Server(), &service{session, vesselClient})

	//Run Server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
