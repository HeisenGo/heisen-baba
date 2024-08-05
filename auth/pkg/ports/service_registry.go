package ports

type IServiceRegistry interface {
	RegisterServiceRegisterService(serviceName, serviceHostAddress, servicePrefixPath, serviceHTTPHealthPath string, serviceGRPCPort, serviceHTTPPort int) error
}
