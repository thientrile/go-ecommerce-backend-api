package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	Databases []struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
	} `mapstructure:"databases"`
}

func main() {
	viper := viper.New()
	viper.AddConfigPath("./configs/") // path to look for the config file in
	viper.SetConfigName("local")      // name of config file (without extension)
	viper.SetConfigType("yaml")       // REQUIRED if the config file does not have the extension in the name
	// read file config
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read configuation %w \n ", err))
	}

	// read server config
	fmt.Println("Server Port::", viper.GetString("server.port"))
	fmt.Println("Server Host::", viper.GetString("security.jwt.key"))

	// get config to struct

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Unable to decode into struct, %v\n", err)
	}
	fmt.Println("Config Server Port::", config.Server.Port)

	for _, db := range config.Databases {
		fmt.Printf("Database User: %s, Password: %s, Host: %s\n", db.User, db.Password, db.Host)
	}

}
