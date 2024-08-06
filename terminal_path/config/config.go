package config

type Config struct {
	Server Server `mapstructure:"server"`
	DB     DB     `mapstructure:"db"`
	Redis  Redis  `mapstructure:"redis"`
}

type Server struct {
	HTTPPort              int             `mapstructure:"http_port"`
	GRPCPort              int             `mapstructure:"grpc_port"`
	Host                  string          `mapstructure:"host"`
	ServiceRegistry       ServiceRegistry `mapstructure:"service_registry"`
	ServiceHostAddress    string          `mapstructure:"service_host_address"`
	ServiceHostName       string          `mapstructure:"service_host_name"`
	ServiceHTTPHealthPath string          `mapstructure:"service_http_health_path"`
	ServiceHTTPPrefixPath string          `mapstructure:"service_http_prefix_path"`
}

type ServiceRegistry struct {
	Address                string `mapstructure:"address"`
	AuthServiceName        string `mapstructure:"auth_service_name"`
	TripCompanyServiceName string `mapstructure:"trip_company_service_name"`
	ServiceName            string `mapstructure:"service_name"`
}

type DB struct {
	User          string `mapstructure:"user"`
	Pass          string `mapstructure:"pass"`
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	DBName        string `mapstructure:"db_name"`
	AppCommission uint   `mapstructure:"app_commission"`
}

type Redis struct {
	Pass string `mapstructure:"pass"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
