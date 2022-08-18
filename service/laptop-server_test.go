package service

import (
	"context"
	"testing"

	"github.com/srjcshv/grpc/pb"
	"github.com/srjcshv/grpc/sample"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestService_CreateLaptop(t *testing.T) {
	t.Parallel()

	laptopNoId := sample.NewLaptop()
	laptopNoId.Id = ""

	laptopInvalidID := sample.NewLaptop()
	laptopInvalidID.Id = "invalid-uuid"

	laptopDuplicatedID := sample.NewLaptop()
	storeDuclicatedID := NewInMemoryLaptopStore()
	err := storeDuclicatedID.Save(laptopDuplicatedID)
	require.Nil(t, err)

	tests := []struct {
		name   string
		laptop *pb.Laptop
		store  LaptopStore
		code   codes.Code
	}{
		{
			name:   "ok",
			laptop: sample.NewLaptop(),
			store:  NewInMemoryLaptopStore(),
			code:   codes.OK,
		},
		{
			name:   "ok-noId",
			laptop: sample.NewLaptop(),
			store:  NewInMemoryLaptopStore(),
			code:   codes.OK,
		},
		{
			name:   "fail-invalidId",
			laptop: laptopInvalidID,
			store:  NewInMemoryLaptopStore(),
			code:   codes.InvalidArgument,
		},
		{
			name:   "fail-duplicatedId",
			laptop: laptopDuplicatedID,
			store:  storeDuclicatedID,
			code:   codes.AlreadyExists,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			req := &pb.CreateLaptopRequest{
				Laptop: test.laptop,
			}
			server := NewLaptopServer(test.store)
			res, err := server.CreateLaptop(context.Background(), req)
			if test.code == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotEmpty(t, res.Id)
				if len(test.laptop.Id) > 0 {
					require.Equal(t, test.laptop.Id, res.Id)
				}
			} else {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.code, st.Code())
			}
		})
	}
}
