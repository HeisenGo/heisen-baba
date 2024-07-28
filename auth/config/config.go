package config

type Config struct {
	Server Server `mapstructure:"server"`
	DB     DB     `mapstructure:"db"`
}

type Server struct {
	GRPCPort               int             `mapstructure:"grpc_port"`
	Host                   string          `mapstructure:"host"`
	TokenExpMinutes        uint            `mapstructure:"token_exp_minutes"`
	RefreshTokenExpMinutes uint            `mapstructure:"refresh_token_exp_minutes"`
	TokenSecret            string          `mapstructure:"token_secret"`
	ServiceRegistry        ServiceRegistry `mapstructure:"service_registry"`
	ServiceHostName        string          `mapstructure:"service_host_name"`
	ServiceGRPCHealthPath  string          `mapstructure:"service_grpc_health_path"`
	ServiceGRPCPrefixPath  string          `mapstructure:"service_grpc_prefix_path"`
}
type ServiceRegistry struct {
	Address string `mapstructure:"address"`
}

type DB struct {
	User   string `mapstructure:"user"`
	Pass   string `mapstructure:"pass"`
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	DBName string `mapstructure:"db_name"`
}
