package service

import (
	"context"
	"net"
	"testing"

	"github.com/srjcshv/grpc/pb"
	"github.com/srjcshv/grpc/sample"
	"github.com/srjcshv/grpc/serializer"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startTestLaptopServer(t *testing.T, laptopStore LaptopStore) string {
	laptopServer := NewLaptopServer(laptopStore)
	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	go grpcServer.Serve(listener)

	return listener.Addr().String()
}

func newTestLaptopClient(t *testing.T, serverAddress string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(conn)
}

func TestService_CreateLaptopClient(t *testing.T) {
	t.Parallel()

	laptopStore := NewInMemoryLaptopStore()
	serverAddress := startTestLaptopServer(t, laptopStore)
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

func requireSameLaptop(t *testing.T, laptop1 *pb.Laptop, laptop2 *pb.Laptop) {
	json1, err := serializer.ProtobufToJson(laptop1)
	require.NoError(t, err)

	json2, err := serializer.ProtobufToJson(laptop2)
	require.NoError(t, err)

	require.Equal(t, json1, json2)

}
