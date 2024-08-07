package ports

type IServiceRegistry interface {
	RegisterService(serviceName, serviceHostName, servicePrefixPath, serviceHTTPHealthPath string, serviceHTTPPort int) error //
	DiscoverService(serviceName string) (port int, ip string, err error)
}
