package memoryServer

import (
	"context"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/sirupsen/logrus"
)

type grpcMemoryServiceServer struct {
	pb.UnimplementedMemoryServiceServer
}

func NewGrpcMemoryServer() *grpcMemoryServiceServer {
	return &grpcMemoryServiceServer{}
}

// List
// list memory files
func (s *grpcMemoryServiceServer) List(ctx context.Context, in *pb.ListMemoryRequest) (*pb.ListMemoryResponse, error) {

	logrus.Infof("List called on Memory Server %v", in)
	memories, err := memory.ListMemories()
	if err != nil {
		return nil, err
	}

	return &pb.ListMemoryResponse{
		Memories: memories,
	}, nil

}

func (s *grpcMemoryServiceServer) Get(ctx context.Context, in *pb.GetMemoryRequest) (*pb.GetMemoryResponse, error) {

	logrus.Infof("Get called on Memory Server %v", in)

	memory, err := memory.LoadMemory(in.Name)
	if err != nil {
		return nil, err
	}

	return &pb.GetMemoryResponse{
		Meta: memory.Meta,
	}, nil

}

func (s *grpcMemoryServiceServer) Delete(ctx context.Context, in *pb.DeleteMemoryRequest) (*pb.DeleteMemoryResponse, error) {

	logrus.Infof("Get called on Memory Server %v", in)

	err := memory.DeleteMemory(in.Name)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteMemoryResponse{
		Success: true,
	}, nil

}

func (s *grpcMemoryServiceServer) Add(ctx context.Context, in *pb.AddMemoryRequest) (*pb.AddMemoryResponse, error) {

	logrus.Infof("Add called on Memory Server %v", in)

	memory := memory.NewMemory(in.Meta.Name)
	memory.Meta.Uuid = in.Meta.Uuid
	memory.Meta.Description = in.Meta.Description
	memory.Meta.Model = in.Meta.Model
	memory.Meta.Seed = in.Meta.Seed
	memory.Meta.Temperature = in.Meta.Temperature
	memory.Meta.SystemPrompt = in.Meta.SystemPrompt

	return &pb.AddMemoryResponse{
		Meta: memory.Meta,
	}, nil

}

func (s *grpcMemoryServiceServer) Update(ctx context.Context, in *pb.UpdateMemoryRequest) (*pb.UpdateMemoryResponse, error) {

	logrus.Infof("Update called on Memory Server  %v", in)

	memory := memory.NewMemory(in.Meta.Name)
	memory.Meta.Uuid = in.Meta.Uuid
	memory.Meta.Description = in.Meta.Description
	memory.Meta.Model = in.Meta.Model
	memory.Meta.Seed = in.Meta.Seed
	memory.Meta.Temperature = in.Meta.Temperature
	memory.Meta.SystemPrompt = in.Meta.SystemPrompt

	return &pb.UpdateMemoryResponse{
		Meta: memory.Meta,
	}, nil

}
