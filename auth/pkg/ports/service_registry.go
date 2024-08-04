package ports

type IServiceRegistry interface {

	RegisterServiceRegisterService(serviceName, serviceHostAddress, servicePrefixPath, serviceHTTPHealthPath string, serviceHTTPPort int) error
}
