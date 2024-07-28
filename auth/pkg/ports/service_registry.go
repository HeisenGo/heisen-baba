package ports

type IServiceRegistry interface {
	RegisterServiceRegisterService(serviceHostName, servicePrefixPath, serviceHTTPHealthPath string, serviceHTTPPort int) error
}
