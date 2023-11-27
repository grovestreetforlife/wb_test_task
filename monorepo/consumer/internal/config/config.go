package config

import (
	"github.com/creasty/defaults"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Database Database `yaml:"database"`
	Consumer Consumer `yaml:"consumer"`
	Cache    Cache    `yaml:"cache"`
	Jaeger   Jaeger   `yaml:"jaeger"`
}

type Database struct {
	PostgresDatabase PostgresDatabase `yaml:"postgres"`
}

type PostgresDatabase struct {
	Url         string `yaml:"url"`
	MaxOpenConn int    `yaml:"max_open_conn"`
	MaxConnLife int    `yaml:"max_conn_life"`
}

type Cache struct {
	RedisCache RedisCache `yaml:"redis"`
}

type RedisCache struct {
	Address   string `yaml:"address"`
	Password  string `yaml:"password"`
	TtlSecond int    `yaml:"ttl_second" default:"3600"`
}

type Consumer struct {
	NatsConsumer NatsConsumer `yaml:"nats"`
}

type NatsConsumer struct {
	Url                  string   `yaml:"url"`
	Subjects             []string `yaml:"subjects"`
	RetryOfFailedConnect bool     `yaml:"retry_of_failed_connect"`
	StreamName           string   `yaml:"stream_name"`
	CountConsumers       int      `yaml:"count_consumers" default:"2"`
}

type Jaeger struct {
	ServiceName              string  `yaml:"service_name"`
	Host                     string  `yaml:"host"`
	Port                     string  `yaml:"port"`
	Environment              string  `yaml:"environment"`
	TraceRatioFraction       float64 `yaml:"trace_ratio_fraction"`
	OTELExporterOTLPEndpoint string  `yaml:"otel_exporter_otlp_endpoint"`
}

func New(configFile string) (*Config, error) {
	var config Config
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "fail to upload yaml file")
	}

	err = defaults.Set(&config)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal yaml file")
	}
	return &config, nil
}
