package config

type Config struct {
	App      App      `mapstructure:"app" validate:"required"`
	Database Database `mapstructure:"database" validate:"required"`
	Redis    Redis    `mapstructure:"redis" validate:"required"`
	Logar    Logar    `mapstructure:"logar" validate:"required"`
	Provider Provider `mapstructure:"provider" validate:"required"`
}

type App struct {
	URL              string `mapstructure:"url" validate:"required"`
	Port             int    `mapstructure:"port" validate:"required"`
	JwtSecret        string `mapstructure:"jwt_secret" validate:"required"`
	AdminJwtSecret   string `mapstructure:"admin_jwt_secret" validate:"required"`
	AuthExpireHours  int    `mapstructure:"auth_expire_hours" validate:"required"`
	OtpExpireSeconds int    `mapstructure:"otp_expire_seconds" validate:"required"`
}

type Database struct {
	Name    string `mapstructure:"name" validate:"required"`
	Host    string `mapstructure:"host" validate:"required"`
	Pass    string `mapstructure:"pass" validate:"required"`
	User    string `mapstructure:"user" validate:"required"`
	Port    int16  `mapstructure:"port" validate:"required"`
	Debug   bool   `mapstructure:"debug"`
	Migrate bool   `mapstructure:"migrate"`
}

type Redis struct {
	Host string `mapstructure:"host" validate:"required"`
	Port int    `mapstructure:"port" validate:"required"`
	Pass string `mapstructure:"pass" validate:"required"`
}

type Logar struct {
	AdminUsername string `mapstructure:"admin_username" validate:"omitempty"`
	AdminPassword string `mapstructure:"admin_password" validate:"omitempty"`
}
type Provider struct {
	Provider1 string `mapstructure:"provider1" validate:"omitempty"`
	Provider2 string `mapstructure:"provider2" validate:"omitempty"`
}

var config Config

func Get() Config {
	return config
}

func Set(c Config) {
	config = c
}
