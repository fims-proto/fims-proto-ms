package config

import "github.com/spf13/viper"

func setDefaults() {
	viper.SetDefault("profile", "dev")

	viper.SetDefault("app.port", "80")
	viper.SetDefault("app.multiTenancy", false)

	viper.SetDefault("logger.debug", false)
	viper.SetDefault("logger.jsonEncoding", true)
	viper.SetDefault("logger.showSql", false)
}
