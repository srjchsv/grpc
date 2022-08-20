package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"github.com/srjcshv/grpc/pb"
	"github.com/srjcshv/grpc/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	serversAddress := flag.String("address", "", "the servers address")
	maxPriceUsd := flag.Float64("price", 3000, "max price")
	minCpuCores := flag.Uint64("cores", 4, "min cpu cores")
	minCpuGhz := flag.Float64("ghz", 2.5, "min cpu ghz")
	minRam := flag.Uint64("ram", 8, "min ram")
	flag.Parse()
	
	log.Printf("dial server %v", *serversAddress)

	conn, err := grpc.Dial(*serversAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	laptopClient := pb.NewLaptopServiceClient(conn)
	for i := 0; i < 10; i++ {
		createLaptop(laptopClient)
	}
	filter := &pb.Filter{
		MaxPriceUsd: *maxPriceUsd,
		MinCpuCores: uint32(*minCpuCores),
		MinCpuGhz:   *minCpuGhz,
		MinRam:      &pb.Memory{Value: *minRam, Unit: pb.Memory_GIGABYTE},
	}

	SearchLaptop(laptopClient, filter)

}

func SearchLaptop(laptopClient pb.LaptopServiceClient, filter *pb.Filter) {
	log.Printf("search filter: %v", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatal("cannot search laptop: ", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("cannot recieve response ", err)
		}

		laptop := res.GetLaptop()
		log.Print("-found: ", laptop.GetId())
		log.Print(" +brand: ", laptop.GetBrand())
		log.Print(" +name: ", laptop.GetName())
		log.Print(" +cpu cores: ", laptop.GetCpu().GetNumberCores())
		log.Print(" +price:", laptop.GetPriceUsd())
	}
}

func createLaptop(laptopClient pb.LaptopServiceClient) {
	laptop := sample.NewLaptop()
	laptop.Id = ""
	req := &pb.CreateLaptopRequest{Laptop: laptop}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := laptopClient.CreateLaptop(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Printf("laptop elready exists")
		} else {
			log.Fatal("cannot create laptop: ", err)
		}
		return
	}

	log.Printf("created laptop with id: %v", res.Id)
}
