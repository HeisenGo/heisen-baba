package ports

type IServiceRegistry interface {
	RegisterService(serviceName, serviceHostAddress, servicePrefixPath, serviceHTTPHealthPath string, serviceGRPCPort, serviceHTTPPort int) error
}
