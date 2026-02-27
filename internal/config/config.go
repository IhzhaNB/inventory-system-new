package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// AppConfig holds application-specific settings
type AppConfig struct {
	Port string `mapstructure:"APP_PORT"`
	Env  string `mapstructure:"APP_ENV"`
}

// DBConfig holds database connection settings
type DBConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSL_MODE"`
}

// GetURL generates the standard PostgreSQL connection string
func (db *DBConfig) GetURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		db.User, db.Password, db.Host, db.Port, db.Name, db.SSLMode)
}

// JWTConfig holds authentication settings
type JWTConfig struct {
	Secret string `mapstructure:"JWT_SECRET"`
	Expire string `mapstructure:"JWT_EXPIRE"`
}

// Config is the master struct that groups all configurations
type Config struct {
	App AppConfig `mapstructure:",squash"`
	DB  DBConfig  `mapstructure:",squash"`
	JWT JWTConfig `mapstructure:",squash"`
}

// LoadConfig reads the configuration from the provided path.
func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(path + "/.env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
