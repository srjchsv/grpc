package service

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/srjcshv/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopServer struct {
	LaptopStore LaptopStore
	pb.UnimplementedLaptopServiceServer
}

func NewLaptopServer(laptopStore LaptopStore) *LaptopServer {
	return &LaptopServer{LaptopStore: laptopStore}
}

func (server *LaptopServer) mustEmbedUnimplementedLaptopServiceServer() {}

func (server *LaptopServer) CreateLaptop(ctx context.Context, req *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Printf("received a create laptop request with id:%v", laptop.Id)

	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is no a valid uuid %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop  uuid %v", err)
		}
		laptop.Id = id.String()
	}

	switch ctx.Err() {
	case context.Canceled:
		log.Println("request is canceled")
		return nil, status.Error(codes.Canceled, "request is canceled")
	case context.DeadlineExceeded:
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	err := server.LaptopStore.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop to store: %v", err)
	}

	log.Printf("saved laptop with id: %v", laptop.Id)

	return &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}, nil
}

func (server *LaptopServer) SearchLaptop(req *pb.SearchLaptopRequest, stream pb.LaptopService_SearchLaptopServer) error {
	filter := req.GetFilter()
	log.Printf("recieved a search laptop request with filter: %v", filter)

	err := server.LaptopStore.Search(
		stream.Context(),
		filter,
		func(laptop *pb.Laptop) error {
			res := &pb.SearchLaptopResponse{Laptop: laptop}
			err := stream.Send(res)
			if err != nil {
				return err
			}
			log.Printf("sent laptop with id: %v", laptop.GetId())
			return nil
		},
	)
	if err != nil {
		return status.Errorf(codes.Internal, "unexpected errror: %v", err)
	}
	return nil
}
