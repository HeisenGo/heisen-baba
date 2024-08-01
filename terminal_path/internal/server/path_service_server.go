// internal/server/path_service_server.go
package server

import (
	"context"
	"errors"

	"terminalpathservice/internal/terminal"
	"terminalpathservice/service"
	pb "terminalpathservice/terminal_path/internal/protobufs"
)

type PathServiceServer struct {
    pb.UnimplementedPathServiceServer
    pathService *service.PathService
}

func NewPathServiceServer(pathService *service.PathService) *PathServiceServer {
    return &PathServiceServer{pathService: pathService}
}

func (s *PathServiceServer) GetFullPathByID(ctx context.Context, req *pb.GetFullPathByIDRequest) (*pb.GetFullPathByIDResponse, error) {
    pathID := req.GetPathID()

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
        DistanceKm:       path.DistanceKM,
        Code:             path.Code,
        Name:             path.Name,
        Type:             path.Type,
        FromCountry:      path.FromTerminal.Country,
        FromCity:         path.FromTerminal.City,
        FromTerminalName: path.FromTerminal.Name,
        ToCountry:        path.ToTerminal.Country,
        ToCity:           path.ToTerminal.City,
        ToTerminalName:   path.ToTerminal.Name,
    }

    return &pb.GetFullPathByIDResponse{
        Message: "Path fetched successfully",
        Data:    data,
    }, nil
}
