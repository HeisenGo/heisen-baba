package ports

type IServiceRegistry interface {
	RegisterService(serviceName, serviceHostAddress, servicePrefixPath, serviceHTTPHealthPath string, serviceGRPCPort, serviceHTTPPort int) error
	DiscoverService(serviceName string) (port int, ip string, err error)
}
