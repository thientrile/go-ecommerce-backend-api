package setting

type Config struct {
	Mysql   MSQLSetting   `mapstruct:"mysql"`
	Logger LoggerSetting `mapstruct:"logger"`
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
	Max_size      int    `mapstruct:"max_size"`
	Max_backups   int    `mapstruct:"max_backups"`
	Max_age       int    `mapstruct:"max_age"`
	Compress      bool   `mapstruct:"compress"`
}
