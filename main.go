package main

import (
	"errors"
	"go-htmx-test/db"
	"go-htmx-test/web/home"
	"log"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

func main() {
	// Get environment variables
	env, err := getEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to database
	db.Connect(env.DB_Host, env.DB_User, env.DB_Pass, env.DB_Name, env.DB_Port, env.Timezone)

	// Echo instance
	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.Logger())

	// Home routes
	homeHandler := home.HomeHandler{}
	e.Any("/", homeHandler.Any)
	e.Static("/assets", "assets")

	e.Logger.Fatal(e.Start(":" + env.Port))
}

type envConfig struct {
	Port     string `mapstructure:"PORT"`
	DB_Host  string `mapstructure:"DB_HOST"`
	DB_User  string `mapstructure:"DB_USER"`
	DB_Pass  string `mapstructure:"DB_PASS"`
	DB_Name  string `mapstructure:"DB_NAME"`
	DB_Port  int    `mapstructure:"DB_PORT"`
	Timezone string `mapstructure:"TZ"`
}

// getEnvFile gets the environment variables from a file
func getEnvFile() (envConfig, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		return envConfig{}, errors.New("error reading env file")
	}

	config := envConfig{}

	if err := viper.Unmarshal(&config); err != nil {
		return envConfig{}, errors.New("error unmarshaling env")
	}

	return config, nil
}

// getEnvOS gets the environment variables from the OS
func getEnvOS() (envConfig, error) {
	config := envConfig{}

	config.Port = os.Getenv("PORT")
	config.DB_Host = os.Getenv("DB_HOST")
	config.DB_User = os.Getenv("DB_USER")
	config.DB_Pass = os.Getenv("DB_PASS")
	config.DB_Name = os.Getenv("DB_NAME")
	config.Timezone = os.Getenv("TZ")

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return envConfig{}, errors.New("env error; DB_PORT must be an integer")
	}
	config.DB_Port = dbPort

	return config, nil
}

// getEnv gets the environment variables from a file or the OS
func getEnv() (envConfig, error) {
	// Try to get env from file first
	env, err := getEnvFile()
	if err != nil {
		log.Println(err)
		log.Println("Trying OS env...")

		// If file doesn't exist, try to get env from OS
		env, err = getEnvOS()
		if err != nil {
			return envConfig{}, err
		}
	}

	// Validate env
	if err := validateEnv(env); err != nil {
		return envConfig{}, err
	}

	return env, nil
}

// validateEnv validates the environment variables
func validateEnv(env envConfig) error {
	if env.Port == "" {
		return errors.New("PORT is required")
	}

	if env.DB_Host == "" {
		return errors.New("DB_HOST is required")
	}

	if env.DB_User == "" {
		return errors.New("DB_USER is required")
	}

	if env.DB_Pass == "" {
		return errors.New("DB_PASS is required")
	}

	if env.DB_Name == "" {
		return errors.New("DB_NAME is required")
	}

	if env.DB_Port == 0 {
		return errors.New("DB_PORT is required")
	}

	return nil
}
