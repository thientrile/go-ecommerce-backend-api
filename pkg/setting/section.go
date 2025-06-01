package setting

type Config struct {
	Msql MSQLSetting `mapstruct:"mysql"`
}

type MSQLSetting struct {
	Host string `mapstruct:"host"`
	Port int    `mapstruct:"port"`
	Username string `mapstruct:"username"`
	Password string `mapstruct:"password"`
	Dbname string `mapstruct:"dbname"`
	MaxIdleConns int `mapstruct:"maxIdleConns"`
	MaxOpenConns int `mapstruct:"maxOpenConns"`
	ConnMaxLifetime int `mapstruct:"connMaxLifetime"`
}