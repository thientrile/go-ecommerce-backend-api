package setting

type Config struct {
	Server ServerSetting `mapstruct:"server"`
	Mysql  MSQLSetting   `mapstruct:"mysql"`
	Logger LoggerSetting `mapstruct:"logger"`
	Redis  RedisSetting  `mapstruct:"redis"`
	Kafka  KafkaSetting  `mapstruct:"kafka"`
	JWT    JWTSetting    `mapstruct:"jwt"`
}

type ServerSetting struct {
	Port int    `mapstruct:"port"`
	Mode string `mapstruct:"mode"`
}

type RedisSetting struct {
	Host      string `mapstruct:"host"`
	Port      int    `mapstruct:"port"`
	Password  string `mapstruct:"password"`
	Database  int    `mapstruct:"database"`
	Pool_size int    `mapstruct:"pool_size"`
}

type MSQLSetting struct {
	Host            string `mapstruct:"host"`
	Port            int    `mapstruct:"port"`
	Username        string `mapstruct:"username"`
	Password        string `mapstruct:"password"`
	Dbname          string `mapstruct:"dbname"`
	MaxIdleConns    int    `mapstruct:"maxIdleConns"`
	MaxOpenConns    int    `mapstruct:"maxOpenConns"`
	ConnMaxLifetime int    `mapstruct:"connMaxLifetime"`
	ConnMaxIdleTime int    `mapstruct:"connMaxIdleTime"`
}

type LoggerSetting struct {
	Log_Level     string `mapstruct:"log_level"`
	File_log_name string `mapstruct:"file_log_name"`
	File_log_path string `mapstruct:"file_log_path"`
	Max_size      int    `mapstruct:"max_size"`
	Max_backups   int    `mapstruct:"max_backups"`
	Max_age       int    `mapstruct:"max_age"`
	Compress      bool   `mapstruct:"compress"`
}

type KafkaSetting struct {
	Host  string `mapstruct:"host"`
	Port  int    `mapstruct:"port"`
	Topic struct {
		Product string `mapstruct:"product"`
		Order   string `mapstruct:"order"`
		Auth    string `mapstruct:"auth"`
	} `mapstruct:"topic"`
}

// JWTSetting holds the configuration for JWT authentication.
type JWTSetting struct {
	TokenHourLifespan int    `mapstruct:"TOKEN_HOUR_LIFETIME"`
	JwtExpiration     string `mapstruct:"JWT_EXPIRATION"`
	ApiSecret         string `mapstruct:"API_SECRET"`
	Issuer            string `mapstruct:"ISSUER"`
}
