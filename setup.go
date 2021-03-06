package configmapper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v1"
)

const (
	EnableLocalConfigFlagEnv    = "LOCAL"
	FileEnvKey                  = "CONFIG_FILE_PATH"
	LocalConfigPath             = "LOCAL_CONFIG_FILE_PATH"
	DefaultConfigFilePath       = "/etc/app/config.toml"
	LocalConfigFilePathManifest = "manifests/configmap.yaml"
)

func loadViperFromConfig(configFilePath string, c interface{}) interface{} {
	viper.SetConfigFile(configFilePath)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if !errors.Is(err, nil) {
		panic(fmt.Errorf("config file read error: %w", err))
	}
	return loadViper(c)
}

func LoadViperFromManifest(manifestFilePath string, c interface{}) (interface{}, error) {
	buf, err := ioutil.ReadFile(manifestFilePath)
	if !errors.Is(err, nil) {
		return nil, fmt.Errorf("config file read error: %w", err)
	}
	data := make(map[string]map[string]string, 20)
	err = yaml.Unmarshal(buf, &data)
	if !errors.Is(err, nil) {
		return nil, fmt.Errorf("config file read error: %w", err)
	}

	b, ok := data["data"]["config.toml"]
	if !ok {
		panic(fmt.Errorf("config file read error"))
	}

	viper.SetConfigType("toml")
	err = viper.ReadConfig(bytes.NewBuffer([]byte(b)))
	if !errors.Is(err, nil) {
		return nil, fmt.Errorf("config file read error: %w", err)
	}
	viper.AutomaticEnv()

	return loadViper(c), nil
}

func loadViper(c interface{}) interface{} {
	if err := viper.Unmarshal(&c); !errors.Is(err, nil) {
		panic(fmt.Errorf("config file read error: %w", err))
	}

	// validation check
	validate := validator.New()
	if err := validate.Struct(c); !errors.Is(err, nil) {
		panic(fmt.Errorf("config file validation error: %w", err))
	}
	return c
}

// Initialize ...
func Initialize(config interface{}) (interface{}, error) {
	var confPath string
	if confPath = os.Getenv(FileEnvKey); confPath == "" {
		confPath = DefaultConfigFilePath
	}

	var conf interface{}
	var err error
	if localFlag := os.Getenv(EnableLocalConfigFlagEnv); localFlag != "" {
		if localConfigFilePath := os.Getenv(LocalConfigPath); localConfigFilePath != "" {
			fmt.Printf("config file is %v\n", localConfigFilePath)
			conf, err = LoadViperFromManifest(localConfigFilePath, config)
		} else {
			fmt.Printf("config file is %v\n", LocalConfigFilePathManifest)
			conf, err = LoadViperFromManifest(LocalConfigFilePathManifest, config)
		}
	} else {
		// Output the path to standard output in case of a read error in the configuration file
		// (Logger is not used because its before Logger initialization)
		fmt.Printf("config file is %v\n", confPath)
		conf = loadViperFromConfig(confPath, config)
	}
	if err != nil {
		return nil, err
	}
	b, e := json.Marshal(conf)
	if e != nil {
		return nil, fmt.Errorf("config output json marshal err %w", e)
	}
	fmt.Printf("config loaded %s\n", string(b))

	return conf, nil
}
