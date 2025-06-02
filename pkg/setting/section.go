package setting

type Config struct {
	Server ServerSetting `mapstruct:"server"`
	Mysql  MSQLSetting   `mapstruct:"mysql"`
	Logger LoggerSetting `mapstruct:"logger"`
	Redis  RedisSetting  `mapstruct:"redis"`
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
