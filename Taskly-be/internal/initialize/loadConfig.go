package initialize

import (
	"fmt"
	"strings"

	"Taskly.com/m/global"

	"github.com/spf13/viper"
)

func LoadConfigProd() {
	v := viper.New()

	// map db.host -> DB_HOST, cloud_name -> CLOUD_NAME
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Bind explicit: local key -> ENV NAME
	_ = v.BindEnv("cloud_name", "CLOUD_NAME")
	_ = v.BindEnv("api_key", "API_KEY")
	_ = v.BindEnv("api_secret", "API_SECRET")
	_ = v.BindEnv("database_url_internal", "DATABASE_URL_INTERNAL")
	_ = v.BindEnv("database_url_external", "DATABASE_URL_EXTERNAL")
	_ = v.BindEnv("vnp_tmncode", "VNP_TMNCODE")
	_ = v.BindEnv("vnp_hashsecret", "VNP_HASHSECRET")
	_ = v.BindEnv("vnp_url", "VNP_URL")
	_ = v.BindEnv("vnp_url_callback", "VNP_URL_CALLBACK")
	_ = v.BindEnv("vnp_ipn_url", "VNP_IPN_URL")
	_ = v.BindEnv("redis_url", "REDIS_URL")

	// Log / JWT...
	_ = v.BindEnv("log_level", "LOG_LEVEL")
	_ = v.BindEnv("log_file_name", "LOG_FILE_NAME")
	_ = v.BindEnv("token_hour_lifespan", "TOKEN_HOUR_LIFESPAN")
	_ = v.BindEnv("api_secret_jwt", "API_SECRET_JWT") // note: you used API_SECRET twice in struct

	// Unmarshal vào struct. IMPORTANT: mapstructure tags trong struct phải
	// trùng với "local key" bạn bind ở trên (không bắt buộc in hoa).
	if err := v.Unmarshal(&global.ENVSetting); err != nil {
		panic(fmt.Errorf("unable to decode env settings: %w", err))
	}

	// debug
	fmt.Println("CLOUD_NAME =", v.GetString("cloud_name"))
}

func LoadConfigDev() {
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
	if err := envViper.Unmarshal(&global.ENVSetting); err != nil {
		fmt.Printf("unable to decode configuration %v", err)
	}
}
