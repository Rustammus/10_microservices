package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"sync"
)

type Config struct {
	Log     LogConfig     `yaml:"log"`
	REST    RESTConfig    `yaml:"rest"`
	RPC     RPCConfig     `yaml:"rpc"`
	Metrics MetricsConfig `yaml:"metrics"`
	//TokenTTL time.Duration `yaml:"token_ttl" env:"APP_TOKEN_TTL" env-default:"1h"`
}

type LogConfig struct {
	Level  slog.Level `yaml:"level" env:"APP_LOG_LEVEL" env-default:"debug"`
	Output string     `yaml:"output" env:"APP_LOG_OUTPUT"`
}

type RPCConfig struct {
	Port string `yaml:"port" env:"APP_RPC_PORT" env-default:"50051"`
}

type RESTConfig struct {
	Port string `yaml:"port" env:"APP_REST_PORT" env-default:"8082"`
}

type MetricsConfig struct {
	Path string `yaml:"path" env:"APP_METRICS_PATH" env-default:"/metrics"`
	Port string `yaml:"port" env:"APP_METRICS_PORT" env-default:"2112"`
}

var conf Config

var once sync.Once

func GetConfig() Config {
	once.Do(func() {
		err := cleanenv.ReadEnv(&conf)
		if err != nil {
			panic(err)
		}
	})
	return conf
}
