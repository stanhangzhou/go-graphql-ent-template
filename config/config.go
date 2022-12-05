package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
	"gitlab.com/trustify/core/pkg/util/environment"
)

type config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSL      string
	}
	HttpServer struct {
		Port string
	}
}

var C config

type ReadConfigOption struct {
	AppEnv string
}

func ReadConfig(options ReadConfigOption) {
	Config := &C

	if environment.IsDev() {
		viper.AddConfigPath(filepath.Join(rootDir(), "config"))
		viper.SetConfigName("config")
	} else if environment.IsTest() || (options.AppEnv == environment.Test) {
		viper.AddConfigPath(filepath.Join(rootDir(), "config"))
		viper.SetConfigName("config.test")
	} else if environment.IsE2E() || options.AppEnv == environment.E2E {
		viper.AddConfigPath(filepath.Join(rootDir(), "config"))
		viper.SetConfigName("config.e2e")
	} else {
		// TODO: Handle production config
		fmt.Println("production config not implemented")
	}

	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("could not read configuration: %v", err)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalf("could not unmarshal config: %v", err)
		os.Exit(1)
	}
}

func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
