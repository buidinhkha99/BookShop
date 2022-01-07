package pkg

import "github.com/spf13/viper"

func LoadConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("conf")
	viper.SetConfigType("toml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
