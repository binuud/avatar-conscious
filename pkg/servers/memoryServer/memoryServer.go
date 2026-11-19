package memoryServer

import (
	"context"
	"fmt"

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

func (s *grpcMemoryServiceServer) GetMemoryMeta(ctx context.Context, in *pb.GetMemoryMetaRequest) (*pb.GetMemoryMetaResponse, error) {

	logrus.Infof("Get called on Memory Server %v", in)

	if in.Name == "" {
		// return empty memory
		m := memory.NewEmtpyMemory()
		return &pb.GetMemoryMetaResponse{
			Meta: m.Memory.Meta,
		}, nil
	}

	memory, err := memory.LoadMemory(in.Name)
	if err != nil {
		return nil, err
	}

	return &pb.GetMemoryMetaResponse{
		Meta: memory.Meta,
	}, nil

}

func (s *grpcMemoryServiceServer) Get(ctx context.Context, in *pb.GetMemoryRequest) (*pb.GetMemoryResponse, error) {

	logrus.Infof("Get called on Memory Server %v", in)

	if in.Name == "" {
		// return empty memory
		m := memory.NewEmtpyMemory()
		return &pb.GetMemoryResponse{
			Memory: &m.Memory,
		}, nil
	}

	m, err := memory.LoadMemory(in.Name)
	if err != nil {
		return nil, err
	}

	return &pb.GetMemoryResponse{
		Memory: &m.Memory,
	}, nil

}

func (s *grpcMemoryServiceServer) Delete(ctx context.Context, in *pb.DeleteMemoryRequest) (*pb.DeleteMemoryResponse, error) {

	logrus.Infof("Delete called on Memory Server %v", in)
	if in.Name == "" {
		return nil, fmt.Errorf("memory name is emtpy")
	}
	err := memory.DeleteMemory(in.Name)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteMemoryResponse{
		Success: true,
	}, nil

}

// ClearMemory
// Deletes the memory from storage and creates new
// Handle mutex ##TBD Binu
func (s *grpcMemoryServiceServer) ClearMemory(ctx context.Context, in *pb.ClearMemoryRequest) (*pb.ClearMemoryResponse, error) {

	logrus.Info("Recieved ClearMemory")
	m := memory.NewMemory(in.Name)
	m.ClearMemory()
	m.DumpMemory(nil)
	return &pb.ClearMemoryResponse{Success: true}, nil
}

func (s *grpcMemoryServiceServer) Add(ctx context.Context, in *pb.AddMemoryRequest) (*pb.AddMemoryResponse, error) {

	logrus.Infof("Add called on Memory Server %v", in)

	if in.Meta.Name == "" {
		return nil, fmt.Errorf("memory name is emtpy")
	}

	memory := memory.NewMemory(in.Meta.Name)
	memory.Meta.Uuid = in.Meta.Uuid
	memory.Meta.Description = in.Meta.Description
	memory.Meta.Model = in.Meta.Model
	memory.Meta.Seed = in.Meta.Seed
	memory.Meta.Temperature = in.Meta.Temperature
	memory.SetSystemPrompt(in.Meta.SystemPrompt)
	memory.Save()
	memory.DumpMemory(nil)
	return &pb.AddMemoryResponse{
		Meta: memory.Meta,
	}, nil

}

func (s *grpcMemoryServiceServer) Update(ctx context.Context, in *pb.UpdateMemoryRequest) (*pb.UpdateMemoryResponse, error) {

	logrus.Infof("Update called on Memory Server  %v", in)

	if in.Meta.Name == "" {
		return nil, fmt.Errorf("memory name is emtpy")
	}

	memory := memory.NewMemory(in.Meta.Name)
	memory.Meta.Uuid = in.Meta.Uuid
	memory.Meta.Description = in.Meta.Description
	memory.Meta.Model = in.Meta.Model
	memory.Meta.Seed = in.Meta.Seed
	memory.Meta.Temperature = in.Meta.Temperature
	memory.SetSystemPrompt(in.Meta.SystemPrompt)
	memory.Save()
	memory.DumpMemory(nil)
	return &pb.UpdateMemoryResponse{
		Meta: memory.Meta,
	}, nil

}
