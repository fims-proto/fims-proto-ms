package config

import "github.com/spf13/viper"

func setDefaults() {
	viper.SetDefault("profile", "dev")

	viper.SetDefault("app.port", "80")
	viper.SetDefault("app.multiTenancy", false)

	viper.SetDefault("gin.releaseMode", "debug")

	viper.SetDefault("logger.debug", true)
	viper.SetDefault("logger.jsonEncoding", false)
	viper.SetDefault("logger.showSql", true)
}
