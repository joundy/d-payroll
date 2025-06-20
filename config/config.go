package config

import (
	"github.com/spf13/viper"
)

type PostgresConfig struct {
	Host     string
	Port     int32
	User     string
	Password string
	Db       string
}

type AdminUserConfig struct {
	Username string
	Password string
}

type HttpConfig struct {
	Port int32
	Host string
}

type AuthConfig struct {
	JwtSecret string
}

type OvertimeConfig struct {
	MaxDurationPerDayMilis int
}

type PayrollConfig struct {
	DayPerMonthProrate    int
	MaxWorkingMilisPerDay int
}

type Config struct {
	Postgres  *PostgresConfig
	AdminUser *AdminUserConfig
	Http      *HttpConfig
	Auth      *AuthConfig
	Overtime  *OvertimeConfig
	Payroll   *PayrollConfig
}

// TODO: config error handling and logging
func NewConfig() *Config {
	v := viper.New()

	// Set defaults and config file first
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")

	v.AutomaticEnv()
	v.ReadInConfig()

	return &Config{
		Postgres:  initPostgresConfig(v),
		AdminUser: initAdminUser(v),
		Http:      initHttpConfig(v),
		Auth:      initAuthConfig(v),
		Overtime: &OvertimeConfig{
			MaxDurationPerDayMilis: 1000 * 60 * 60 * 3,
		},
		Payroll: &PayrollConfig{
			DayPerMonthProrate:    22, // preference, could be 20, 30, etc..
			MaxWorkingMilisPerDay: 8 * 60 * 60 * 1000,
		},
	}
}

func initPostgresConfig(v *viper.Viper) *PostgresConfig {
	v.SetDefault("POSTGRES_HOST", "localhost")
	v.SetDefault("POSTGRES_PORT", "5432")
	v.SetDefault("POSTGRES_USER", "postgres")
	v.SetDefault("POSTGRES_PASSWORD", "password")
	v.SetDefault("POSTGRES_DB", "d-payroll")

	return &PostgresConfig{
		Host:     v.GetString("POSTGRES_HOST"),
		Port:     v.GetInt32("POSTGRES_PORT"),
		User:     v.GetString("POSTGRES_USER"),
		Password: v.GetString("POSTGRES_PASSWORD"),
		Db:       v.GetString("POSTGRES_DB"),
	}
}

func initAdminUser(v *viper.Viper) *AdminUserConfig {
	v.SetDefault("ADMIN_USERNAME", "admin")
	v.SetDefault("ADMIN_PASSWORD", "VERY_STRONG_PASSWORD")

	return &AdminUserConfig{
		Username: v.GetString("ADMIN_USERNAME"),
		Password: v.GetString("ADMIN_PASSWORD"),
	}
}

func initHttpConfig(v *viper.Viper) *HttpConfig {
	v.SetDefault("HTTP_PORT", "3000")
	v.SetDefault("HTTP_HOST", "0.0.0.0")

	return &HttpConfig{
		Port: v.GetInt32("HTTP_PORT"),
		Host: v.GetString("HTTP_HOST"),
	}
}

func initAuthConfig(v *viper.Viper) *AuthConfig {
	v.SetDefault("AUTH_JWT_SECRET", "secret")

	return &AuthConfig{
		JwtSecret: v.GetString("AUTH_JWT_SECRET"),
	}
}
