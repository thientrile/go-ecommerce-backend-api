package setting

type Config struct {
	Server  ServerSetting  `mapstructure:"server"`
	Mysql   MSQLSetting    `mapstructure:"mysql"`
	Logger  LoggerSetting  `mapstructure:"logger"`
	Redis   RedisSetting   `mapstructure:"redis"`
	Kafka   KafkaSetting   `mapstructure:"kafka"`
	JWT     JWTSetting     `mapstructure:"jwt"`
	Limiter LimiterSetting `mapstructure:"limiter"`
}

type ServerSetting struct {
	Port    int    `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
	Domain  string `mapstructure:"domain"`
	Version string `mapstructure:"version"`
}

type RedisSetting struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Password  string `mapstructure:"password"`
	Database  int    `mapstructure:"database"`
	Pool_size int    `mapstructure:"pool_size"`
}

type MSQLSetting struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Dbname          string `mapstructure:"dbname"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns"`
	ConnMaxLifetime int    `mapstructure:"connMaxLifetime"`
	ConnMaxIdleTime int    `mapstructure:"connMaxIdleTime"`
}

type LoggerSetting struct {
	Log_Level     string `mapstructure:"log_level"`
	File_log_name string `mapstructure:"file_log_name"`
	File_log_path string `mapstructure:"file_log_path"`
	Max_size      int    `mapstructure:"max_size"`
	Max_backups   int    `mapstructure:"max_backups"`
	Max_age       int    `mapstructure:"max_age"`
	Compress      bool   `mapstructure:"compress"`
}

type KafkaSetting struct {
	Host  string            `mapstructure:"host"`  // ✅
	Port  int               `mapstructure:"port"`  // ✅
	Topic map[string]string `mapstructure:"topic"` // ✅
}

// JWTSetting holds the configuration for JWT authentication.
type JWTSetting struct {
	TokenHourLifespan int    `mapstructure:"TOKEN_HOUR_LIFETIME"`
	JwtExpiration     string `mapstructure:"JWT_EXPIRATION"`
	ApiSecret         string `mapstructure:"API_SECRET"`
	Issuer            string `mapstructure:"ISSUER"`
}

type LimiterSetting struct {
	Store         int                    `mapstructure:"store" json:"store"`
	DefaultConfig map[string]interface{} `mapstructure:"default_config" json:"default_config"`
	Rules         map[string]RuleConfig  `mapstructure:"rules" json:"rules"`
	URLPath       struct {
		Public  []string `mapstructure:"public"`
		Private []string `mapstructure:"private"`
	} `mapstructure:"url_path"`
}

type RuleConfig struct {
	Rate            string   `mapstructure:"rate" json:"rate"`
	Description     string   `mapstructure:"description" json:"description"`
	Enabled         bool     `mapstructure:"enabled" json:"enabled"`
	BurstMultiplier int      `mapstructure:"burst_multiplier" json:"burst_multiplier"`
	StrictMode      bool     `mapstructure:"strict_mode" json:"strict_mode"`
	IPWhitelist     []string `mapstructure:"ip_whitelist" json:"ip_whitelist"`
	MaxFileSize     string   `mapstructure:"max_file_size" json:"max_file_size"`
}
