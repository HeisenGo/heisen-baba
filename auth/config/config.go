package config

type Config struct {
	Server Server `mapstructure:"server"`
	DB     DB     `mapstructure:"db"`
	Redis  Redis  `mapstructure:"redis"`
}

type Server struct {
	HTTPPort               int             `mapstructure:"http_port"`
	Host                   string          `mapstructure:"host"`
	TokenExpMinutes        uint            `mapstructure:"token_exp_minutes"`
	RefreshTokenExpMinutes uint            `mapstructure:"refresh_token_exp_minutes"`
	TokenSecret            string          `mapstructure:"token_secret"`
	ServiceRegistry        ServiceRegistry `mapstructure:"service_registry"`
	ServiceHostName        string          `mapstructure:"service_host_name"`
	ServiceHTTPHealthPath  string          `mapstructure:"service_http_health_path"`
	ServiceHTTPPrefixPath  string          `mapstructure:"service_http_prefix_path"`
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

type Redis struct {
	Pass string `mapstructure:"pass"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
