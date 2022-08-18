package serializer

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/srjcshv/grpc/pb"
	"github.com/srjcshv/grpc/sample"
	"github.com/stretchr/testify/require"
)

func TestSerializer_WriteProtobufToBinaryFile(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	jsonFile := "../tmp/laptop.json"

	laptop1 := sample.NewLaptop()
	err := WriteProtobufToBinaryFile(laptop1, binaryFile)
	require.NoError(t, err)

	err = WriteProtobufToJsonFile(laptop1, jsonFile)
	require.NoError(t, err)

	laptop2 := &pb.Laptop{}
	err = ReadProtobufFromBinaryFile(binaryFile, laptop2)
	require.NoError(t, err)

	require.True(t, proto.Equal(laptop1, laptop2))
}
