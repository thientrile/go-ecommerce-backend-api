package initialize

import (
	"fmt"

	"github.com/spf13/viper"
	"go-ecommerce-backend-api.com/global"
)


func LoadConfig() {
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

	if err := viper.Unmarshal(&global.Config); err != nil {
		fmt.Printf("Unable to decode into struct, %v\n", err)
	}
}
