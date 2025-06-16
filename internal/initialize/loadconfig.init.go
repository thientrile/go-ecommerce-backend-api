package initialize

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go-ecommerce-backend-api.com/global"
)

func LoadConfig() {
	// // Load file .env trước
	// if err := godotenv.Load(); err != nil {
	// 	fmt.Printf("Warning: Could not load .env file: %v\n", err)
	// 	// Không return error vì không load được .env không phải lỗi nghiêm trọng
	// } else {
	// 	fmt.Println("Successfully loaded .env file")
	// }
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}
	fmt.Printf("Loading configuration for environment: %s\n", env)
	viper := viper.New()
	viper.AddConfigPath("./configs/") // path to look for the config file in
	viper.SetConfigName(env)          // name of config file (without extension)
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
