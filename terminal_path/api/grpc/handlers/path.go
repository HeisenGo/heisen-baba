package handlers

import (
	"context"
	"errors"
	"terminalpathservice/internal/terminal"
	pb "terminalpathservice/protobufs"
	"terminalpathservice/service"

	"google.golang.org/grpc/health/grpc_health_v1"
)

type PathServiceServer struct {
	pb.UnimplementedPathServiceServer
	pathService *service.PathService
}

func NewPathServiceServer(pathService *service.PathService) *PathServiceServer {
	return &PathServiceServer{pathService: pathService}
}

func (s *PathServiceServer) GetFullPathByID(ctx context.Context, req *pb.GetFullPathByIDRequest) (*pb.GetFullPathByIDResponse, error) {
	pathID := req.GetPathId()

	if pathID < 0 {
		return nil, errors.New("invalid path ID")
	}

	path, err := s.pathService.GetFullPathByID(ctx, uint(pathID))
	if err != nil {
		if errors.Is(err, terminal.ErrRecordsNotFound) {
			return nil, errors.New("path not found")
		}
		return nil, errors.New("internal server error")
	}

	data := &pb.PathResponse{
		Id:               uint32(path.ID),
		DistanceKm:       path.DistanceKM, // in kilometers
		Code:             path.Code,
		Name:             path.Name,
		Type:             string(path.Type),
		FromCountry:      path.FromTerminal.Country,
		ToCountry:        path.ToTerminal.Country,
		FromCity:         path.FromTerminal.City,
		ToCity:           path.ToTerminal.City,
		FromTerminalName: path.FromTerminal.Name,
		ToTerminalName:   path.ToTerminal.Name,
	}

	return &pb.GetFullPathByIDResponse{
		Message: "Path fetched successfully",
		Data:    data,
	}, nil
}

type HealthServer struct {
	grpc_health_v1.HealthServer
}

// Check implements Health.Check
func (s *HealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}
