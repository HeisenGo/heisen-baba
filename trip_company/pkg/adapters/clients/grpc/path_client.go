package grpc

import (
	"context"
	"fmt"
	"tripcompanyservice/internal/trip"
	"tripcompanyservice/pkg/adapters/clients/grpc/mappers.go"
	"tripcompanyservice/pkg/ports"
	"tripcompanyservice/protobufs"

	"google.golang.org/grpc"
)

type GRPCPathClient struct {
	ServiceRegistry ports.IServiceRegistry
	PathServiceName string
}

func NewGRPCPathClient(serviceRegistry ports.IServiceRegistry, pathServiceName string) *GRPCPathClient {
	return &GRPCPathClient{ServiceRegistry: serviceRegistry, PathServiceName: pathServiceName}
}

func (g *GRPCPathClient)GetFullPathByID(pathID uint32) (*trip.Path, error) {

	port, ip, err := g.ServiceRegistry.DiscoverService(g.PathServiceName)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", ip, port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	client := protobufs.NewPathServiceClient(conn)

	ctx := context.Background()

	request := &protobufs.GetFullPathByIDRequest{
		PathId: pathID,
	}

	// Call the GetUserByToken method
	response, err := client.GetFullPathByID(ctx, request)
	if err != nil {
		return nil, err
	}
	domainPath, err := mappers.GRPCPathToDomain(response)
	if err != nil {
		return nil, err
	}
	return domainPath, nil
}
