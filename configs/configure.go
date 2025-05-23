package configs

import (
	"bytes"
	"fmt"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
	"text/template"
	"time"
)

var instance Config

func Instance() Config {
	return instance
}

func init() {
	out := renderConfigTemplate()

	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewBuffer((*out).Bytes())); err != nil {
		panic(fmt.Errorf("error reading config: %w", err))
	}

	var config configImp
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringToDurationParserHook(),
		),
		Result:           &config,
		TagName:          "mapstructure",
		WeaklyTypedInput: true,
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		panic(fmt.Errorf("failed to create decoder: %w", err))
	}

	if err := decoder.Decode(viper.AllSettings()); err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}

	config.Provider = readInjectedProviderConfig()
	instance = config
}

func readInjectedProviderConfig() ProviderConfig {
	providerConfigMap := viper.GetStringMap("provider")
	providerConfig, err := providerConfigFactory(&providerConfigMap)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}
	return providerConfig
}

func renderConfigTemplate() *bytes.Buffer {
	configData, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %w", err))
	}

	tmpl, err := template.New("config").
		Option("missingkey=zero").
		Funcs(template.FuncMap{
			"default": templateDefaultFunction,
		}).
		Parse(string(configData))

	if err != nil {
		panic(fmt.Errorf("template parsing error: %w", err))
	}

	env := parseTemplateEnvironmentVariables()

	var out bytes.Buffer
	if err := tmpl.Execute(&out, env); err != nil {
		panic(fmt.Errorf("template execution error: %w", err))
	}

	return &out
}

func templateDefaultFunction(value, defaultVal interface{}) interface{} {
	if value == nil || value == "" {
		return defaultVal
	}
	return value
}

func parseTemplateEnvironmentVariables() map[string]string {
	envMap := make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}
	return envMap
}

func stringToDurationParserHook() mapstructure.DecodeHookFuncType {
	return func(
		from reflect.Type,
		to reflect.Type,
		data interface{},
	) (interface{}, error) {
		if from.Kind() == reflect.String && to == reflect.TypeOf(time.Duration(0)) {
			return time.ParseDuration(data.(string))
		}
		return data, nil
	}
}
