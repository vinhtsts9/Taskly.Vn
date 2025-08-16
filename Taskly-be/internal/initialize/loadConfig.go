package initialize

import (
	"fmt"

	"Taskly.com/m/global"

	"github.com/spf13/viper"
)

func LoadConfig() {
	// Load local.yaml
	yamlViper := viper.New()
	yamlViper.AddConfigPath("./configs")
	yamlViper.SetConfigName("local")
	yamlViper.SetConfigType("yaml")

	err := yamlViper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read config %w", err))
	}

	// Load app.env
	envViper := viper.New()
	envViper.AddConfigPath(".")
	envViper.SetConfigName("app")
	envViper.SetConfigType("env")

	if err := envViper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config from app.env:", err)
	}

	// Gộp biến môi trường từ .env vào
	envViper.AutomaticEnv()

	fmt.Println("server port", yamlViper.GetInt("server.port"))
	fmt.Println("security jwt key", envViper.GetString("CLOUD_NAME"))

	// Kết hợp cả YAML và ENV vào global.Config
	if err := yamlViper.Unmarshal(&global.Config); err != nil {
		fmt.Printf("unable to decode configuration %v", err)
	}
	if err := envViper.Unmarshal(&global.CloudinarySetting); err != nil {
		fmt.Printf("unable to decode configuration %v", err)
	}
}
