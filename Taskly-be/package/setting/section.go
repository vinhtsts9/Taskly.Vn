package setting

type Config struct {
	Server        ServerSetting     `mapstructure:"server"`
	Redis         RedisSetting      `mapstructure:"redis"`
	PostgreSQL    PostgreSQLSetting `mapstructure:"postgresql"`
	KafkaBroker   Kafka             `mapstructure:"kafka"`
	ElasticSearch ElasticSearch     `mapstructure:"elasticsearch"`
}
type ENV struct {
	CloudName             string `mapstructure:"CLOUD_NAME"`
	ApiKey                string `mapstructure:"API_KEY"`
	ApiSecret             string `mapstructure:"API_SECRET"`
	Database_url_internal string `mapstructure:"DATABASE_URL_INTERNAL"`
	Database_url_external string `mapstructure:"DATABASE_URL_EXTERNAL"`
	Vnp_TmnCode           string `mapstructure:"VNP_TMNCODE"`
	Vnp_HashSecret        string `mapstructure:"VNP_HASHSECRET"`
	Vnp_Url               string `mapstructure:"VNP_URL"`
	Vnp_UrlCallBack       string `mapstructure:"VNP_URL_CALLBACK"`
	Vnp_IpnUrl            string `mapstructure:"VNP_IPN_URL"`
	Redis_Url             string `mapstructure:"REDIS_URL"`
	Brevo_ApiKey          string `mapstructure:"BREVO_APIKEY"`

	// Log
	LogLevel      string `mapstructure:"LOG_LEVEL"`
	LogFileName   string `mapstructure:"LOG_FILE_NAME"`
	LogMaxSize    int    `mapstructure:"LOG_MAX_SIZE"`
	LogMaxBackups int    `mapstructure:"LOG_MAX_BACKUPS"`
	LogMaxAge     int    `mapstructure:"LOG_MAX_AGE"`
	LogCompress   bool   `mapstructure:"LOG_COMPRESS"`

	// JWT
	TokenHourLifespan uint   `mapstructure:"TOKEN_HOUR_LIFESPAN"`
	JwtExpiration     string `mapstructure:"JWT_EXPIRATION"`
	RefreshExpiration string `mapstructure:"REFRESH_EXPIRATION"`
	ApiSecretJwt      string `mapstructure:"API_SECRET_JWT"`

	// Mail
	SMTP_HOST     string `mapstructure:"SMTP_HOST"`
	SMTP_PORT     string `mapstructure:"SMTP_PORT"`
	SMTP_USERNAME string `mapstructure:"SMTP_USERNAME"`
	SMTP_PASSWORD string `mapstructure:"SMTP_PASSWORD"`
	SMTP_FROM     string `mapstructure:"SMTP_FROM"`

	Fe_api string `mapstructure:"FE_API"`
	Be_api string `mapstructure:"BE_API"`
}
type ServerSetting struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}
type RedisSetting struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}
type PostgreSQLSetting struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
}
type Kafka struct {
	Brokers string `mapstructure:"brokers"`
}
type ElasticSearch struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
