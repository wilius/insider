package configs

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
	"text/template"
)

var instance Config

func init() {
	configData, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %w", err))
	}

	tmpl, err := template.New("config").
		Option("missingkey=zero").
		Funcs(template.FuncMap{
			"default": defaultFunc,
		}).
		Parse(string(configData))

	if err != nil {
		panic(fmt.Errorf("template parsing error: %w", err))
	}

	env := loadEnvVars()

	var out bytes.Buffer
	if err := tmpl.Execute(&out, env); err != nil {
		panic(fmt.Errorf("template execution error: %w", err))
	}

	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewBuffer(out.Bytes())); err != nil {
		panic(fmt.Errorf("error reading config: %w", err))
	}

	var config configImp
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}

	instance = config
}

func Instance() Config {
	return instance
}

func defaultFunc(value, defaultVal interface{}) interface{} {
	if value == nil || value == "" {
		return defaultVal
	}
	return value
}

func loadEnvVars() map[string]string {
	envMap := make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}
	return envMap
}
