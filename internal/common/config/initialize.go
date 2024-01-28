package config

import (
	"bytes"
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

func Initialize() {
	// defaulting
	setDefaults()

	// bind environment variables
	flag.Parse()
	viper.MustBindEnv("profile", "PROFILE")
	viper.MustBindEnv("postgres.dsn", "DSN")

	// read config
	profile := viper.GetString("profile")
	viper.SetConfigName(fmt.Sprintf("application-%s", profile))
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to load config file: %w", err))
	}

	// check mandatory:
	checkResult := bytes.Buffer{}
	// postgres
	if !viper.IsSet("postgres.dsn") {
		checkResult.WriteString("postgres.dsn; ")
	}

	if checkResult.Len() > 0 {
		panic("config missing: " + checkResult.String())
	}
}
