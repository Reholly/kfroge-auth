package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server ServerConfig   `yaml:"server"`
	Db     DatabaseConfig `yaml:"database"`
	Smtp   SmtpConfig     `yaml:"smtp"`
	Auth   AuthConfig     `yaml:"auth"`
	Seed   SeedConfig     `yaml:"seed"`
	Jwt    JwtConfig      `yaml:"jwt"`
	Cache  CacheConfig    `yaml:"cache"`
}

type CacheConfig struct {
	CodeExpirationTimeInMinutes int `yaml:"code_expiration_time_in_minutes"`
	CleanupIntervalInMinutes    int `yaml:"cleanup_interval_in_minutes"`
}

type ServerConfig struct {
	Rest string `yaml:"rest"`
	Grpc string `yaml:"grpc"`
}

type DatabaseConfig struct {
	ConnectionString string `yaml:"connection_string"`
	MigrationDir     string `yaml:"migration_dir"`
	DriverName       string `yaml:"driver_name"`
}

type SmtpConfig struct {
	SmtpFrom     string `yaml:"smtp_from"`
	SmtpPassword string `yaml:"smtp_password"`
	SmtpHost     string `yaml:"smtp_host"`
	SmtpPort     string `yaml:"smtp_port"`
}

type AuthConfig struct {
	CodeSalt                 string `yaml:"code_salt"`
	PasswordSalt             string `yaml:"password_salt"`
	EmailConfirmationUrlBase string `yaml:"email_confirmation_url_base"`
}

type JwtConfig struct {
	JwtSecret                      string `yaml:"jwt_secret"`
	AccessTokenTimeToLiveInSeconds int    `yaml:"access_token_ttl_seconds"`
	RefreshTokenTimeToLiveInHours  int    `yaml:"refresh_token_ttl_hours"`
}

type SeedConfig struct {
	AdminUsername string `yaml:"default_admin_username"`
	AdminPassword string `yaml:"default_admin_password"`
	AdminEmail    string `yaml:"default_admin_email"`
}

func LoadConfig(configPath string) (Config, error) {
	config := Config{}
	file, err := os.ReadFile(configPath)

	if err != nil {
		return Config{}, err
	}

	err = yaml.Unmarshal(file, &config)

	if err != nil {
		return Config{}, err
	}

	return config, nil
}
