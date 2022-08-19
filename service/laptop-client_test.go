package service

import (
	"context"
	"io"
	"net"
	"testing"

	"github.com/srjcshv/grpc/pb"
	"github.com/srjcshv/grpc/sample"
	"github.com/srjcshv/grpc/serializer"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startTestLaptopServer(t *testing.T, laptopStore LaptopStore) (*LaptopServer, string) {
	laptopServer := NewLaptopServer(laptopStore)
	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	go grpcServer.Serve(listener)

	return laptopServer, listener.Addr().String()
}

func newTestLaptopClient(t *testing.T, serverAddress string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(conn)
}

func TestService_CreateLaptopClient(t *testing.T) {
	t.Parallel()

	laptopStore := NewInMemoryLaptopStore()
	_, serverAddress := startTestLaptopServer(t, laptopStore)
	laptopClient := newTestLaptopClient(t, serverAddress)

	laptop1 := sample.NewLaptop()
	expectedID := laptop1.Id

	req := &pb.CreateLaptopRequest{
		Laptop: laptop1,
	}

	res, err := laptopClient.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedID, res.Id)

	laptop2, err := laptopStore.Find(res.Id)
	require.NoError(t, err)
	require.NotNil(t, laptop2)

	requireSameLaptop(t, laptop1, laptop2)
}

func TestService_ClientSearchLaptop(t *testing.T) {
	t.Parallel()
	filter := &pb.Filter{
		MaxPriceUsd: 2000,
		MinCpuCores: 4,
		MinCpuGhz:   2.2,
		MinRam: &pb.Memory{
			Value: 8,
			Unit:  pb.Memory_GIGABYTE,
		},
	}
	store := NewInMemoryLaptopStore()
	expectedIDs := make(map[string]bool)
	for i := 0; i < 6; i++ {
		laptop := sample.NewLaptop()

		switch i {
		case 0:
			laptop.PriceUsd = 2500
		case 1:
			laptop.Cpu.NumberCores = 2
		case 2:
			laptop.Cpu.MinGhz = 2.0
		case 3:
			laptop.Ram = &pb.Memory{Value: 4096, Unit: pb.Memory_MEGABYTE}
		case 4:
			laptop.PriceUsd = 1999
			laptop.Cpu.NumberCores = 4
			laptop.Cpu.MinGhz = 2.5
			laptop.Cpu.MaxGhz = laptop.Cpu.MinGhz + 2.0
			laptop.Ram = &pb.Memory{Value: 16, Unit: pb.Memory_GIGABYTE}
			expectedIDs[laptop.Id] = true
		case 5:
			laptop.PriceUsd = 2000
			laptop.Cpu.NumberCores = 6
			laptop.Cpu.MinGhz = 2.8
			laptop.Cpu.MaxGhz = laptop.Cpu.MinGhz + 2.0
			laptop.Ram = &pb.Memory{Value: 64, Unit: pb.Memory_GIGABYTE}
			expectedIDs[laptop.Id] = true
		}

		err := store.Save(laptop)
		require.NoError(t, err)

		_, serverAddress := startTestLaptopServer(t, store)
		laptopClient := newTestLaptopClient(t, serverAddress)

		req := &pb.SearchLaptopRequest{Filter: filter}
		stream, err := laptopClient.SearchLaptop(context.Background(), req)
		require.NoError(t, err)

		found := 0
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			require.NoError(t, err)
			require.Contains(t, expectedIDs, res.GetLaptop().GetId())
			found++
		}
		require.Equal(t, len(expectedIDs), found)
	}
}

func requireSameLaptop(t *testing.T, laptop1 *pb.Laptop, laptop2 *pb.Laptop) {
	json1, err := serializer.ProtobufToJson(laptop1)
	require.NoError(t, err)

	json2, err := serializer.ProtobufToJson(laptop2)
	require.NoError(t, err)

	require.Equal(t, json1, json2)

}
